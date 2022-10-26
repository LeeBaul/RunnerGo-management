package stress

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/go-omnibus/omnibus"
	"github.com/go-omnibus/proof"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gen"

	"kp-controller/internal/pkg/conf"
	"kp-controller/internal/pkg/consts"
	"kp-controller/internal/pkg/dal"
	"kp-controller/internal/pkg/dal/mao"
	"kp-controller/internal/pkg/dal/model"
	"kp-controller/internal/pkg/dal/prometheus"
	"kp-controller/internal/pkg/dal/query"
	"kp-controller/internal/pkg/dal/rao"
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
	balance         *WeightRoundRobinBalance
	stress          []*rao.Stress
}

type Stress interface {
	Execute(baton *Baton) error
	SetNext(Stress)
}

type CheckIdleMachine struct {
	next Stress
}

func (s *CheckIdleMachine) Execute(baton *Baton) error {
	tx := query.Use(dal.DB()).Machine
	machines, err := tx.WithContext(baton.Ctx).Where(tx.Status.Eq(consts.MachineStatusIdle), tx.UpdatedAt.Between(time.Now().Add(-15*time.Second), time.Now())).Find()
	if err != nil {
		return err
	}
	if len(machines) == 0 {
		return fmt.Errorf("not idle machine")
	}

	baton.balance = &WeightRoundRobinBalance{}
	for _, machine := range machines {
		// 过滤cpu
		c, err := prometheus.GetCPUUsage(baton.Ctx, machine.IP)
		if err != nil {
			proof.Errorf("prometheus get cpu err: %s", err.Error())
			continue
		}
		if c > conf.Conf.Machine.Threshold.CPU || c == 0 {
			continue
		}

		// 过滤内存
		//m, err := prometheus.GetMemUsage(baton.Ctx, machine.IP)
		//if err != nil {
		//	proof.Errorf("prometheus get mem err: %s", err.Error())
		//	continue
		//}
		//if m > conf.Conf.Machine.Threshold.MEM || m == 0 {
		//	continue
		//}

		// 过滤磁盘io
		//d, err := prometheus.GetDiskIOUsage(baton.Ctx, machine.IP)
		//if err != nil {
		//	proof.Errorf("prometheus get disk io err: %s", err.Error())
		//	continue
		//}
		//if d > conf.Conf.Machine.Threshold.Disk || d == 0 {
		//	continue
		//}

		// 过滤网络io
		//n, err := prometheus.GetNetIOUsage(baton.Ctx, machine.IP)
		//if err != nil {
		//	proof.Errorf("prometheus get net io err: %s", err.Error())
		//	continue
		//}
		//if n > conf.Conf.Machine.Threshold.Net || n == 0 {
		//	continue
		//}

		baton.balance.Add(fmt.Sprintf("%s:%d", machine.IP, machine.Port), omnibus.DefiniteString(machine.Weight))
	}

	if len(baton.balance.rss) == 0 {
		return fmt.Errorf("empty idle machine")
	}

	return s.next.Execute(baton)
}

func (s *CheckIdleMachine) SetNext(stress Stress) {
	s.next = stress
}

type AssemblePlan struct {
	next Stress
}

