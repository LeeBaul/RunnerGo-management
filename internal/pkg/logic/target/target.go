package target

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
)

func APICount(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq("api"),
		tx.TeamID.Eq(teamID),
	).Count()
}
