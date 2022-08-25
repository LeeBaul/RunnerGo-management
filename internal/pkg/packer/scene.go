package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSceneReqToScene(scene *rao.SaveSceneReq) *mao.Scene {
	request, err := sonic.MarshalString(scene.Request)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("scene.request json marshal err %w", err))
	}

	script, err := sonic.MarshalString(scene.Script)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("scene.script json marshal err %w", err))
	}

	return &mao.Scene{
		TargetID: scene.TargetID,
		Request:  request,
		Script:   script,
	}
}