func (s *AssemblePlan) Execute(baton *Baton) error {
	tx := query.Use(dal.DB()).Plan
	p, err := tx.WithContext(baton.Ctx).Where(tx.ID.Eq(baton.PlanID)).First()
	if err != nil {
		return err
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

func (s *AssembleScenes) Execute(baton *Baton) error {
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
		return err
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

func (s *AssembleTask) Execute(baton *Baton) error {
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
	cur, err := collection.Find(baton.Ctx, bson.D{{"plan_id", baton.PlanID}})
	if err != nil {
		return err
	}

	var task []*mao.Task
	if err := cur.All(baton.Ctx, &task); err != nil {
		return err
	}

	memo := make(map[int64]*mao.Task)
	for _, t := range task {
		memo[t.SceneID] = t
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

func (s *AssembleGlobalVariables) Execute(baton *Baton) error {
	tx := query.Use(dal.DB()).Variable
	variables, err := tx.WithContext(baton.Ctx).Where(
		tx.TeamID.Eq(baton.TeamID),
		tx.Type.Eq(consts.VariableTypeGlobal),
	).Find()

	if err != nil {
		return err
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

func (s *AssembleFlows) Execute(baton *Baton) error {
	var sceneIDs []int64
	for _, scene := range baton.scenes {
		sceneIDs = append(sceneIDs, scene.ID)
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	cur, err := collection.Find(baton.Ctx, bson.D{{"scene_id", bson.D{{"$in", sceneIDs}}}})
	if err != nil {
		return err
	}

	var flows []*mao.Flow
	if err := cur.All(baton.Ctx, &flows); err != nil {
		return err
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

func (s *AssembleSceneVariables) Execute(baton *Baton) error {
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
		return err
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

func (s *AssembleImportVariables) Execute(baton *Baton) error {
	var sceneIDs []int64
	for _, scene := range baton.scenes {
		sceneIDs = append(sceneIDs, scene.ID)
	}

	tx := query.Use(dal.DB()).VariableImport
	vis, err := tx.WithContext(baton.Ctx).Where(tx.SceneID.In(sceneIDs...)).Find()
	if err != nil {
		return err
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

func (s *MakeReport) Execute(baton *Baton) error {
	tx := query.Use(dal.DB()).Report

	cnt, err := tx.WithContext(baton.Ctx).Unscoped().Where(tx.TeamID.Eq(baton.TeamID)).Count()
	if err != nil {
		return err
	}

	reports := make([]*model.Report, 0)
	for i, scene := range baton.scenes {
		reports = append(reports, &model.Report{
			Rank:      cnt + 1 + omnibus.DefiniteInt64(i),
			TeamID:    scene.TeamID,
			PlanID:    baton.plan.ID,
			PlanName:  baton.plan.Name,
			SceneID:   scene.ID,
			SceneName: scene.Name,
			TaskType:  baton.plan.TaskType,
			TaskMode:  baton.plan.Mode,
			Status:    consts.ReportStatusNormal,
			RanAt:     time.Now(),
			RunUserID: baton.UserID,
		})
	}

	if err := tx.WithContext(baton.Ctx).CreateInBatches(reports, 10); err != nil {
		return err
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
			return err
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

func (s *MakeStress) Execute(baton *Baton) error {

	for _, report := range baton.reports {
		for _, scene := range baton.scenes {
			for _, flow := range baton.flows {

				if scene.ID == report.SceneID && scene.ID == flow.SceneID {

					globalVariables := make([]*rao.Variable, 0)
					for _, v := range baton.globalVariables {
						globalVariables = append(globalVariables, &rao.Variable{
							Var: v.Var,
							Val: v.Val,
						})
					}

					var nodes rao.Nodes
					if err := bson.Unmarshal(flow.Nodes, &nodes); err != nil {
						proof.Errorf("node bson unmarshal err:%v", err)
						continue
					}

					sceneVariables := make([]*rao.Variable, 0)
					for _, v := range baton.sceneVariables {
						sceneVariables = append(sceneVariables, &rao.Variable{
							Var: v.Var,
							Val: v.Val,
						})
					}

					importVariables := make([]string, 0)
					for _, v := range baton.importVariables {
						importVariables = append(importVariables, v.URL)
					}

					req := rao.Stress{
						PlanID:     baton.plan.ID,
						PlanName:   baton.plan.Name,
						ReportID:   omnibus.DefiniteString(report.ID),
						TeamID:     baton.TeamID,
						ReportName: baton.plan.Name,
						ConfigTask: &rao.ConfigTask{
							TaskType: baton.plan.TaskType,
							Mode:     baton.plan.Mode,
							Remark:   baton.plan.Remark,
							CronExpr: baton.plan.CronExpr,
							ModeConf: &rao.ModeConf{
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
						Scene: &rao.Scene{
							SceneID:                 scene.ID,
							EnablePlanConfiguration: false,
							SceneName:               scene.Name,
							TeamID:                  baton.TeamID,
							Nodes:                   nodes.Nodes,
							Configuration: &rao.SceneConfiguration{
								ParameterizedFile: &rao.SceneVariablePath{
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

func (s *SplitStress) Execute(baton *Baton) error {
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

func (s *SplitImportVariable) Execute(baton *Baton) error {

	reportMemo := make(map[string]int)
	pathMemo := make(map[string]string)
	for _, stress := range baton.stress {
		for _, path := range stress.Scene.Configuration.ParameterizedFile.Path {
			pathMemo[stress.ReportID] = path
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
			return err
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		files := omnibus.Explode("/", p)
		localFilePath := fmt.Sprintf("/tmp/%s", files[len(files)-1])
		if err := ioutil.WriteFile(localFilePath, data, 0644); err != nil {
			return err
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

func (s *RunMachineStress) Execute(baton *Baton) error {
	for _, stress := range baton.stress {
		t := baton.balance.Next()

		tx := query.Use(dal.DB()).ReportMachine
		_, err := tx.WithContext(baton.Ctx).Where(tx.ReportID.Eq(omnibus.DefiniteInt64(stress.ReportID))).Assign(
			tx.ReportID.Value(omnibus.DefiniteInt64(stress.ReportID)),
			tx.IP.Value(omnibus.Explode(":", t)[0]),
		).FirstOrCreate()

		if err != nil {
			return err
		}

		// 增加分区字段判断
		partition := GetPartition()
		if partition == -1 {
			return errors.New("当前没有可用的kafka分区")
		}
		stress.Partition = partition

		_, err = resty.New().R().SetBody(stress).Post(fmt.Sprintf("http://%s/runner/run_plan", t))
		proof.Infof("runner err %+v req %+v", err, proof.Render("req", stress))
		if err != nil {
			return err
		}

		p := query.Use(dal.DB()).Plan
		_, err = p.WithContext(baton.Ctx).Where(p.ID.Eq(baton.PlanID)).UpdateColumn(p.Status, consts.PlanStatusUnderway)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RunMachineStress) SetNext(stress Stress) {
	s.next = stress
}

func GetPartition() int32 {
	//默认分区为0
	var partition int32 = -1 //默认为-1 表示不可用分区锁
	// kafka全局的报告分区key名
	partitionLock := "kafka:report:partition"
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
