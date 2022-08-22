package report

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func CountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Report

	return tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
}

func ListByTeamID(ctx context.Context, teamID int64, limit, offset int) ([]*rao.Report, int64, error) {
	tx := query.Use(dal.DB()).Report
	reports, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).
		Order(tx.UpdatedAt.Desc(), tx.CreatedAt.Desc()).
		Limit(limit).Offset(offset).Find()

	if err != nil {
		return nil, 0, err
	}

	cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransReportModelToResp(reports), cnt, nil
}
