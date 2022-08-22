package api

import (
	"context"
	"kp-management/internal/pkg/biz/consts"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.CreateTargetReq, userID int64) error {
	target := packer.TransTargetReqToTarget(req, userID)
	api := packer.TransTargetReqToAPI(req)

	collection := dal.GetMongo().Database(dal.MongoD()).Collection(consts.CollectAPI)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			api.TargetID = target.ID
			_, err := collection.InsertOne(ctx, api)

			return err
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": api})

		return err
	})
}
