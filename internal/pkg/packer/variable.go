package packer

import (
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransModelVariablesToRaoVariables(vs []*model.Variable) []*rao.Variable {
	ret := make([]*rao.Variable, 0)
	for _, v := range vs {
		ret = append(ret, &rao.Variable{
			VarID:       v.ID,
			TeamID:      v.TeamID,
			Var:         v.Var,
			Val:         v.Val,
			Description: v.Description,
		})
	}
	return ret
}

func TransRaoVariablesToModelVariables(teamID int64, vs []*rao.Variable) []*model.Variable {
	ret := make([]*model.Variable, 0)
	for _, v := range vs {
		ret = append(ret, &model.Variable{
			TeamID:      teamID,
			Var:         v.Var,
			Val:         v.Val,
			Description: v.Description,
			Type:        consts.VariableTypeGlobal,
		})
	}
	return ret
}

func TransSceneRaoVariablesToModelVariables(teamID, sceneID int64, vs []*rao.Variable) []*model.Variable {
	ret := make([]*model.Variable, 0)
	for _, v := range vs {
		ret = append(ret, &model.Variable{
			TeamID:      teamID,
			SceneID:     sceneID,
			Var:         v.Var,
			Val:         v.Val,
			Description: v.Description,
			Type:        consts.VariableTypeScene,
		})
	}
	return ret
}
