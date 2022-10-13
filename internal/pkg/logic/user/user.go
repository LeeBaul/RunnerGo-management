package user

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
)

func FirstByUserID(ctx context.Context, userID, teamID int64) (*rao.Member, error) {
	tx := query.Use(dal.DB()).User
	u, err := tx.WithContext(ctx).Where(tx.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	ut := dal.GetQuery().UserTeam
	r, err := ut.WithContext(ctx).Where(ut.UserID.Eq(u.ID), ut.TeamID.Eq(teamID)).First()
	if err != nil {
		return nil, err
	}

	return &rao.Member{
		Avatar:   u.Avatar,
		Email:    u.Email,
		Nickname: u.Nickname,
		UserID:   userID,
		RoleID:   r.RoleID,
	}, nil
}
