package report

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"kp-management/internal/pkg/dal/mao"
	"time"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
)

func ListMachines(ctx context.Context, reportID int64) (*rao.ListMachineResp, error) {
	r := query.Use(dal.DB()).Report
	report, err := r.WithContext(ctx).Where(r.ID.Eq(reportID)).First()
	if err != nil {
		return nil, err
	}

	startTimeSec := report.RanAt.Unix() - 60
	var endTimeSec int64
	// 判断报告是否完成
	if report.Status == consts.ReportStatusNormal { // 进行中
		endTimeSec = time.Now().Unix()
	} else { // 已完成
		endTimeSec = report.UpdatedAt.Unix() + 60
	}

	resp := rao.ListMachineResp{
		StartTimeSec: startTimeSec,
		EndTimeSec:   endTimeSec,
		ReportStatus: report.Status,
		Metrics:      make([]*rao.Metric, 0),
	}

	rm := query.Use(dal.DB()).ReportMachine
	rms, err := rm.WithContext(ctx).Where(rm.ReportID.Eq(reportID)).Find()
	if err != nil {
		return nil, err
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectMachineMonitorData)
	for _, machine := range rms {
		// 从mg里面查出来压力机监控数据
		mmd, err := collection.Find(ctx, bson.D{{"machine_ip", machine.IP}, {"created_at", bson.D{{"$gte", startTimeSec}}}, {"created_at", bson.D{{"$lte", endTimeSec}}}})
		if err != nil {
			return nil, err
		}
		var machineMonitorSlice []*mao.MachineMonitorData
		if err = mmd.All(ctx, &machineMonitorSlice); err != nil {
			return nil, err
		}

		cpu := make([][]interface{}, 0, len(machineMonitorSlice))
		mem := make([][]interface{}, 0, len(machineMonitorSlice))
		net := make([][]interface{}, 0, len(machineMonitorSlice))
		disk := make([][]interface{}, 0, len(machineMonitorSlice))
		for _, machineMonitorInfo := range machineMonitorSlice {
			cpuTmp := make([]interface{}, 0, 2)
			cpuTmp = append(cpuTmp, machineMonitorInfo.MonitorData.CreateTime, machineMonitorInfo.MonitorData.CpuUsage)
			cpu = append(cpu, cpuTmp)

			memTmp := make([]interface{}, 0, 2)
			memTmp = append(memTmp, machineMonitorInfo.MonitorData.CreateTime, machineMonitorInfo.MonitorData.MemInfo[0].UsedPercent)
			mem = append(mem, memTmp)

			//netTmp := make([]interface{}, 0, 2)
			//netTmp[0] = machineMonitorInfo.MonitorData.CreateTime
			//netTmp[1] = machineMonitorInfo.MonitorData.CpuUsage
			//net = append(net, netTmp)

			diskTmp := make([]interface{}, 0, 2)
			diskTmp = append(diskTmp, machineMonitorInfo.MonitorData.CreateTime, machineMonitorInfo.MonitorData.DiskInfos[0].UsedPercent)
			disk = append(disk, diskTmp)

		}
		resp.Metrics = append(resp.Metrics, &rao.Metric{
			CPU:    cpu,
			Mem:    mem,
			NetIO:  net,
			DiskIO: disk,
		})
	}

	return &resp, nil
}
