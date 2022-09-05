package packer

import (
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveFolderReqToTargetModel(folder *rao.SaveFolderReq, userID int64) *model.Target {
	return &model.Target{
		ID:            folder.TargetID,
		TeamID:        folder.TeamID,
		TargetType:    consts.TargetTypeFolder,
		Name:          folder.Name,
		ParentID:      folder.ParentID,
		Method:        folder.Method,
		Sort:          folder.Sort,
		TypeSort:      folder.TypeSort,
		Status:        consts.TargetStatusNormal,
		Version:       folder.Version,
		CreatedUserID: userID,
		RecentUserID:  userID,
		Source:        consts.TargetSourceNormal,
	}
}

func TransSaveTargetReqToTargetModel(target *rao.SaveTargetReq, userID int64) *model.Target {
	return &model.Target{
		ID:            target.TargetID,
		TeamID:        target.TeamID,
		TargetType:    consts.TargetTypeAPI,
		Name:          target.Name,
		ParentID:      target.ParentID,
		Method:        target.Method,
		Sort:          target.Sort,
		TypeSort:      target.TypeSort,
		Status:        consts.TargetStatusNormal,
		Version:       target.Version,
		CreatedUserID: userID,
		RecentUserID:  userID,
		Source:        consts.TargetSourceNormal,
	}
}

func TransSaveGroupReqToTargetModel(group *rao.SaveGroupReq, userID int64) *model.Target {
	return &model.Target{
		ID:            group.TargetID,
		TeamID:        group.TeamID,
		TargetType:    consts.TargetTypeGroup,
		Name:          group.Name,
		ParentID:      group.ParentID,
		Method:        group.Method,
		Sort:          group.Sort,
		TypeSort:      group.TypeSort,
		Status:        consts.TargetStatusNormal,
		Version:       group.Version,
		CreatedUserID: userID,
		RecentUserID:  userID,
		Source:        group.Source,
	}
}

func TransSaveSceneReqToTargetModel(scene *rao.SaveSceneReq, userID int64) *model.Target {
	return &model.Target{
		ID:            scene.TargetID,
		TeamID:        scene.TeamID,
		TargetType:    consts.TargetTypeScene,
		Name:          scene.Name,
		ParentID:      scene.ParentID,
		Method:        scene.Method,
		Sort:          scene.Sort,
		TypeSort:      scene.TypeSort,
		Status:        consts.TargetStatusNormal,
		Version:       scene.Version,
		CreatedUserID: userID,
		RecentUserID:  userID,
		Source:        scene.Source,
		PlanID:        scene.PlanID,
	}
}
