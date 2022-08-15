package operation

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func List(ctx context.Context, teamID int64, size, limit int32) ([]*rao.Operation, error) {
	tx := query.Use(dal.DB()).Operation
	operations, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransOperationModelToResp(operations), nil
}
