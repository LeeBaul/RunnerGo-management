package api

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.SaveTargetReq, userID int64) (int64, error) {
	target := packer.TransSaveTargetReqToTargetModel(req, userID)
	api := packer.TransSaveTargetReqToMaoAPI(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPI)

	err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			api.TargetID = target.ID
			_, err := collection.InsertOne(ctx, api)
			if err != nil {
				return err
			}

			return record.InsertCreate(ctx, target.TeamID, userID, record.OperationOperateCreateAPI, target.Name)
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": api})
		if err != nil {
			return err
		}

		return record.InsertUpdate(ctx, target.TeamID, userID, record.OperationOperateUpdateAPI, target.Name)
	})

	return target.ID, err
}

func DetailByTargetIDs(ctx context.Context, teamID int64, targetIDs []int64) ([]*rao.APIDetail, error) {
	tx := query.Use(dal.DB()).Target
	targets, err := tx.WithContext(ctx).Where(
		tx.ID.In(targetIDs...),
		tx.TeamID.Eq(teamID),
		tx.TargetType.Eq(consts.TargetTypeAPI),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).Order(tx.Sort.Desc(), tx.CreatedAt.Desc()).Find()

	if err != nil {
		return nil, err
	}

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPI)
	cursor, err := collection.Find(ctx, bson.D{{"target_id", bson.D{{"$in", targetIDs}}}})
	if err != nil {
		return nil, err
	}
	var apis []*mao.API
	if err = cursor.All(ctx, &apis); err != nil {
		return nil, err
	}

	return packer.TransTargetsToRaoAPIDetails(targets, apis), nil
}
