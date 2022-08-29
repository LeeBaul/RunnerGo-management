package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveFlowReqToFlow(req *rao.SaveFlowReq) *mao.Flow {
	flowStr, err := sonic.MarshalString(req.Flows)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("folder.script json marshal err %w", err))
	}

	return &mao.Flow{
		SceneID: req.SceneID,
		TeamID:  req.TeamID,
		Version: req.Version,
		Flows:   flowStr,
	}
}
