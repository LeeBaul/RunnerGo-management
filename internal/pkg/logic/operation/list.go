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

	var userIDs []int64
	for _, o := range operations {
		userIDs = append(userIDs, o.UserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransOperationModelToResp(operations, users), nil
}
