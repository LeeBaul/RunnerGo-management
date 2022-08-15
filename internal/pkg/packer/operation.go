package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransOperationModelToResp(operations []*model.Operation) []*rao.Operation {
	var ret []*rao.Operation
	for _, o := range operations {
		ret = append(ret, &rao.Operation{
			UserID: o.UserID,
			//todo user
			UserName:       "",
			UserAvatar:     "",
			UserStatus:     0,
			Category:       o.Category,
			Name:           o.Name,
			CreatedTimeSec: o.CreatedAt.Unix(),
		})
	}
	return ret
}
