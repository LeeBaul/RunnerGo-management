package packer

import (
	"encoding/json"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
	"time"
)

func TransPlansToRaoPlanList(plans []*model.Plan, users []*model.User) []*rao.Plan {
	ret := make([]*rao.Plan, 0)

	memo := make(map[int64]*model.User)
	for _, user := range users {
		memo[user.ID] = user
	}

	for _, p := range plans {
		ret = append(ret, &rao.Plan{
			PlanID:            p.ID,
			Rank:              p.Rank,
			TeamID:            p.TeamID,
			Name:              p.Name,
			TaskType:          p.TaskType,
			Mode:              p.Mode,
			Status:            p.Status,
			CreatedUserName:   memo[p.CreateUserID].Nickname,
			CreatedUserAvatar: memo[p.CreateUserID].Avatar,
			CreatedUserID:     p.CreateUserID,
			Remark:            p.Remark,
			CreatedTimeSec:    p.CreatedAt.Unix(),
			UpdatedTimeSec:    p.UpdatedAt.Unix(),
		})
	}

	return ret
}

func TransSavePlanReqToPlanModel(req *rao.SavePlanConfReq, userID int64) *model.Plan {
	return &model.Plan{
		ID:           req.PlanID,
		TeamID:       req.TeamID,
		Name:         req.Name,
		TaskType:     req.TaskType,
		Mode:         req.Mode,
		Status:       consts.PlanStatusNormal,
		CreateUserID: userID,
		Remark:       req.Remark,
		CronExpr:     req.CronExpr,
	}
}

func TransSavePlanReqToMaoTask(req *rao.SavePlanConfReq) *mao.Task {
	mc := req.ModeConf

	return &mao.Task{
		PlanID:   req.PlanID,
		SceneID:  req.SceneID,
		TaskMode: req.Mode,
		TaskType: req.TaskType,
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

func TransTaskToRaoPlan(p *model.Plan, t *mao.Task, u *model.User) *rao.Plan {

	var mc rao.ModeConf
	//var n []*rao.Node
	//var e []*rao.Edge
	if t != nil {
		mc = rao.ModeConf{
			ReheatTime:       t.ModeConf.ReheatTime,
			RoundNum:         t.ModeConf.RoundNum,
			Concurrency:      t.ModeConf.Concurrency,
			ThresholdValue:   t.ModeConf.ThresholdValue,
			StartConcurrency: t.ModeConf.StartConcurrency,
			Step:             t.ModeConf.Step,
			StepRunTime:      t.ModeConf.StepRunTime,
			MaxConcurrency:   t.ModeConf.MaxConcurrency,
			Duration:         t.ModeConf.Duration,
		}

		//var nb mao.Node
		//if err := bson.Unmarshal(t.Nodes, &nb); err != nil {
		//	proof.Errorf("plan.nodes json marshal err %w", err)
		//}
		//n = nb.Nodes
		//
		//var eb mao.Edge
		//if err := bson.Unmarshal(t.Nodes, &eb); err != nil {
		//	proof.Errorf("plan.edges json marshal err %w", err)
		//}
		//e = eb.Edges

	}

	return &rao.Plan{
		PlanID:            p.ID,
		TeamID:            p.TeamID,
		Name:              p.Name,
		TaskType:          p.TaskType,
		Mode:              p.Mode,
		Status:            p.Status,
		CreatedUserID:     p.CreateUserID,
		CreatedUserAvatar: u.Avatar,
		CreatedUserName:   u.Nickname,
		Remark:            p.Remark,
		CreatedTimeSec:    p.CreatedAt.Unix(),
		UpdatedTimeSec:    p.UpdatedAt.Unix(),
		CronExpr:          p.CronExpr,
		ModeConf:          &mc,
		//Nodes:             n,
		//Edges:             e,
	}
}

func TransSaveTimingTaskConfigReqToModelData(req *rao.SavePlanConfReq, userID int64) (*model.TimedTaskConf, error) {
	// 把mode_conf压缩成字符串
	modeConfString, err := json.Marshal(req.ModeConf)
	if err != nil {
		return nil, err
	}
	return &model.TimedTaskConf{
		PlanID:        req.PlanID,
		SenceID:       req.SceneID,
		TeamID:        req.TeamID,
		UserID:        userID,
		Frequency:     req.TimedTaskConf.Frequency,
		TaskExecTime:  req.TimedTaskConf.TaskExecTime,
		TaskCloseTime: req.TimedTaskConf.TaskCloseTime,
		TaskType:      req.TaskType,
		TaskMode:      req.Mode,
		ModeConf:      string(modeConfString),
		Status:        0,
	}, nil
}

func TransChangeReportConfRunToMao(req rao.ChangeTaskConfReq) *mao.ChangeTaskConf {
	return &mao.ChangeTaskConf{
		ReportID: req.ReportID,
		TeamID:   req.TeamID,
		PlanID:   req.PlanID,
		ModeConf: &mao.ModeConf{
			ReheatTime:       req.ModeConf.ReheatTime,
			RoundNum:         req.ModeConf.RoundNum,
			Concurrency:      req.ModeConf.Concurrency,
			ThresholdValue:   req.ModeConf.ThresholdValue,
			StartConcurrency: req.ModeConf.StartConcurrency,
			Step:             req.ModeConf.Step,
			StepRunTime:      req.ModeConf.StepRunTime,
			MaxConcurrency:   req.ModeConf.MaxConcurrency,
			Duration:         req.ModeConf.Duration,
			CreatedTimeSec:   time.Now().Unix(),
		},
	}
}
