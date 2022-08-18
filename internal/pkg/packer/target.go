package packer

import (
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransFolderReqToTarget(folder *rao.SaveFolderReq, userID int64) *model.Target {

	return &model.Target{
		ID:            folder.TargetID,
		TeamID:        folder.TeamID,
		TargetType:    consts.TargetTypeFolder,
		Name:          folder.Name,
		ParentID:      folder.ParentID,
		Method:        folder.Method,
		Sort:          folder.Sort,
		TypeSort:      folder.TypeSort,
		Status:        1,
		Version:       folder.Version,
		CreatedUserID: userID,
		RecentUserID:  userID,
	}
}

func TransTargetReqToTarget(target *rao.CreateTargetReq) *model.Target {
	return &model.Target{
		ID:            target.TargetID,
		TargetType:    target.TargetType,
		Name:          target.Name,
		ParentID:      target.ParentID,
		Method:        target.Method,
		Sort:          target.Sort,
		TypeSort:      target.TypeSort,
		Status:        1,
		Version:       target.Version,
		CreatedUserID: 0,
		RecentUserID:  0,
		// todo user_id
	}
}

func TransGroupReqToTarget(group *rao.SaveGroupReq) *model.Target {
	return &model.Target{
		ID:            group.TargetID,
		TargetType:    group.TargetType,
		Name:          group.Name,
		ParentID:      group.ParentID,
		Method:        group.Method,
		Sort:          group.Sort,
		TypeSort:      group.TypeSort,
		Status:        1,
		Version:       group.Version,
		CreatedUserID: 0,
		RecentUserID:  0,
		// todo user_id
	}
}
