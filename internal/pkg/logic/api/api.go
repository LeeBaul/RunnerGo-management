package api

import (
	"context"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SaveAPI(ctx context.Context, req *rao.CreateTargetReq) error {
	target := packer.TransTargetReqToTarget(req)
	api := packer.TransTargetReqToAPI(req)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if err := tx.Target.WithContext(ctx).Save(target); err != nil {
			return err
		}

		if err := tx.API.WithContext(ctx).Where(tx.API.TargetID.Eq(req.TargetID)).Save(api); err != nil {
			return err
		}

		return nil
	})
}
