package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransUsersToRaoMembers(users []*model.User, userTeams []*model.UserTeam) []*rao.Member {
	ret := make([]*rao.Member, 0)
	for _, u := range users {
		for _, ut := range userTeams {
			if ut.UserID == u.ID {
				ret = append(ret, &rao.Member{
					Avatar:      "", // todo avatar
					UserID:      u.ID,
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
