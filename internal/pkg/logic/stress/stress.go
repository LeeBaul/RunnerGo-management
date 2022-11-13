package stress

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-omnibus/omnibus"
	"github.com/go-omnibus/proof"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gen"
	"io"
	"io/ioutil"
	"kp-management/internal/pkg/biz/errno"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/load"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/run_plan"
)

type Baton struct {
	Ctx      context.Context
	PlanID   int64
	TeamID   int64
	UserID   int64
	SceneIDs []int64

	plan            *model.Plan
	scenes          []*model.Target
	task            map[int64]*mao.Task // sceneID 对应任务配置
	globalVariables []*model.Variable
	flows           []*mao.Flow
	sceneVariables  []*model.Variable
	importVariables []*model.VariableImport
	reports         []*model.Report
	balance         *DispatchMachineBalance
	stress          []*run_plan.Stress
}

type UsableMachineMap struct {
	IP               string // IP地址(包含端口号)
	Region           string // 机器所属区域
	Weight           int64  // 权重
	UsableGoroutines int64  // 可用协程数
}

// 压力机心跳上报数据
type HeartBeat struct {
	Name              string        `json:"name"`               // 机器名称
	CpuUsage          float64       `json:"cpu_usage"`          // CPU使用率
	CpuLoad           *load.AvgStat `json:"cpu_load"`           // CPU负载信息
	MemInfo           []MemInfo     `json:"mem_info"`           // 内存使用情况
	Networks          []Network     `json:"networks"`           // 网络连接情况
	DiskInfos         []DiskInfo    `json:"disk_infos"`         // 磁盘IO情况
	MaxGoroutines     int64         `json:"max_goroutines"`     // 当前机器支持最大协程数
	CurrentGoroutines int64         `json:"current_goroutines"` // 当前已用协程数
	ServerType        int64         `json:"server_type"`        // 压力机类型：0-主力机器，1-备用机器
	CreateTime        int64         `json:"create_time"`        // 数据上报时间（时间戳）
}

type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type Network struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

type Stress interface {
	Execute(baton *Baton) (int, error)
	SetNext(Stress)
}

type CheckIdleMachine struct {
	next Stress
}

