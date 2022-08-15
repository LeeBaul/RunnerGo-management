package folder

import (
	"context"
	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
)

func Delete(ctx context.Context, targetID int64) error {
	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.Target.WithContext(ctx).Where(tx.Target.ID.Eq(targetID), tx.Target.TargetType.Eq(consts.TargetTypeFolder)).Delete(); err != nil {
			return err
		}

		if _, err := tx.Folder.WithContext(ctx).Where(tx.Folder.TargetID.Eq(targetID)).Delete(); err != nil {
			return err
		}

		return nil
	})
}
