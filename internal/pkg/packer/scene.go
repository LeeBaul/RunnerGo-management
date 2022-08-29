package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSceneReqToScene(scene *rao.SaveSceneReq) *mao.Scene {
	request, err := bson.Marshal(scene.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("scene.request json marshal err %w", err))
	}

	script, err := bson.Marshal(scene.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("scene.script json marshal err %w", err))
	}

	return &mao.Scene{
		TargetID: scene.TargetID,
		Request:  request,
		Script:   script,
	}
}
