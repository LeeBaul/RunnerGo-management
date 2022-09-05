package packer

import (
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveSceneReqToMaoScene(scene *rao.SaveSceneReq) *mao.Scene {
	//request, err := bson.Marshal(scene.Request)
	//if err != nil {
	//	fmt.Sprintln(fmt.Errorf("scene.request json marshal err %w", err))
	//}
	//
	//script, err := bson.Marshal(scene.Script)
	//if err != nil {
	//	fmt.Sprintln(fmt.Errorf("scene.script json marshal err %w", err))
	//}

	return &mao.Scene{
		TargetID: scene.TargetID,
		//Request:  request,
		//Script:   script,
	}
}

func TransTargetToRaoScene(targets []*model.Target, scenes []*mao.Scene) []*rao.Scene {

	ret := make([]*rao.Scene, 0)
	for _, t := range targets {
		//for _, scene := range scenes {
		//	if t.ID == scene.TargetID {
		//		var r rao.Request
		//		if err := bson.Unmarshal(scene.Request, &r); err != nil {
		//			fmt.Sprintln(fmt.Errorf("scene.request json UnMarshal err %w", err))
		//		}
		//
		//		var s rao.Script
		//		if err := bson.Unmarshal(scene.Script, &s); err != nil {
		//			fmt.Sprintln(fmt.Errorf("scene.script json UnMarshal err %w", err))
		//		}

		ret = append(ret, &rao.Scene{
			TeamID:   t.TeamID,
			TargetID: t.ID,
			ParentID: t.ParentID,
			Name:     t.Name,
			Method:   t.Method,
			Sort:     t.Sort,
			TypeSort: t.TypeSort,
			Version:  t.Version,
			//Request:  &r,
			//Script:   &s,
		})

		//	}
		//}
	}

	return ret
}
