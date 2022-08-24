package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransPlansToResp(plans []*model.Plan, users []*model.User) []*rao.Plan {
	ret := make([]*rao.Plan, 0)
	for _, p := range plans {
		for _, u := range users {
			if p.RunUserID == u.ID {
				ret = append(ret, &rao.Plan{
					PlanID:         p.ID,
					TeamID:         p.TeamID,
					Name:           p.Name,
					TaskType:       p.TaskType,
					Mode:           p.Mode,
					Status:         p.Status,
					RunUserID:      p.RunUserID,
					RunUserName:    u.Nickname,
					Remark:         p.Remark,
					CreatedTimeSec: p.CreatedAt.Unix(),
					UpdatedTimeSec: p.UpdatedAt.Unix(),
				})
			}
		}
	}

	return ret
}
