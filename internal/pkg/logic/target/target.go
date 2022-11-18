package target

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gen"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/dal/runner"
	"kp-management/internal/pkg/packer"
)

func SendSceneAPI(ctx context.Context, teamID, sceneID int64, nodeID string) (string, error) {
	var f mao.Flow
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	err := collection.FindOne(ctx, bson.D{{"scene_id", sceneID}}).Decode(&f)
	if err != nil {
		return "", err
	}

	var n mao.Node
	if err := bson.Unmarshal(f.Nodes, &n); err != nil {
		return "", err
	}

	for _, node := range n.Nodes {
		if node.ID == nodeID {

			tx := dal.GetQuery().Variable
			variables, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.Type.Eq(consts.VariableTypeGlobal)).Find()
			if err != nil {
				return "", err
			}

			var vs []*rao.KVVariable
			for _, v := range variables {
				vs = append(vs, &rao.KVVariable{
					Key:   v.Var,
					Value: v.Val,
				})
			}

			node.API.Variable = vs
			return runner.RunAPI(ctx, node.API)
		}
	}

	return "", nil
}

func SendScene(ctx context.Context, teamID, sceneID, userID int64) (string, error) {
	tx := dal.GetQuery().Target
	t, err := tx.WithContext(ctx).Where(tx.ID.Eq(sceneID), tx.TargetType.Eq(consts.TargetTypeScene)).First()
	if err != nil {
		return "", err
	}

	var f mao.Flow
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
	err = collection.FindOne(ctx, bson.D{{"scene_id", sceneID}}).Decode(&f)
	if err != nil {
		return "", err
	}

	vi := dal.GetQuery().VariableImport
	vis, err := vi.WithContext(ctx).Where(vi.SceneID.Eq(sceneID)).Limit(5).Find()
	if err != nil {
		return "", err
	}

	sv := dal.GetQuery().Variable
	sceneVariables, err := sv.WithContext(ctx).Where(sv.SceneID.Eq(sceneID)).Find()
	if err != nil {
		return "", err
	}

	variables, err := sv.WithContext(ctx).Where(sv.TeamID.Eq(teamID)).Find()
	if err != nil {
		return "", err
	}

	if err := record.InsertRun(ctx, teamID, userID, record.OperationOperateRunScene, t.Name); err != nil {
		return "", err
	}

	req := packer.TransMaoFlowToRaoSceneFlow(t, &f, vis, sceneVariables, variables)
	return runner.RunScene(ctx, req)
}

func GetSendSceneResult(ctx context.Context, retID string) ([]*rao.SceneDebug, error) {
	cur, err := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectSceneDebug).
		Find(ctx, bson.D{{"uuid", retID}})
	if err != nil {
		return nil, err
	}

	var sds []*mao.SceneDebug
	if err := cur.All(ctx, &sds); err != nil {
		return nil, err
	}

	if len(sds) == 0 {
		return nil, nil
	}

	return packer.TransMaoSceneDebugsToRaoSceneDebugs(sds), nil

}

func SendAPI(ctx context.Context, teamID, targetID int64) (string, error) {
	tx := dal.GetQuery().Target
	t, err := tx.WithContext(ctx).Where(tx.ID.Eq(targetID)).First()
	if err != nil {
		return "", err
	}

	var a mao.API
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPI)
	err = collection.FindOne(ctx, bson.D{{"target_id", targetID}}).Decode(&a)
	if err != nil {
		return "", err
	}

	v := dal.GetQuery().Variable
	variables, err := v.WithContext(ctx).Where(v.TeamID.Eq(teamID)).Find()
	if err != nil {
		return "", err
	}

	return runner.RunAPI(ctx, packer.TransTargetToRaoAPIDetail(t, &a, variables))
}

func GetSendAPIResult(ctx context.Context, retID string) (*rao.APIDebug, error) {
	var ad mao.APIDebug
	err := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPIDebug).
		FindOne(ctx, bson.D{{"uuid", retID}}).Decode(&ad)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return packer.TransMaoAPIDebugToRaoAPIDebug(&ad), nil
}

