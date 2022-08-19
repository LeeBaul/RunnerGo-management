package api

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.CreateTargetReq, userID int64) error {
	target := packer.TransTargetReqToTarget(req, userID)
	api := packer.TransTargetReqToAPI(req)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			api.TargetID = target.ID
			if err := tx.API.WithContext(ctx).Create(api); err != nil {
				return err
			}

			return nil
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		if _, err := tx.API.WithContext(ctx).Where(tx.API.TargetID.Eq(api.TargetID)).Updates(api); err != nil {
			return err
		}

		return nil
	})
}
