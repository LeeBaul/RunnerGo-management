package scene

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

func Save(ctx context.Context, req *rao.SaveSceneReq, userID int64) error {
	target := packer.TransSceneReqToTarget(req, userID)
	scene := packer.TransSceneReqToScene(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectScene)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			scene.TargetID = target.ID

			_, err := collection.InsertOne(ctx, scene)

			tx.Operation.WithContext(ctx).Create(&model.Operation{
				TeamID:   target.TeamID,
				UserID:   userID,
				Category: consts.OperationCategoryCreate,
				Name:     fmt.Sprintf("创建场景 - %s", target.Name),
			})

			return err
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": scene})

		tx.Operation.WithContext(ctx).Create(&model.Operation{
			TeamID:   target.TeamID,
			UserID:   userID,
			Category: consts.OperationCategoryCreate,
			Name:     fmt.Sprintf("修改场景 - %s", target.Name),
		})

		return err
	})
}
