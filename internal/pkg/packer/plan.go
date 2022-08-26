package packer

import (
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal/mao"
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

func TransSavePlanReqToModel(req *rao.SavePlanReq, userID int64) *model.Plan {
	return &model.Plan{
		ID:           req.PlanID,
		TeamID:       req.TeamID,
		Name:         req.Name,
		TaskType:     req.TaskType,
		Mode:         req.Mode,
		Status:       consts.PlanStatusNormal,
		CreateUserID: userID,
		Remark:       req.Remark,
	}
}

func TransSavePlanReqToModeConf(req *rao.SavePlanReq) *mao.PlanModeConf {
	mc := req.ModeConf

	return &mao.PlanModeConf{
		PlanID: req.PlanID,
		ModeConf: &mao.ModeConf{
			ReheatTime:       mc.ReheatTime,
			RoundNum:         mc.RoundNum,
			Concurrency:      mc.Concurrency,
			ThresholdValue:   mc.ThresholdValue,
			StartConcurrency: mc.StartConcurrency,
			Step:             mc.Step,
			StepRunTime:      mc.StepRunTime,
			MaxConcurrency:   mc.MaxConcurrency,
			Duration:         mc.Duration,
		},
	}

}
