package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransReportModelToRaoReportList(reports []*model.Report, users []*model.User, plans []*model.Plan, scenes []*model.Target) []*rao.Report {
	ret := make([]*rao.Report, 0)
	for _, r := range reports {
		for _, u := range users {
			for _, p := range plans {
				for _, s := range scenes {

					if u.ID == r.RunUserID && p.ID == r.PlanID && s.ID == r.SceneID {
						ret = append(ret, &rao.Report{
							ReportID:    r.ID,
							TaskType:    r.TaskType,
							TaskMode:    r.TaskMode,
							Status:      r.Status,
							RunTimeSec:  r.RanAt.Unix(),
							LastTimeSec: r.UpdatedAt.Unix(),
							RunUserID:   r.RunUserID,
							RunUserName: u.Nickname,
							TeamID:      r.TeamID,
							PlanID:      r.PlanID,
							PlanName:    p.Name,
							SceneID:     r.SceneID,
							SceneName:   s.Name,
						})
					}

				}
			}
		}
	}
	return ret
}
