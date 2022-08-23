package folder

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, userID int64, req *rao.SaveFolderReq) error {
	target := packer.TransFolderReqToTarget(req, userID)
	folder := packer.TransFolderReqToFolder(req)

	collection := dal.GetMongo().Database(dal.MongoD()).Collection(consts.CollectFolder)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			folder.TargetID = target.ID
			_, err := collection.InsertOne(ctx, folder)

			tx.Operation.WithContext(ctx).Create(&model.Operation{
				TeamID:   target.TeamID,
				UserID:   userID,
				Category: consts.OperationCategoryCreate,
				Name:     fmt.Sprintf("创建文件夹 - %s", target.Name),
			})

			return err
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": folder})

		tx.Operation.WithContext(ctx).Create(&model.Operation{
			TeamID:   target.TeamID,
			UserID:   userID,
			Category: consts.OperationCategoryUpdate,
			Name:     fmt.Sprintf("修改文件夹 - %s", target.Name),
		})

		return err
	})
}
