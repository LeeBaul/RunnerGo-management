package variable

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
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
