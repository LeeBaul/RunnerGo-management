package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/proof"
	"github.com/goccy/go-json"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/machine"
	"kp-management/internal/pkg/logic/stress"
	"kp-management/internal/pkg/packer"
	"strconv"
	"strings"
	"time"
)

// GetMachineList 获取机器列表
func GetMachineList(ctx *gin.Context) {
	var req rao.GetMachineListParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	res, total, err := machine.GetMachineList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.GetMachineListResponse{
		MachineList: res,
		Total:       total,
	})
	return
}

// ChangeMachineOnOff 压力机启用或卸载
func ChangeMachineOnOff(ctx *gin.Context) {
	var req rao.ChangeMachineOnOff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := machine.ChangeMachineOnOff(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}

// MachineDataInsert 把压力机上报的机器数据插入数据库
func MachineDataInsert() {
	ctx := context.Background()
	for {
		// 从Redis获取压力机列表
		machineListRes := dal.RDB.HGetAll(consts.MachineListRedisKey)
		if len(machineListRes.Val()) == 0 || machineListRes.Err() != nil {
			proof.Errorf("压力机数据入库--没有获取到任何压力机上报数据，err:", machineListRes.Err())
			time.Sleep(5 * time.Second) // 5秒循环一次
			continue
		}

		// 有数据，则入库
		for machineAddr, machineDetail := range machineListRes.Val() {
			// 获取机器IP，端口号，区域
			machineAddrSlice := strings.Split(machineAddr, "_")
			if len(machineAddrSlice) != 3 {
				continue
			}

			// 把机器详情信息解析成格式化数据
			var runnerMachineInfo stress.HeartBeat
			err := json.Unmarshal([]byte(machineDetail), &runnerMachineInfo)
			if err != nil {
				proof.Infof("压力机数据入库--压力机详情数据解析失败，err：", err)
				continue
			}

			ip := machineAddrSlice[0]
			port := machineAddrSlice[1]
			portInt, err := strconv.Atoi(port)
			if err != nil {
				proof.Errorf("压力机数据入库--转换类型失败，err:", err)
				continue
			}
			region := machineAddrSlice[2]

			// 查询当前机器信息是否存在数据库
			tx := dal.GetQuery().Machine

			// 查询数据
			_, err = tx.WithContext(ctx).Where(tx.IP.Eq(ip)).First()
			if err != nil && err != gorm.ErrRecordNotFound {
				proof.Errorf("压力机数据入库--查询数据出错，err:", err)
				continue
			}

			if err == nil { // 查到了，修改数据
				updateData := model.Machine{
					Port:              int32(portInt),
					Region:            region,
					Name:              runnerMachineInfo.Name,
					CPUUsage:          float32(runnerMachineInfo.CpuUsage),
					CPULoadOne:        float32(runnerMachineInfo.CpuLoad.Load1),
					CPULoadFive:       float32(runnerMachineInfo.CpuLoad.Load5),
					CPULoadFifteen:    float32(runnerMachineInfo.CpuLoad.Load15),
					MemUsage:          float32(runnerMachineInfo.MemInfo[0].UsedPercent),
					DiskUsage:         float32(runnerMachineInfo.DiskInfos[0].UsedPercent),
					MaxGoroutines:     runnerMachineInfo.MaxGoroutines,
					CurrentGoroutines: runnerMachineInfo.CurrentGoroutines,
					ServerType:        int32(runnerMachineInfo.ServerType),
				}
				_, err := tx.WithContext(ctx).Where(tx.IP.Eq(ip)).Updates(&updateData)
				if err != nil {
					proof.Errorf("压力机数据入库--更新数据失败，err:", err)
					continue
				}
			} else { // 没查到，新增数据
				insertData := model.Machine{
					IP:                ip,
					Port:              int32(portInt),
					Region:            region,
					Name:              runnerMachineInfo.Name,
					CPUUsage:          float32(runnerMachineInfo.CpuUsage),
					CPULoadOne:        float32(runnerMachineInfo.CpuLoad.Load1),
					CPULoadFive:       float32(runnerMachineInfo.CpuLoad.Load5),
					CPULoadFifteen:    float32(runnerMachineInfo.CpuLoad.Load15),
					MemUsage:          float32(runnerMachineInfo.MemInfo[0].UsedPercent),
					DiskUsage:         float32(runnerMachineInfo.DiskInfos[0].UsedPercent),
					MaxGoroutines:     runnerMachineInfo.MaxGoroutines,
					CurrentGoroutines: runnerMachineInfo.CurrentGoroutines,
					ServerType:        int32(runnerMachineInfo.ServerType),
				}
				err := tx.WithContext(ctx).Create(&insertData)
				if err != nil {
					proof.Errorf("压力机数据入库")
					continue
				}
			}
		}

		time.Sleep(5 * time.Second) // 5秒循环一次
	}

}

// MachineMonitorInsert 压力机监控数据入库
func MachineMonitorInsert() {
	ctx := context.Background()
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectMachineMonitorData)
	for {
		machineList, _ := dal.GetRDB().Keys(ctx, consts.MachineMonitorPrefix+"*").Result()

		for _, MachineMonitorKey := range machineList {
			machineAddrSlice := strings.Split(MachineMonitorKey, ":")
			if len(machineAddrSlice) != 2 {
				continue
			}
			machineIP := machineAddrSlice[1]
			// 从Redis获取压力机列表
			machineListRes := dal.RDB.LRange(MachineMonitorKey, 0, -1).Val()
			if len(machineListRes) == 0 {
				continue
			}
			for _, monitorData := range machineListRes {
				var runnerMachineInfo mao.HeartBeat
				// 把机器详情信息解析成格式化数据
				err := json.Unmarshal([]byte(monitorData), &runnerMachineInfo)
				if err != nil {
					proof.Infof("压力机监控数据入库--数据解析失败 err：", err)
					continue
				}

				machineMonitorInsertData := packer.TransMachineMonitorToMao(machineIP, runnerMachineInfo, runnerMachineInfo.CreateTime)
				_, err = collection.InsertOne(ctx, machineMonitorInsertData)
				if err != nil {
					proof.Infof("压力机监控数据入库--插入mg数据失败，err:", err)
					continue
				}
			}
			// 数据入库完毕，把redis列表删掉
			err := dal.GetRDB().Del(ctx, MachineMonitorKey)
			if err.Err() != nil {
				proof.Errorf("压力机监控数据入库--删除redis列表失败，err:", err.Err())
			}
		}
		time.Sleep(5 * time.Second)
	}
}
