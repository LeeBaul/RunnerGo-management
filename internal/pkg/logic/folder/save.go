package folder

import (
	"context"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, userID int64, req *rao.SaveFolderReq) error {
	target := packer.TransFolderReqToTarget(req, userID)
	folder := packer.TransFolderReqToFolder(req)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			folder.TargetID = target.ID
			if err := tx.Folder.WithContext(ctx).Create(folder); err != nil {
				return err
			}

			return nil
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		if _, err := tx.Folder.WithContext(ctx).Where(tx.Folder.TargetID.Eq(folder.ID)).Updates(folder); err != nil {
			return err
		}

		return nil
	})
}