func (s *CheckIdleMachine) Execute(baton *Baton) (int, error) {
	// 从Redis获取压力机列表
	machineListRes := dal.RDB.HGetAll(consts.MachineListRedisKey)
	if len(machineListRes.Val()) == 0 || machineListRes.Err() != nil {
		// todo 后面可能增加兜底策略
		proof.Infof("没有上报上来的空闲压力机可用")
		return errno.ErrResourceNotEnough, fmt.Errorf("资源不足")
	}

	baton.balance = &DispatchMachineBalance{}

	usableMachineMap := new(UsableMachineMap)                                     // 单个压力机基本数据
	usableMachineSlice := make([]*UsableMachineMap, 0, len(machineListRes.Val())) // 所有上报过来的压力机切片
	var minWeight int64                                                           // 所有可用压力机里面最小的权重的值
	var inUseMachineNum int                                                       // 所有有任务在运行的压力机数量

	var breakFor = false

	// 查到了机器列表，然后判断可用性
	var runnerMachineInfo HeartBeat
	for machineAddr, machineDetail := range machineListRes.Val() {
		// 把机器详情信息解析成格式化数据
		err := json.Unmarshal([]byte(machineDetail), &runnerMachineInfo)
		if err != nil {
			proof.Infof("runner_machine_detail 数据解析失败 err：", err)
			continue
		}

		// 压力机数据上报时间超过3秒，则认为服务不可用，不参与本次压力测试
		nowTime := time.Now().Unix()
		if nowTime-runnerMachineInfo.CreateTime > 3 {
			proof.Infof("当前压力机上报心跳数据超时，暂不可用")
			continue
		}

		// 判断当前压力机性能是否爆满,如果某个指标爆满，则不参与本次压力测试
		if runnerMachineInfo.CpuUsage >= 65 { // CPU使用判断
			continue
		}
		for _, memInfo := range runnerMachineInfo.MemInfo { // 内存使用判断
			if memInfo.UsedPercent >= 65 {
				breakFor = true
				break
			}
		}
		for _, diskInfo := range runnerMachineInfo.DiskInfos { // 磁盘使用判断
			if diskInfo.UsedPercent >= 55 {
				breakFor = true
				break
			}
		}

		// 最后判断是否结束当前循环
		if breakFor {
			continue
		}

		machineAddrSlice := strings.Split(machineAddr, "_")
		if len(machineAddrSlice) != 3 {
			continue
		}

		// 当前机器可用协程数
		usableGoroutines := runnerMachineInfo.MaxGoroutines - runnerMachineInfo.CurrentGoroutines

		// 组装可用机器结构化数据
		usableMachineMap.IP = machineAddrSlice[0] + ":" + machineAddrSlice[1]
		usableMachineMap.UsableGoroutines = usableGoroutines
		usableMachineMap.Weight = usableGoroutines
		usableMachineSlice = append(usableMachineSlice, usableMachineMap)

		// 获取当前压力机当中最小的权重值
		if minWeight == 0 || minWeight > usableGoroutines {
			minWeight = usableGoroutines
		}

		// 获取当前机器是否使用当中
		machineUseStateKey := consts.MachineUseStatePrefix + machineAddrSlice[0] + ":" + machineAddrSlice[1]
		useStateVal, _ := dal.RDB.Get(machineUseStateKey).Result()
		if useStateVal != "" {
			inUseMachineNum++
		}

	}

	for _, machineInfo := range usableMachineSlice {
		if inUseMachineNum < len(usableMachineSlice) {
			// 获取当前机器是否使用当中
			machineUseStateKey := consts.MachineUseStatePrefix + machineInfo.IP
			useStateVal, _ := dal.RDB.Get(machineUseStateKey).Result()
			if useStateVal != "" {
				machineInfo.UsableGoroutines = int64(minWeight) - 1
				if machineInfo.UsableGoroutines <= 0 {
					machineInfo.UsableGoroutines = 1
				}
			}
		}
	}

	sort.Slice(usableMachineSlice, func(i, j int) bool {
		return usableMachineSlice[i].UsableGoroutines > usableMachineSlice[j].UsableGoroutines
	})

	// 按当前顺序把机器放到备用列表
	for _, machineInfo := range usableMachineSlice {
		addErr := baton.balance.AddMachine(fmt.Sprintf("%s", machineInfo.IP))
		if addErr != nil {
			continue
		}
	}

	if len(baton.balance.rss) == 0 {
		proof.Infof("当前没有空闲压力机可用")
		return errno.ErrResourceNotEnough, fmt.Errorf("资源不足")
	}

	return s.next.Execute(baton)
}

func (s *CheckIdleMachine) SetNext(stress Stress) {
	s.next = stress
}

type AssemblePlan struct {
	next Stress
}

func (s *AssemblePlan) Execute(baton *Baton) (int, error) {
	tx := query.Use(dal.DB()).Plan
	p, err := tx.WithContext(baton.Ctx).Where(tx.ID.Eq(baton.PlanID)).First()
	if err != nil {
		return errno.ErrMysqlFailed, err
	}
	baton.plan = p
	return s.next.Execute(baton)
}

func (s *AssemblePlan) SetNext(stress Stress) {
	s.next = stress
}

type AssembleScenes struct {
	next Stress
}

func (s *AssembleScenes) Execute(baton *Baton) (int, error) {
	tx := query.Use(dal.DB()).Target

	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.PlanID.Eq(baton.PlanID))
	conditions = append(conditions, tx.TargetType.Eq(consts.TargetTypeScene))
	conditions = append(conditions, tx.Status.Eq(consts.TargetStatusNormal))
	if len(baton.SceneIDs) > 0 {
		conditions = append(conditions, tx.ID.In(baton.SceneIDs...))
	}

	scenes, err := tx.WithContext(baton.Ctx).Where(conditions...).Find()
	if err != nil {
		return errno.ErrMysqlFailed, err
	}

	baton.scenes = scenes
	return s.next.Execute(baton)
}

func (s *AssembleScenes) SetNext(stress Stress) {
	s.next = stress
}

type AssembleTask struct {
	next Stress
}

