package target

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gen"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/dal/runner"
	"kp-management/internal/pkg/packer"
)

func SendAPI(ctx context.Context, targetID int64) (string, error) {
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

	return runner.RunAPI(ctx, packer.TransTargetToRaoAPIDetail(t, &a))
}

func GetSendAPIResult(ctx context.Context, retID string) (*rao.APIDebug, error) {
	var ad mao.APIDebug
	err := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectAPIDebug).
		FindOne(ctx, bson.D{{"uuid", retID}}).Decode(&ad)
	if err != nil {
		return nil, err
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
	).Order(tx.Sort.Desc(), tx.CreatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return packer.TransTargetToRaoFolderAPIList(targets), cnt, nil
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

func Trash(ctx context.Context, targetID int64) error {
	t := query.Use(dal.DB()).Target
	_, err := t.WithContext(ctx).Where(t.ID.Eq(targetID)).UpdateColumn(t.Status, consts.TargetStatusTrash)

	return err
}

func Recall(ctx context.Context, targetID int64) error {
	t := query.Use(dal.DB()).Target
	_, err := t.WithContext(ctx).Where(t.ID.Eq(targetID)).UpdateColumn(t.Status, consts.TargetStatusNormal)

	return err
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
	).Count()
}

func SceneCountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Target

	return tx.WithContext(ctx).Where(
		tx.TargetType.Eq(consts.TargetTypeScene),
		tx.TeamID.Eq(teamID),
	).Count()
}
