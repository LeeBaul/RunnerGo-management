package scene

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
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
