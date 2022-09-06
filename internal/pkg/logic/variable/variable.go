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

func ListVariables(ctx context.Context, teamID int64, limit, offset int) ([]*rao.Variable, int64, error) {
	tx := query.Use(dal.DB()).Variable

	v, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return packer.TransVariablesToRaoVariables(v), cnt, nil
}

func DeleteVariable(ctx context.Context, teamID, varID int64) error {
	tx := query.Use(dal.DB()).Variable

	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(varID)).Delete()
	return err
}
