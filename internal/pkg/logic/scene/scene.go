package scene

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.SaveSceneReq, userID int64) error {
	target := packer.TransSaveSceneReqToTargetModel(req, userID)
	scene := packer.TransSaveSceneReqToMaoScene(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectScene)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			scene.TargetID = target.ID
			_, err := collection.InsertOne(ctx, scene)

			record.InsertCreate(ctx, target.TeamID, userID, fmt.Sprintf("创建场景 - %s", target.Name))

			return err
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": scene})

		record.InsertUpdate(ctx, target.TeamID, userID, fmt.Sprintf("修改场景 - %s", target.Name))

		return err
	})
}

func GetByTargetID(ctx context.Context, teamID, targetID int64) (*rao.Scene, error) {
	tx := query.Use(dal.DB()).Target
	t, err := tx.WithContext(ctx).Where(
		tx.ID.Eq(targetID),
		tx.TeamID.Eq(teamID),
		tx.TargetType.Eq(consts.TargetTypeScene),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).First()

	if err != nil {
		return nil, err
	}

	var s *mao.Scene
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectScene)
	err = collection.FindOne(ctx, bson.D{{"target_id", targetID}}).Decode(&s)
	if err != nil {
		return nil, err
	}

	return packer.TransTargetToRaoScene(t, s), nil
}
