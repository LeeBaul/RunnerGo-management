package report

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-omnibus/omnibus"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/prometheus"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
)

func ListMachines(ctx context.Context, reportID int64) (*rao.ListMachineResp, error) {
	r := query.Use(dal.DB()).Report
	report, err := r.WithContext(ctx).Where(r.ID.Eq(reportID)).First()
	if err != nil {
		return nil, err
	}

	startTimeSec, endTimeSec := report.RanAt.Unix(), time.Now().Unix()
	if report.Status == consts.ReportStatusFinish {
		endTimeSec = report.UpdatedAt.Unix()
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

	for _, machine := range rms {
		cpu, err := prometheus.GetCPURangeUsage(machine.IP, startTimeSec, endTimeSec)
		if err != nil {
			return nil, err
		}

		mem, err := prometheus.GetMemRangeUsage(machine.IP, startTimeSec, endTimeSec)
		if err != nil {
			return nil, err
		}

		net, err := prometheus.GetNetIORangeUsage(machine.IP, startTimeSec, endTimeSec)
		if err != nil {
			return nil, err
		}

		disk, err := prometheus.GetDiskRangeUsage(machine.IP, startTimeSec, endTimeSec)
		if err != nil {
			return nil, err
		}

		for _, c := range cpu {
			c[1], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", omnibus.DefiniteFloat64(c[1])*100), 64)
		}

		for _, m := range mem {
			m[1], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", omnibus.DefiniteFloat64(m[1])*100), 64)
		}

		for _, d := range disk {
			d[1], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", omnibus.DefiniteFloat64(d[1])*100), 64)
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
