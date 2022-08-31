package packer

import (
	"encoding/json"
	"fmt"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveFlowReqToMaoFlow(req *rao.SaveFlowReq) *mao.Flow {
	flowByte, err := json.Marshal(req.Flows)
	if err != nil {
		fmt.Errorf("flow.flows json marshal err %w", err)
	}

	return &mao.Flow{
		SceneID: req.SceneID,
		TeamID:  req.TeamID,
		Version: req.Version,
		Flows:   string(flowByte),
	}
}

func TransMaoFlowToRaoGetFowResp(f *mao.Flow) *rao.GetFlowResp {
	var flows []*rao.Flow
	if err := json.Unmarshal([]byte(f.Flows), &flows); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.flows json unmarshal err %w", err))
	}

	return &rao.GetFlowResp{
		SceneID: f.SceneID,
		TeamID:  f.TeamID,
		Version: f.Version,
		Flows:   flows,
	}
}
