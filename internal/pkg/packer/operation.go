package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransOperationsToRaoOperationList(operations []*model.Operation, users []*model.User) []*rao.Operation {
	ret := make([]*rao.Operation, 0)

	memo := make(map[int64]*model.User)
	for _, user := range users {
		memo[user.ID] = user
	}

	for _, o := range operations {
		ret = append(ret, &rao.Operation{
			UserID:         o.UserID,
			UserName:       memo[o.UserID].Nickname,
			UserAvatar:     memo[o.UserID].Avatar,
			UserStatus:     0,
			Category:       o.Category,
			Operate:        o.Operate,
			Name:           o.Name,
			CreatedTimeSec: o.CreatedAt.Unix(),
		})
	}

	return ret
}
