package report

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gen"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func CountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Report

	return tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
}

func ListUnderway(ctx context.Context, teamID int64, limit, offset int) ([]*rao.Report, int64, error) {
	tx := query.Use(dal.DB()).Report

	reports, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.Status.Eq(consts.ReportStatusNormal)).
		Order(tx.UpdatedAt.Desc(), tx.CreatedAt.Desc()).
		FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range reports {
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	var planIDs []int64
	var sceneIDs []int64
	for _, report := range reports {
		planIDs = append(planIDs, report.PlanID)
		sceneIDs = append(sceneIDs, report.SceneID)
	}

	p := dal.GetQuery().Plan
	plans, err := p.WithContext(ctx).Where(p.ID.In(planIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	s := dal.GetQuery().Target
	scenes, err := s.WithContext(ctx).Where(s.ID.In(sceneIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransReportModelToRaoReportList(reports, users, plans, scenes), cnt, nil
}

func ListByTeamID(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64) ([]*rao.Report, int64, error) {
	tx := query.Use(dal.DB()).Report

	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.TeamID.Eq(teamID))

	if keyword != "" {
		p := dal.GetQuery().Plan
		plans, err := p.WithContext(ctx).Where(p.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}
		var planIDs []int64
		for _, plan := range plans {
			planIDs = append(planIDs, plan.ID)
		}
		conditions = append(conditions, tx.PlanID.In(planIDs...))

		s := dal.GetQuery().Target
		scenes, err := s.WithContext(ctx).Where(s.Name.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}
		var sceneIDs []int64
		for _, scene := range scenes {
			sceneIDs = append(sceneIDs, scene.ID)
		}
		if len(sceneIDs) > 0 {
			conditions[1] = tx.SceneID.In(sceneIDs...)
		}

		u := query.Use(dal.DB()).User
		users, err := u.WithContext(ctx).Where(u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}

		if len(users) > 0 {
			conditions[1] = tx.RunUserID.Eq(users[0].ID)
		}
	}

	if startTimeSec > 0 && endTimeSec > 0 {
		startTime := time.Unix(startTimeSec, 0)
		endTime := time.Unix(endTimeSec, 0)
		conditions = append(conditions, tx.CreatedAt.Between(startTime, endTime))
	}

	reports, cnt, err := tx.WithContext(ctx).Where(conditions...).
		Order(tx.UpdatedAt.Desc(), tx.CreatedAt.Desc()).
		FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range reports {
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	var planIDs []int64
	var sceneIDs []int64
	for _, report := range reports {
		planIDs = append(planIDs, report.PlanID)
		sceneIDs = append(sceneIDs, report.SceneID)
	}

	p := dal.GetQuery().Plan
	plans, err := p.WithContext(ctx).Where(p.ID.In(planIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	s := dal.GetQuery().Target
	scenes, err := s.WithContext(ctx).Where(s.ID.In(sceneIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransReportModelToRaoReportList(reports, users, plans, scenes), cnt, nil
}

func DeleteReport(ctx context.Context, teamID, reportID int64) error {
	tx := query.Use(dal.DB()).Report
	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(reportID)).Delete()

	return err
}
