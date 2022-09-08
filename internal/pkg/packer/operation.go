package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransOperationsToRaoOperationList(operations []*model.Operation, users []*model.User) []*rao.Operation {
	ret := make([]*rao.Operation, 0)
	for _, o := range operations {

		for _, u := range users {
			if u.ID == o.UserID {
				ret = append(ret, &rao.Operation{
					UserID:         o.UserID,
					UserName:       u.Nickname,
					UserAvatar:     u.Avatar,
					UserStatus:     0,
					Category:       o.Category,
					Name:           o.Name,
					CreatedTimeSec: o.CreatedAt.Unix(),
				})
			}
		}

	}
	return ret
}
