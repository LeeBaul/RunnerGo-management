package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransReportModelToResp(reports []*model.Report) []*rao.Report {
	ret := make([]*rao.Report, 0)
	for _, r := range reports {
		ret = append(ret, &rao.Report{
			ReportID: r.ID,
			Name:     r.Name,
			Mode:     r.Mode,
			Status:   r.Status,
		})
	}
	return ret
}
