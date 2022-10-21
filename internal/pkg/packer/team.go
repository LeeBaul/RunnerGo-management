package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

type TeamMemberCount struct {
	TeamID int64
	Cnt    int64
}

func TransTeamsModelToRaoTeam(teams []*model.Team, userTeams []*model.UserTeam, teamCnt []*TeamMemberCount) []*rao.Team {
	ret := make([]*rao.Team, 0)

	memo := make(map[int64]*model.UserTeam)
	for _, team := range userTeams {
		memo[team.TeamID] = team
	}

	cntMemo := make(map[int64]int64)
	for _, count := range teamCnt {
		cntMemo[count.TeamID] = count.Cnt
	}

	for _, t := range teams {
		ret = append(ret, &rao.Team{
			Name:            t.Name,
			Type:            t.Type,
			Sort:            memo[t.ID].Sort,
			TeamID:          t.ID,
			RoleID:          memo[t.ID].RoleID,
			CreatedUserID:   t.CreatedUserID,
			CreatedUserName: "", // todo user name
			CreatedTimeSec:  t.CreatedAt.Unix(),
			Cnt:             cntMemo[t.ID],
		})
	}

	return ret
}
