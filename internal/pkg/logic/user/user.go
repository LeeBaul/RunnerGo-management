package user

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
)

func FirstByUserID(ctx context.Context, userID int64) (*rao.Member, error) {
	tx := query.Use(dal.DB()).User
	u, err := tx.WithContext(ctx).Where(tx.ID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}

	return &rao.Member{
		Avatar:   "", // todo avatar
		Email:    u.Email,
		Nickname: u.Nickname,
	}, nil
}
