package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransReportModelToRaoReportList(reports []*model.Report, users []*model.User) []*rao.Report {
	ret := make([]*rao.Report, 0)

	memo := make(map[int64]*model.User)
	for _, user := range users {
		memo[user.ID] = user
	}

	for _, r := range reports {
		ret = append(ret, &rao.Report{
			ReportID:    r.ID,
			Rank:        r.Rank,
			TaskType:    r.TaskType,
			TaskMode:    r.TaskMode,
			Status:      r.Status,
			RunTimeSec:  r.RanAt.Unix(),
			LastTimeSec: r.UpdatedAt.Unix(),
			RunUserID:   r.RunUserID,
			RunUserName: memo[r.RunUserID].Nickname,
			TeamID:      r.TeamID,
			PlanID:      r.PlanID,
			PlanName:    r.PlanName,
			SceneID:     r.SceneID,
			SceneName:   r.SceneName,
		})
	}
	return ret
}
