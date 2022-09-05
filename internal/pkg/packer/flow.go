package packer

import (
	"encoding/json"
	"fmt"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveFlowReqToMaoFlow(req *rao.SaveFlowReq) *mao.Flow {

	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("flow.nodes json marshal err %w", err))
	}

	edges, err := json.Marshal(req.Edges)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("flow.edges json marshal err %w", err))
	}

	return &mao.Flow{
		SceneID:         req.SceneID,
		TeamID:          req.TeamID,
		Version:         req.Version,
		Nodes:           string(nodes),
		Edges:           string(edges),
		MultiLevelNodes: string(req.MultiLevelNodes),
	}
}

func TransMaoFlowToRaoGetFowResp(f *mao.Flow) *rao.GetFlowResp {
	var nodes []*rao.Node
	if err := json.Unmarshal([]byte(f.Nodes), &nodes); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.nodes json unmarshal err %w", err))

	}

	var edges []*rao.Edge
	if err := json.Unmarshal([]byte(f.Edges), &edges); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.edges json unmarshal err %w", err))
	}

	return &rao.GetFlowResp{
		SceneID:         f.SceneID,
		TeamID:          f.TeamID,
		Version:         f.Version,
		Nodes:           nodes,
		Edges:           edges,
		MultiLevelNodes: []byte(f.MultiLevelNodes),
	}
}
