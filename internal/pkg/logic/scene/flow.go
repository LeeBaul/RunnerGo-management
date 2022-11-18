package scene

import (
	"context"
	"kp-management/internal/pkg/biz/record"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SaveFlow(ctx context.Context, req *rao.SaveFlowReq) error {
	flow := packer.TransSaveFlowReqToMaoFlow(req)
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)

	err := collection.FindOne(ctx, bson.D{{"scene_id", req.SceneID}}).Err()
	if err == mongo.ErrNoDocuments { // 新建
		_, err := collection.InsertOne(ctx, flow)
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.D{
		{"scene_id", flow.SceneID},
	}, bson.M{"$set": flow})

	return err
}

func GetFlow(ctx context.Context, sceneID int64) (*rao.GetFlowResp, error) {
	var ret mao.Flow

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	err := collection.FindOne(ctx, bson.D{{"scene_id", sceneID}}).Decode(&ret)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return packer.TransMaoFlowToRaoGetFowResp(&ret), nil
}

func BatchGetFlow(ctx context.Context, sceneIDs []int64) ([]*rao.Flow, error) {

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	cursor, err := collection.Find(ctx, bson.D{{"scene_id", bson.D{{"$in", sceneIDs}}}})
	if err != nil {
		return nil, err
	}

	var flows []*mao.Flow
	if err := cursor.All(ctx, &flows); err != nil {
		return nil, err
	}

	return packer.TransMaoFlowsToRaoFlows(flows), nil
}

func DeleteScene(ctx context.Context, targetID, userID int64) error {
	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		t, err := tx.Target.WithContext(ctx).Where(tx.Target.ID.Eq(targetID)).First()
		if err != nil {
			return err
		}

		if _, err := tx.Target.WithContext(ctx).Where(tx.Target.ID.Eq(targetID)).Delete(); err != nil {
			return err
		}

		if _, err = tx.Target.WithContext(ctx).Where(tx.Target.ParentID.Eq(targetID)).Delete(); err != nil {
			return err
		}

		// 从mg里面删除当前场景对应的flow
		collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
		_, err = collection.DeleteOne(ctx, bson.D{{"scene_id", targetID}})
		if err != nil {
			return err
		}

		if t.TargetType == consts.TargetTypeScene {
			if err := record.InsertDelete(ctx, t.TeamID, userID, record.OperationOperateDeleteScene, t.Name); err != nil {
				return err
			}
		}
		return nil
	})
}
