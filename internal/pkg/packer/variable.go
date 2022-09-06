package packer

import (
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransVariablesToRaoVariables(vs []*model.Variable) []*rao.Variable {
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
