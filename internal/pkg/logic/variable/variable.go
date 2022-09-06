package variable

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SaveVariable(ctx context.Context, req *rao.SaveVariableReq) error {
	tx := query.Use(dal.DB()).Variable

	_, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.VarID)).Assign(
		tx.TeamID.Value(req.TeamID),
		tx.Var.Value(req.Var),
		tx.Val.Value(req.Val),
		tx.Description.Value(req.Description),
	).FirstOrCreate()

	return err
}

func ListVariables(ctx context.Context, teamID int64) ([]*rao.Variable, error) {
	tx := query.Use(dal.DB()).Variable

	v, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransVariablesToRaoVariables(v), nil
}

func DeleteVariable(ctx context.Context, teamID, varID int64) error {
	tx := query.Use(dal.DB()).Variable

	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(varID)).Delete()
	return err
}
