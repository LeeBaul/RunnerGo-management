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
		fmt.Sprintln(fmt.Errorf("flow.flows json marshal err %w", err))
	}

	return &mao.Flow{
		SceneID: req.SceneID,
		TeamID:  req.TeamID,
		Version: req.Version,
		Flows:   flowStr,
	}
}

func TransMongoFlowToResp(f *mao.Flow) *rao.GetFlowResp {
	var flows []*rao.Flow
	if err := sonic.UnmarshalString(f.Flows, &flows); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.flows json unmarshal err %w", err))
	}

	return &rao.GetFlowResp{
		SceneID: f.SceneID,
		TeamID:  f.TeamID,
		Version: f.Version,
		Flows:   flows,
	}
}
