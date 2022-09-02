package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveFlowReqToMaoFlow(req *rao.SaveFlowReq) *mao.Flow {

	nodes, err := bson.Marshal(req.Nodes)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("flow.nodes json marshal err %w", err))
	}

	edges, err := bson.Marshal(req.Edges)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("flow.edges json marshal err %w", err))
	}

	return &mao.Flow{
		SceneID:         req.SceneID,
		TeamID:          req.TeamID,
		Version:         req.Version,
		Nodes:           nodes,
		Edges:           edges,
		MultiLevelNodes: req.MultiLevelNodes,
	}
}

func TransMaoFlowToRaoGetFowResp(f *mao.Flow) *rao.GetFlowResp {
	var nodes []*rao.Node
	if err := bson.Unmarshal(f.Nodes, &nodes); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.nodes json unmarshal err %w", err))

	}

	var edges []*rao.Edge
	if err := bson.Unmarshal(f.Edges, &edges); err != nil {
		fmt.Sprintln(fmt.Errorf("flow.edges json unmarshal err %w", err))
	}

	return &rao.GetFlowResp{
		SceneID:         f.SceneID,
		TeamID:          f.TeamID,
		Version:         f.Version,
		Nodes:           nodes,
		Edges:           edges,
		MultiLevelNodes: f.MultiLevelNodes,
	}
}
