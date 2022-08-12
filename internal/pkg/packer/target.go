package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToTarget(folder *rao.SaveFolderReq) *model.Target {
	return &model.Target{
		ID:            folder.TargetID,
		TargetType:    folder.TargetType,
		Name:          folder.Name,
		ParentID:      folder.ParentID,
		Method:        folder.Method,
		Sort:          folder.Sort,
		TypeSort:      folder.TypeSort,
		Status:        1,
		Version:       folder.Version,
		CreatedUserID: 0,
		RecentUserID:  0,
		// todo user_id
	}
}
