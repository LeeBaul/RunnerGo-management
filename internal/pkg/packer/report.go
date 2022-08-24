package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransReportModelToResp(reports []*model.Report, users []*model.User) []*rao.Report {
	ret := make([]*rao.Report, 0)
	for _, r := range reports {
		for _, u := range users {
			if u.ID == r.RunUserID {
				ret = append(ret, &rao.Report{
					ReportID:    r.ID,
					Name:        r.Name,
					Mode:        r.Mode,
					Status:      r.Status,
					RunTimeSec:  r.RanAt.Unix(),
					LastTimeSec: r.UpdatedAt.Unix(),
					RunUserID:   r.RunUserID,
					RunUserName: u.Nickname,
					TeamID:      r.TeamID,
					TaskType:    r.TaskType,
					SceneType:   r.SceneType,
				})
			}
		}
	}
	return ret
}