func (s *AssembleTask) Execute(baton *Baton) (int, error) {
	memo := make(map[int64]*mao.Task)

	// 判断参数是否包含scene_ids
	if len(baton.SceneIDs) > 0 { // 包含则说明当前任务时定时任务
		tx := dal.GetQuery().TimedTaskConf
		timedTaskConfInfo, err := tx.WithContext(baton.Ctx).Where(tx.SenceID.Eq(baton.SceneIDs[0])).First()
		if err != nil {
			return errno.ErrMysqlFailed, err
		}

		var modeConf mao.ModeConf
		err = json.Unmarshal([]byte(timedTaskConfInfo.ModeConf), &modeConf)
		if err != nil {
			proof.Errorf("运行定时任务--解析任务配置文件失败")
			return errno.ErrUnMarshalFailed, err
		}
		memo[baton.SceneIDs[0]] = &mao.Task{
			PlanID:   timedTaskConfInfo.PlanID,
			SceneID:  timedTaskConfInfo.SenceID,
			TaskType: timedTaskConfInfo.TaskType,
			TaskMode: timedTaskConfInfo.TaskMode,
			ModeConf: &modeConf,
		}
	} else { // 普通任务
		collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
		cur, err := collection.Find(baton.Ctx, bson.D{{"plan_id", baton.PlanID}})
		if err != nil {
			return errno.ErrMongoFailed, err
		}

		var task []*mao.Task
		if err := cur.All(baton.Ctx, &task); err != nil {
			return errno.ErrMongoFailed, err
		}

		for _, t := range task {
			memo[t.SceneID] = t
		}
	}

	baton.task = memo
	return s.next.Execute(baton)
}

func (s *AssembleTask) SetNext(stress Stress) {
	s.next = stress
}

type AssembleGlobalVariables struct {
	next Stress
}

func (s *AssembleGlobalVariables) Execute(baton *Baton) (int, error) {
	tx := query.Use(dal.DB()).Variable
	variables, err := tx.WithContext(baton.Ctx).Where(
		tx.TeamID.Eq(baton.TeamID),
		tx.Type.Eq(consts.VariableTypeGlobal),
	).Find()

	if err != nil {
		return errno.ErrMysqlFailed, err
	}

	baton.globalVariables = variables
	return s.next.Execute(baton)
}

func (s *AssembleGlobalVariables) SetNext(stress Stress) {
	s.next = stress
}

type AssembleFlows struct {
	next Stress
}

