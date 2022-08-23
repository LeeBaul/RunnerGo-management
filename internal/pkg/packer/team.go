package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransTeamsModelToResp(teams []*model.Team, userTeams []*model.UserTeam) []*rao.Team {
	ret := make([]*rao.Team, 0)
	for _, t := range teams {
		for _, ut := range userTeams {
			if ut.TeamID == t.ID {
				ret = append(ret, &rao.Team{
					Name:   t.Name,
					Sort:   ut.Sort,
					TeamID: t.ID,
					RoleID: ut.RoleID,
				})
			}
		}
	}

	return ret
}

func TransUsersModelToMembers(users []*model.User, userTeams []*model.UserTeam) []*rao.Member {
	ret := make([]*rao.Member, 0)
	for _, u := range users {
		for _, ut := range userTeams {
			if ut.UserID == u.ID {
				ret = append(ret, &rao.Member{
					Avatar:      "", // todo avatar
					Email:       u.Email,
					Nickname:    u.Nickname,
					JoinTimeSec: ut.CreatedAt.Unix(),
					RoleID:      ut.RoleID,
				})
			}
		}
	}
	return ret
}