func ListFolderAPI(ctx context.Context, teamID int64, limit, offset int) ([]*rao.FolderAPI, int64, error) {
	tx := query.Use(dal.DB()).Target
	targets, cnt, err := tx.WithContext(ctx).Where(
		tx.TeamID.Eq(teamID),
		tx.TargetType.In(consts.TargetTypeFolder, consts.TargetTypeAPI),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).Order(tx.Sort, tx.CreatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return packer.TransTargetToRaoFolderAPIList(targets), cnt, nil
}

func SortTarget(ctx context.Context, req *rao.SortTargetReq) error {
	tx := dal.GetQuery().Target

	for _, target := range req.Targets {
		_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(target.TeamID), tx.ID.Eq(target.TargetID)).UpdateSimple(tx.Sort.Value(target.Sort), tx.ParentID.Value(target.ParentID))
		if err != nil {
			return err
		}
	}

	return nil
}

func ListGroupScene(ctx context.Context, teamID int64, source int32, limit, offset int, planID int64) ([]*rao.GroupScene, int64, error) {
	tx := query.Use(dal.DB()).Target

	condition := make([]gen.Condition, 0)
	condition = append(condition, tx.TeamID.Eq(teamID))
	condition = append(condition, tx.TargetType.In(consts.TargetTypeGroup, consts.TargetTypeScene))
	condition = append(condition, tx.Status.Eq(consts.TargetStatusNormal))
	condition = append(condition, tx.Source.Eq(source))

	if source == consts.TargetSourcePlan {
		condition = append(condition, tx.PlanID.Eq(planID))
	}

	targets, cnt, err := tx.WithContext(ctx).Where(condition...).
		Order(tx.Sort.Desc(), tx.CreatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return packer.TransTargetsToRaoGroupSceneList(targets), cnt, nil
}

func ListTrashFolderAPI(ctx context.Context, teamID int64, limit, offset int) ([]*rao.FolderAPI, int64, error) {
	tx := query.Use(dal.DB()).Target
	targets, cnt, err := tx.WithContext(ctx).Where(
		tx.TeamID.Eq(teamID),
		tx.TargetType.In(consts.TargetTypeFolder, consts.TargetTypeAPI),
		tx.Status.Eq(consts.TargetStatusTrash),
	).Order(tx.Sort.Desc(), tx.CreatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return packer.TransTargetToRaoFolderAPIList(targets), cnt, nil
}

func Trash(ctx context.Context, targetID, userID int64) error {
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

		if t.TargetType == consts.TargetTypeScene {
			if err := record.InsertDelete(ctx, t.TeamID, userID, record.OperationOperateDeleteScene, t.Name); err != nil {
				return err
			}
		}
		return nil
	})
}

func Recall(ctx context.Context, targetID int64) error {
	t := query.Use(dal.DB()).Target
	_, err := t.WithContext(ctx).Where(t.ID.Eq(targetID)).UpdateColumn(t.Status, consts.TargetStatusNormal)
	if err != nil {
		return err
	}

	_, err = t.WithContext(ctx).Where(t.ParentID.Eq(targetID)).UpdateColumn(t.Status, consts.TargetStatusNormal)
	if err != nil {
		return err
	}

	return nil
}

func Delete(ctx context.Context, targetID int64) error {
	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.Target.WithContext(ctx).Where(tx.Target.ID.Eq(targetID)).Delete(); err != nil {
			return err
		}

		filter := bson.D{{"target_id", targetID}}

		//if _, err := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFolder).DeleteOne(ctx, filter); err != nil {
		//	return err
		//}

		if _, err := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPI).DeleteOne(ctx, filter); err != nil {
			return err
		}

		return nil
	})
}

func APICountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq(consts.TargetTypeAPI),
		tx.TeamID.Eq(teamID),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
	).Count()
}

func SceneCountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq(consts.TargetTypeScene),
		tx.Status.Eq(consts.TargetStatusNormal),
		tx.Source.Eq(consts.TargetSourceNormal),
		tx.TeamID.Eq(teamID),
	).Count()
}
