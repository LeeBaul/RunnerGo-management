package target

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/query"

	"go.mongodb.org/mongo-driver/bson"
)

func Trash(ctx context.Context, targetID int64) error {
	t := query.Use(dal.DB()).Target
	_, err := t.WithContext(ctx).Where(t.ID.Eq(targetID)).UpdateColumn(t.Status, consts.TargetStatusTrash)

	return err
}

func Delete(ctx context.Context, targetID int64) error {
	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.Target.WithContext(ctx).Where(tx.Target.ID.Eq(targetID)).Delete(); err != nil {
			return err
		}

		filter := bson.D{{"target_id", targetID}}

		if _, err := dal.GetMongo().Database(dal.MongoD()).Collection(consts.CollectFolder).DeleteOne(ctx, filter); err != nil {
			return err
		}

		if _, err := dal.GetMongo().Database(dal.MongoD()).Collection(consts.CollectAPI).DeleteOne(ctx, filter); err != nil {
			return err
		}

		//if _, err := tx.Folder.WithContext(ctx).Where(tx.Folder.TargetID.Eq(targetID)).Delete(); err != nil {
		//	return err
		//}

		//if _, err := tx.API.WithContext(ctx).Where(tx.API.TargetID.Eq(targetID)).Delete(); err != nil {
		//	return err
		//}

		return nil
	})
}
