package group

import (
	"context"
	"fmt"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"

	"go.mongodb.org/mongo-driver/bson"
)

func Save(ctx context.Context, req *rao.SaveGroupReq, userID int64) error {
	target := packer.TransSaveGroupReqToTargetModel(req, userID)
	//group := packer.TransSaveGroupReqToMaoGroup(req)

	//collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectGroup)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {

		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			//group.TargetID = target.ID
			//_, err := collection.InsertOne(ctx, group)
			//if err != nil {
			//	return err
			//}

			return record.InsertCreate(ctx, target.TeamID, userID, fmt.Sprintf("创建分组 - %s", target.Name))
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		//_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": group})
		//if err != nil {
		//	return err
		//}

		return record.InsertUpdate(ctx, target.TeamID, userID, fmt.Sprintf("修改分组 - %s", target.Name))
	})
}

func GetByTargetID(ctx context.Context, teamID, targetID int64) (*rao.Group, error) {
	tx := query.Use(dal.DB()).Target
	t, err := tx.WithContext(ctx).Where(
		tx.ID.Eq(targetID),
		tx.TeamID.Eq(teamID),
		tx.TargetType.Eq(consts.TargetTypeGroup),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).First()

	if err != nil {
		return nil, err
	}

	var g *mao.Group
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectGroup)
	err = collection.FindOne(ctx, bson.D{{"target_id", targetID}}).Decode(&g)
	if err != nil {
		return nil, err
	}

	return packer.TransTargetToRaoGroup(t, g), nil
}
