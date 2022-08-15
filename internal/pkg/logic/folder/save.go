package folder

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.SaveFolderReq) error {
	target := packer.TransFolderReqToTarget(req)
	folder := packer.TransFolderReqToFolder(req)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if err := tx.Target.WithContext(ctx).Save(target); err != nil {
			return err
		}

		if err := tx.Folder.WithContext(ctx).Where(tx.Folder.TargetID.Eq(folder.TargetID)).Save(folder); err != nil {
			return err
		}

		return nil
	})
}
