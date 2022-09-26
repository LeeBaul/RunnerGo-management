package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransReportModelToRaoReportList(reports []*model.Report, users []*model.User) []*rao.Report {
	ret := make([]*rao.Report, 0)
	for _, r := range reports {
		for _, u := range users {

			if u.ID == r.RunUserID {
				ret = append(ret, &rao.Report{
					ReportID:    r.ID,
					Rank:        r.Rank,
					TaskType:    r.TaskType,
					TaskMode:    r.TaskMode,
					Status:      r.Status,
					RunTimeSec:  r.RanAt.Unix(),
					LastTimeSec: r.UpdatedAt.Unix(),
					RunUserID:   r.RunUserID,
					RunUserName: u.Nickname,
					TeamID:      r.TeamID,
					PlanID:      r.PlanID,
					PlanName:    r.PlanName,
					SceneID:     r.SceneID,
					SceneName:   r.SceneName,
				})
			}

		}
	}
	return ret
}
