package target

import (
	"context"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
)

func APICountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq(consts.TargetTypeAPI),
		tx.TeamID.Eq(teamID),
	).Count()
}

func SceneCountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq(consts.TargetTypeScene),
		tx.TeamID.Eq(teamID),
	).Count()
}
