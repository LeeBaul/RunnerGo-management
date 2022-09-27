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
	for _, t := range teams {
		for _, ut := range userTeams {
			for _, cnt := range teamCnt {
				if ut.TeamID == t.ID && cnt.TeamID == t.ID {
					ret = append(ret, &rao.Team{
						Name:            t.Name,
						Type:            t.Type,
						Sort:            ut.Sort,
						TeamID:          t.ID,
						RoleID:          ut.RoleID,
						CreatedUserID:   t.CreatedUserID,
						CreatedUserName: "", // todo user name
						CreatedTimeSec:  t.CreatedAt.Unix(),
						Cnt:             cnt.Cnt,
					})
				}
			}
		}
	}

	return ret
}
