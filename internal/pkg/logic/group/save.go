package group

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.SaveGroupReq) error {
	target := packer.TransGroupReqToTarget(req)
	group := packer.TransGroupReqToGroup(req)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if err := tx.Target.WithContext(ctx).Save(target); err != nil {
			return err
		}

		if err := tx.Group.WithContext(ctx).Where(tx.Group.TargetID.Eq(group.TargetID)).Save(group); err != nil {
			return err
		}

		return nil
	})
}
