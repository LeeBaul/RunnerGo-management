package plan

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func ListByStatus(ctx context.Context, teamID int64, status int32, limit, offset int) ([]*rao.Plan, int64, error) {
	tx := query.Use(dal.DB()).Plan
	ret, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.Status.Eq(status)).
		Order(tx.UpdatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return packer.TransPlansToResp(ret), cnt, nil
}

func CountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Plan

	return tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
}