func (s *AssembleFlows) Execute(baton *Baton) (int, error) {
	var sceneIDs []int64
	for _, scene := range baton.scenes {
		sceneIDs = append(sceneIDs, scene.ID)
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	cur, err := collection.Find(baton.Ctx, bson.D{{"scene_id", bson.D{{"$in", sceneIDs}}}})
	if err != nil {
		return errno.ErrMongoFailed, err
	}

	var flows []*mao.Flow
	if err := cur.All(baton.Ctx, &flows); err != nil {
		return errno.ErrMongoFailed, err
	}

	if len(flows) != len(sceneIDs) {
		proof.Infof("场景不能为空")
		return errno.ErrEmptyScene, errors.New("场景不能为空")
	}

	baton.flows = flows
	return s.next.Execute(baton)
}

func (s *AssembleFlows) SetNext(stress Stress) {
	s.next = stress
}

type AssembleSceneVariables struct {
	next Stress
}

func (s *AssembleSceneVariables) Execute(baton *Baton) (int, error) {
	var sceneIDs []int64
	for _, scene := range baton.scenes {
		sceneIDs = append(sceneIDs, scene.ID)
	}

	tx := query.Use(dal.DB()).Variable
	variables, err := tx.WithContext(baton.Ctx).Where(
		tx.TeamID.Eq(baton.TeamID),
		tx.SceneID.In(sceneIDs...),
		tx.Type.Eq(consts.VariableTypeScene),
	).Find()

	if err != nil {
		return errno.ErrMysqlFailed, err
	}

	baton.sceneVariables = variables
	return s.next.Execute(baton)
}

func (s *AssembleSceneVariables) SetNext(stress Stress) {
	s.next = stress
}

type AssembleImportVariables struct {
	next Stress
}

func (s *AssembleImportVariables) Execute(baton *Baton) (int, error) {
	var sceneIDs []int64
	for _, scene := range baton.scenes {
		sceneIDs = append(sceneIDs, scene.ID)
	}

	tx := query.Use(dal.DB()).VariableImport
	vis, err := tx.WithContext(baton.Ctx).Where(tx.SceneID.In(sceneIDs...)).Find()
	if err != nil {
		return errno.ErrMysqlFailed, err
	}

	baton.importVariables = vis
	return s.next.Execute(baton)
}

func (s *AssembleImportVariables) SetNext(stress Stress) {
	s.next = stress
}

type MakeReport struct {
	next Stress
}

func (s *MakeReport) Execute(baton *Baton) (int, error) {
	tx := query.Use(dal.DB()).Report

	cnt, err := tx.WithContext(baton.Ctx).Unscoped().Where(tx.TeamID.Eq(baton.TeamID)).Count()
	if err != nil {
		return errno.ErrMysqlFailed, err
	}

	reports := make([]*model.Report, 0)
	for i, scene := range baton.scenes {
		if _, ok := baton.task[scene.ID]; !ok {
			return errno.ErrMustTaskInit, errors.New("当前场景没有配置任务类型或任务模式，场景id：" + strconv.Itoa(int(scene.ID)))
		}
		reports = append(reports, &model.Report{
			Rank:      cnt + 1 + omnibus.DefiniteInt64(i),
			TeamID:    scene.TeamID,
			PlanID:    baton.plan.ID,
			PlanName:  baton.plan.Name,
			SceneID:   scene.ID,
			SceneName: scene.Name,
			TaskType:  baton.task[scene.ID].TaskType,
			TaskMode:  baton.task[scene.ID].TaskMode,
			Status:    consts.ReportStatusNormal,
			RanAt:     time.Now(),
			RunUserID: baton.UserID,
		})
	}

	if err := tx.WithContext(baton.Ctx).CreateInBatches(reports, 10); err != nil {
		return errno.ErrMysqlFailed, err
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectReportTask)
	for _, report := range reports {

		_, err := collection.InsertOne(baton.Ctx, &mao.ReportTask{
			ReportID: report.ID,
			TaskType: report.TaskType,
			TaskMode: report.TaskMode,
			PlanID:   baton.plan.ID,
			PlanName: baton.plan.Name,
			ModeConf: &mao.ModeConf{
				ReheatTime:       baton.task[report.SceneID].ModeConf.ReheatTime,
				RoundNum:         baton.task[report.SceneID].ModeConf.RoundNum,
				Concurrency:      baton.task[report.SceneID].ModeConf.Concurrency,
				ThresholdValue:   baton.task[report.SceneID].ModeConf.ThresholdValue,
				StartConcurrency: baton.task[report.SceneID].ModeConf.StartConcurrency,
				Step:             baton.task[report.SceneID].ModeConf.Step,
				StepRunTime:      baton.task[report.SceneID].ModeConf.StepRunTime,
				MaxConcurrency:   baton.task[report.SceneID].ModeConf.MaxConcurrency,
				Duration:         baton.task[report.SceneID].ModeConf.Duration,
			},
		})

		if err != nil {
			return errno.ErrMongoFailed, err
		}
	}

	baton.reports = reports
	return s.next.Execute(baton)
}

func (s *MakeReport) SetNext(stress Stress) {
	s.next = stress
}

type MakeStress struct {
	next Stress
}

func (s *MakeStress) Execute(baton *Baton) (int, error) {

	for _, report := range baton.reports {
		for _, scene := range baton.scenes {
			for _, flow := range baton.flows {

				if scene.ID == report.SceneID && scene.ID == flow.SceneID {

					globalVariables := make([]*run_plan.Variable, 0)
					for _, v := range baton.globalVariables {
						globalVariables = append(globalVariables, &run_plan.Variable{
							Var: v.Var,
							Val: v.Val,
						})
					}

					var nodes run_plan.Nodes
					if err := bson.Unmarshal(flow.Nodes, &nodes); err != nil {
						proof.Errorf("node bson unmarshal err:%v", err)
						continue
					}

					sceneVariables := make([]*run_plan.Variable, 0)
					for _, v := range baton.sceneVariables {
						sceneVariables = append(sceneVariables, &run_plan.Variable{
							Var: v.Var,
							Val: v.Val,
						})
					}

					importVariables := make([]string, 0)
					for _, v := range baton.importVariables {
						importVariables = append(importVariables, v.URL)
					}

					if _, ok := baton.task[scene.ID]; !ok {
						return errno.ErrMustTaskInit, errors.New("请填写任务配置并保存")
					}
					req := run_plan.Stress{
						PlanID:     baton.plan.ID,
						PlanName:   baton.plan.Name,
						ReportID:   omnibus.DefiniteString(report.ID),
						TeamID:     baton.TeamID,
						ReportName: baton.plan.Name,
						ConfigTask: &run_plan.ConfigTask{
							TaskType: baton.task[scene.ID].TaskType,
							Mode:     baton.task[scene.ID].TaskMode,
							Remark:   baton.plan.Remark,
							CronExpr: baton.plan.CronExpr,
							ModeConf: &run_plan.ModeConf{
								ReheatTime:       baton.task[scene.ID].ModeConf.ReheatTime,
								RoundNum:         baton.task[scene.ID].ModeConf.RoundNum,
								Concurrency:      baton.task[scene.ID].ModeConf.Concurrency,
								ThresholdValue:   baton.task[scene.ID].ModeConf.ThresholdValue,
								StartConcurrency: baton.task[scene.ID].ModeConf.StartConcurrency,
								Step:             baton.task[scene.ID].ModeConf.Step,
								StepRunTime:      baton.task[scene.ID].ModeConf.StepRunTime,
								MaxConcurrency:   baton.task[scene.ID].ModeConf.MaxConcurrency,
								Duration:         baton.task[scene.ID].ModeConf.Duration,
							},
						},
						Variable: globalVariables,
						Scene: &run_plan.Scene{
							SceneID:                 scene.ID,
							EnablePlanConfiguration: false,
							SceneName:               scene.Name,
							TeamID:                  baton.TeamID,
							Nodes:                   nodes.Nodes,
							Configuration: &run_plan.SceneConfiguration{
								ParameterizedFile: &run_plan.SceneVariablePath{
									Path: importVariables,
								},
								Variable: sceneVariables,
							},
						},
					}

					baton.stress = append(baton.stress, &req)

				}
			}
		}
	}

	return s.next.Execute(baton)
}

func (s *MakeStress) SetNext(stress Stress) {
	s.next = stress
}

type SplitStress struct {
	next Stress
}

func (s *SplitStress) Execute(baton *Baton) (int, error) {
	memo := make(map[string]int32)
	for i, stress := range baton.stress {
		memo[stress.ReportID]++

		var apiCnt int64
		for _, node := range stress.Scene.Nodes {
			if node.Type == "api" {
				apiCnt++
			}
		}

		maxConcurrency := conf.Conf.Base.MaxConcurrency
		totalConcurrency := apiCnt * stress.ConfigTask.ModeConf.Concurrency
		for totalConcurrency > maxConcurrency {
			baton.stress[i].ConfigTask.ModeConf.Concurrency -= maxConcurrency / apiCnt

			baton.stress = append(baton.stress, baton.stress[i])

			memo[stress.ReportID]++

			totalConcurrency -= maxConcurrency
		}
	}

	for _, stress := range baton.stress {
		stress.MachineNum = memo[stress.ReportID]
	}

	return s.next.Execute(baton)
}

func (s *SplitStress) SetNext(stress Stress) {
	s.next = stress
}

type SplitImportVariable struct {
	next Stress
}

func (s *SplitImportVariable) Execute(baton *Baton) (int, error) {

	reportMemo := make(map[string]int)
	pathMemo := make(map[string]string)
	for _, stress := range baton.stress {
		for _, pathString := range stress.Scene.Configuration.ParameterizedFile.Path {
			pathMemo[stress.ReportID] = pathString
			reportMemo[stress.ReportID] += 1
		}
	}

	var reportPathMut sync.Mutex
	reportPathMemo := make(map[string][]string)
	for reportID, p := range pathMemo {
		fileExt := path.Ext(p)
		if fileExt != ".txt" && fileExt != ".csv" {
			continue
		}

		resp, err := http.Get(p)
		if err != nil {
			return errno.ErrHttpFailed, err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errno.ErrHttpFailed, err
		}

		files := omnibus.Explode("/", p)
		localFilePath := fmt.Sprintf("/tmp/%s", files[len(files)-1])
		if err := ioutil.WriteFile(localFilePath, data, 0644); err != nil {
			return errno.ErrHttpFailed, err
		}

		file, _ := os.Open(localFilePath)
		defer file.Close()

		var wg sync.WaitGroup
		ch := make(chan string)

		for i := 0; i < reportMemo[reportID]; i++ {
			wg.Add(1)

			/*协程任务：从管道中拉取数据并写入到文件中*/
			go func(indx int) {
				f, err := os.OpenFile(localFilePath+strconv.Itoa(indx)+fileExt, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {

				}
				defer f.Close()

				for lineStr := range ch {
					//向文件中写出UTF-8字符串
					f.WriteString(lineStr)
				}

				//todo oss
				reportPathMut.Lock()
				defer reportPathMut.Unlock()
				reportPathMemo[reportID] = append(reportPathMemo[reportID], localFilePath+strconv.Itoa(indx)+fileExt)
				wg.Done()
			}(i)
		}

		//创建缓冲读取器
		reader := bufio.NewReader(file)
		for {
			//读取一行字符串（编码为UTF-8）
			lineStr, err := reader.ReadString('\n')

			//读取完毕时，关闭所有数据管道，并退出读取
			if err == io.EOF {
				close(ch)
				break
			}

			ch <- lineStr
		}

		//阻塞等待所有协程结束任务
		wg.Wait()
	}

	for _, stress := range baton.stress {
		if len(stress.Scene.Configuration.ParameterizedFile.Path) > 0 {
			stress.Scene.Configuration.ParameterizedFile.Path[0] = reportPathMemo[stress.ReportID][0]
			reportPathMemo[stress.ReportID] = reportPathMemo[stress.ReportID][1:]
		}

	}

	return s.next.Execute(baton)
}

func (s *SplitImportVariable) SetNext(stress Stress) {
	s.next = stress
}

type RunMachineStress struct {
	next Stress
}

func (s *RunMachineStress) Execute(baton *Baton) (int, error) {
	// 当前可用压力机数量
	machinesNum := len(baton.balance.rss)
	curIndex := 0 // 当前使用的压力机数组下标

	for _, stress := range baton.stress {
		breakState := false
		forNum := 0
		for {
			if forNum == machinesNum {
				breakState = true
				break
			}

			if curIndex == machinesNum-1 {
				curIndex = 0
			} else {
				curIndex++
			}
			forNum++

			addr := baton.balance.GetMachine(curIndex) // 获取当前压力机地址
			if addr == "" {
				continue
			}
			machinesState := GetRunnerMachineState(addr) // 获取当前压力机可用状态

			if machinesState { // 如果当前机器可用
				// 把当前机器信息写入到数据表当中
				tx := query.Use(dal.DB()).ReportMachine
				_, err := tx.WithContext(baton.Ctx).Where(tx.ReportID.Eq(omnibus.DefiniteInt64(stress.ReportID))).Assign(
					tx.ReportID.Value(omnibus.DefiniteInt64(stress.ReportID)),
					tx.IP.Value(omnibus.Explode(":", addr)[0]),
				).FirstOrCreate()
				if err != nil {
					_ = DeletePlanReport(baton)
					proof.Infof("把报告和对应机器写入到数据库失败，err：", err)
					return errno.ErrMysqlFailed, err
				}

				// 增加分区字段判断
				partition := GetPartition()
				if partition == -1 {
					proof.Infof("当前没有可用的kafka分区")
					_ = DeletePlanReport(baton)
					return errno.ErrResourceNotEnough, fmt.Errorf("资源不足")
				}
				stress.Partition = partition

				_, err = resty.New().R().SetBody(stress).Post(fmt.Sprintf("http://%s/runner/run_plan", addr))
				proof.Infof("请求压力机运行情况，req：%+v。 err： %+v。", proof.Render("req", stress), err)
				if err != nil {
					//// 如果调用施压接口失败，则删除掉当前的这个报告id
					//reportTable := query.Use(dal.DB()).Report
					//_, err2 := reportTable.WithContext(baton.Ctx).Where(reportTable.ID.Eq(omnibus.DefiniteInt64(stress.ReportID))).Delete()
					//if err2 != nil {
					//	return errno.ErrMysqlFailed, err2
					//}
					_ = DeletePlanReport(baton)
					proof.Infof("请求压力机进行压测失败，err：", err)
					return errno.ErrHttpFailed, err
				}

				// 把当前压力机使用状态设置到redis当中
				machineUseStateKey := consts.MachineUseStatePrefix + addr
				dal.RDB.SetNX(machineUseStateKey, 1, 3600*24)

				p := query.Use(dal.DB()).Plan
				_, err = p.WithContext(baton.Ctx).Where(p.ID.Eq(baton.PlanID)).UpdateColumn(p.Status, consts.PlanStatusUnderway)
				if err != nil {
					proof.Infof("把压力机使用情况写入redis是吧，err：", err)
					return errno.ErrMysqlFailed, err
				}
				break
			} else {
				continue
			}
		}

		if breakState {
			// todo 报警
			proof.Infof("当前没有可用的压力机，或所有压力机状态爆满")
			_ = DeletePlanReport(baton)
			return errno.ErrResourceNotEnough, errors.New("资源不足")
		}
	}

	return errno.Ok, nil
}

func (s *RunMachineStress) SetNext(stress Stress) {
	s.next = stress
}

// GetPartition 获取可用分区
func GetPartition() int32 {
	//默认分区为0
	var partition int32 = -1 //默认为-1 表示不可用分区锁
	// kafka全局的报告分区key名
	partitionLock := consts.KafkaReportPartition
	//目前kafka总分有5个分区，随机取出一个
	totalPartitionNum := consts.KafkaReportPartitionNum
	for i := 0; i < totalPartitionNum; i++ {
		// 获取当前时间戳
		nowTime := time.Now().Unix()
		// 把分区转换成字符串
		partitionNumString := strconv.Itoa(i)
		// 尝试获取当前分区锁
		res, _ := dal.RDB.HSetNX(partitionLock, partitionNumString, nowTime).Result()
		if res == false { // 获取失败或者当前分区锁已被占用
			continue
		} else {
			partition = int32(i)
			break
		}
	}
	return partition
}

// GetRunnerMachineState 获取当前压力机是否可用
func GetRunnerMachineState(addr string) bool {
	// 从Redis获取压力机列表
	machineListRes := dal.RDB.HGetAll(consts.MachineListRedisKey)
	if len(machineListRes.Val()) == 0 || machineListRes.Err() != nil {
		proof.Infof("当前没有上报上来的空闲压力机可用")
		return false
	}

	// 退出循环的标识
	var breakFor = false
	//// 查到了机器列表，然后判断可用性
	var runnerMachineInfo *HeartBeat
	// 初始化机器状态map
	machineStateMap := make(map[string]bool, len(machineListRes.Val()))
	for machineAddr, machineDetail := range machineListRes.Val() {
		// 解析hash的field字段
		machineAddrSlice := strings.Split(machineAddr, "_")
		if len(machineAddrSlice) != 3 {
			continue
		}

		// 组装可用机器map的key
		addrKey := machineAddrSlice[0] + ":" + machineAddrSlice[1]
		machineStateMap[addrKey] = false

		// 把机器详情信息解析成格式化数据
		err := json.Unmarshal([]byte(machineDetail), &runnerMachineInfo)
		if err != nil {
			proof.Infof("压力机数据解析失败，err:", err)
			continue
		}

		// 压力机数据上报时间超过3秒，则认为服务不可用，不参与本次压力测试
		nowTime := time.Now().Unix()
		if runnerMachineInfo.CreateTime < nowTime-3 {
			continue
		}

		// 判断当前压力机性能是否爆满,如果某个指标爆满，则不参与本次压力测试
		if runnerMachineInfo.CpuUsage >= 65 { // CPU使用判断
			continue
		}
		for _, memInfo := range runnerMachineInfo.MemInfo { // 内存使用判断
			if memInfo.UsedPercent >= 65 {
				breakFor = true
				break
			}
		}
		for _, diskInfo := range runnerMachineInfo.DiskInfos { // 磁盘使用判断
			if diskInfo.UsedPercent >= 55 {
				breakFor = true
				break
			}
		}

		// 最后判断是否结束当前循环
		if breakFor {
			continue
		}

		// 当前机器可用协程数
		usableGoroutines := runnerMachineInfo.MaxGoroutines - runnerMachineInfo.CurrentGoroutines
		if usableGoroutines <= 0 {
			continue
		}
		machineStateMap[addr] = true
	}

	if machineState, ok := machineStateMap[addr]; !ok {
		return false
	} else {
		return machineState
	}
}

// 删除执行失败的计划下的所有报告
func DeletePlanReport(baton *Baton) error {
	for _, reportInfo := range baton.reports {
		// 如果调用施压接口失败，则删除掉当前的这个报告id
		reportTable := query.Use(dal.DB()).Report
		_, err := reportTable.WithContext(baton.Ctx).Where(reportTable.ID.Eq(omnibus.DefiniteInt64(reportInfo.ID))).Delete()
		if err != nil {
			proof.Infof("运行计划--删除报告失败，报告id为：", reportInfo.ID)
		}
	}
	return nil
}
