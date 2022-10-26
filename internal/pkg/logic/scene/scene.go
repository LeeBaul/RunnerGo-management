package scene

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func Save(ctx context.Context, req *rao.SaveSceneReq, userID int64) (int64, string, error) {
	target := packer.TransSaveSceneReqToTargetModel(req, userID)
	//scene := packer.TransSaveSceneReqToMaoScene(req)

	//collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectScene)

	err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if target.ID == 0 {
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			//scene.TargetID = target.ID
			//_, err := collection.InsertOne(ctx, scene)

			return record.InsertCreate(ctx, target.TeamID, userID, record.OperationOperateCreateScene, target.Name)
		}

		if _, err := tx.Target.WithContext(ctx).Omit(tx.Target.CreatedUserID).Updates(target); err != nil {
			return err
		}

		//_, err := collection.UpdateOne(ctx, bson.D{{"target_id", target.ID}}, bson.M{"$set": scene})

		return record.InsertUpdate(ctx, target.TeamID, userID, record.OperationOperateUpdateScene, target.Name)
	})

	// 拷贝导入场景的原来的变量
	if req.ImportSceneID > 0 {
		tx := dal.GetQuery()
		v, err := tx.Variable.WithContext(ctx).Where(tx.Variable.SceneID.Eq(req.ImportSceneID)).Find()
		if err != nil {
			return 0, "", err
		}
		if len(v) > 0 {
			for _, variable := range v {
				err := tx.Variable.WithContext(ctx).Create(&model.Variable{
					TeamID:      variable.TeamID,
					Type:        variable.Type,
					Var:         variable.Var,
					Val:         variable.Val,
					Description: variable.Description,
					SceneID:     target.ID,
				})
				if err != nil {
					return 0, "", err
				}
			}
		}

		vi, err := tx.VariableImport.WithContext(ctx).Where(tx.VariableImport.SceneID.Eq(req.ImportSceneID)).Find()
		if err != nil {
			return 0, "", err
		}
		if len(vi) > 0 {
			for _, variableImport := range vi {
				err := tx.VariableImport.WithContext(ctx).Create(&model.VariableImport{
					TeamID:     variableImport.TeamID,
					SceneID:    target.ID,
					Name:       variableImport.Name,
					URL:        variableImport.URL,
					UploaderID: variableImport.UploaderID,
				})
				if err != nil {
					return 0, "", err
				}
			}
		}
	}

	return target.ID, target.Name, err
}

func BatchGetByTargetID(ctx context.Context, teamID int64, targetIDs []int64) ([]*rao.Scene, error) {
	tx := query.Use(dal.DB()).Target
	t, err := tx.WithContext(ctx).Where(
		tx.ID.In(targetIDs...),
		tx.TeamID.Eq(teamID),
		tx.TargetType.Eq(consts.TargetTypeScene),
		tx.Status.Eq(consts.TargetStatusNormal),
	).Find()

	if err != nil {
		return nil, err
	}

	//collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectScene)
	//cursor, err := collection.Find(ctx, bson.D{{"target_id", bson.D{{"$in", targetIDs}}}})
	//if err != nil {
	//	return nil, err
	//}
	//var s []*mao.Scene
	//if err := cursor.All(ctx, &s); err != nil {
	//	return nil, err
	//}

	return packer.TransTargetToRaoScene(t, nil), nil
}
