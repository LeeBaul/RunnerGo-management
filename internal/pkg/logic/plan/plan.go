package plan

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gen"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func ListByStatus(ctx context.Context, teamID int64, status int32, limit, offset int) ([]*rao.Plan, int64, error) {
	tx := query.Use(dal.DB()).Plan
	ret, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.Status.Eq(status)).
		Order(tx.UpdatedAt.Desc()).FindByPage(offset, limit)

	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range ret {
		userIDs = append(userIDs, r.CreateUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransPlansToRaoPlanList(ret, users), cnt, nil
}

func CountByTeamID(ctx context.Context, teamID int64) (int64, error) {
	tx := query.Use(dal.DB()).Plan

	return tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID)).Count()
}

func ListByTeamID(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64) ([]*rao.Plan, int64, error) {
	tx := query.Use(dal.DB()).Plan

	conditions := make([]gen.Condition, 0)
	conditions = append(conditions, tx.TeamID.Eq(teamID))

	if keyword != "" {
		conditions = append(conditions, tx.Name.Like(fmt.Sprintf("%%%s%%", keyword)))

		u := query.Use(dal.DB()).User
		users, err := u.WithContext(ctx).Where(u.Nickname.Like(fmt.Sprintf("%%%s%%", keyword))).Find()
		if err != nil {
			return nil, 0, err
		}

		if len(users) > 0 {
			conditions[1] = tx.RunUserID.Eq(users[0].ID)
		}
	}

	if startTimeSec > 0 && endTimeSec > 0 {
		startTime := time.Unix(startTimeSec, 0)
		endTime := time.Unix(endTimeSec, 0)
		conditions = append(conditions, tx.CreatedAt.Between(startTime, endTime))
	}

	conditions = append(conditions, tx.Status.In(consts.PlanStatusNormal, consts.PlanStatusUnderway))
	ret, cnt, err := tx.WithContext(ctx).Where(conditions...).Order(tx.UpdatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range ret {
		userIDs = append(userIDs, r.CreateUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransPlansToRaoPlanList(ret, users), cnt, nil
}

func Save(ctx context.Context, req *rao.SavePlanReq, userID int64) (int64, error) {

	tx := query.Use(dal.DB()).Plan
	p, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.PlanID)).Assign(
		tx.TeamID.Value(req.TeamID),
		tx.Name.Value(req.Name),
		tx.Status.Value(consts.PlanStatusNormal),
		tx.CreateUserID.Value(userID),
		tx.Remark.Value(req.Remark),
	).FirstOrCreate()

	return p.ID, err
}

func SaveTask(ctx context.Context, req *rao.SavePlanConfReq, userID int64) error {
	plan := packer.TransSavePlanReqToPlanModel(req, userID)
	task := packer.TransSavePlanReqToMaoTask(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {

		err := collection.FindOne(ctx, bson.D{{"plan_id", req.PlanID}}).Err()
		if err == mongo.ErrNoDocuments {
			if _, err := tx.Plan.WithContext(ctx).Omit(tx.Plan.CreateUserID).Updates(plan); err != nil {
				return err
			}

			_, err := collection.InsertOne(ctx, task)
			if err != nil {
				return err
			}

			return record.InsertCreate(ctx, plan.TeamID, userID, fmt.Sprintf("创建计划 - %s", plan.Name))
		}

		if _, err := tx.Plan.WithContext(ctx).Omit(tx.Plan.CreateUserID).Updates(plan); err != nil {
			return err
		}

		_, err = collection.UpdateOne(ctx, bson.D{{"plan_id", plan.ID}}, bson.M{"$set": task})
		if err != nil {
			return err
		}

		return record.InsertUpdate(ctx, plan.TeamID, userID, fmt.Sprintf("修改计划 - %s", plan.Name))
	})
}

func GetByPlanID(ctx context.Context, teamID, planID int64) (*rao.Plan, error) {

	tx := query.Use(dal.DB()).Plan
	p, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(planID)).First()
	if err != nil {
		return nil, err
	}

	var t *mao.Task
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
	err = collection.FindOne(ctx, bson.D{{"plan_id", planID}}).Decode(&t)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	u := query.Use(dal.DB()).User
	user, err := u.WithContext(ctx).Where(u.ID.Eq(p.CreateUserID)).First()
	if err != nil {
		return nil, err
	}

	return packer.TransTaskToRaoPlan(p, t, user), nil
}

func DeleteByPlanID(ctx context.Context, teamID, planID int64) error {
	tx := query.Use(dal.DB()).Plan
	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(planID)).Delete()

	return err
}

func SetPreinstall(ctx context.Context, req *rao.SetPreinstallReq) error {
	p := packer.TransSetPreinstallReqToMaoPreinstall(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectPreinstall)
	err := collection.FindOne(ctx, bson.D{{"team_id", req.TeamID}}).Err()
	if err == mongo.ErrNoDocuments { // 新建
		_, err := collection.InsertOne(ctx, p)

		return err
	}

	_, err = collection.UpdateOne(ctx, bson.D{
		{"team_id", req.TeamID},
	}, bson.M{"$set": p})

	return err
}

func GetPreinstall(ctx context.Context, teamID int64) (*rao.Preinstall, error) {

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectPreinstall)
	var p mao.Preinstall
	if err := collection.FindOne(ctx, bson.D{{"team_id", teamID}}).Decode(&p); err != nil {
		return nil, err
	}

	return packer.TransMaoPreinstallToRaoPreinstall(&p), nil

}

func ClonePlan(ctx context.Context, planID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		//克隆计划
		p, err := tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(planID)).First()
		if err != nil {
			return err
		}

		p.ID = 0
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		if err := tx.Plan.WithContext(ctx).Create(p); err != nil {
			return err
		}

		// 克隆场景，分组
		targets, err := tx.Target.WithContext(ctx).Where(tx.Target.PlanID.Eq(planID)).Order(tx.Target.ParentID).Find()
		if err != nil {
			return err
		}

		var sceneIDs []int64
		targetMemo := make(map[int64]int64)
		for _, target := range targets {
			if target.TargetType == consts.TargetTypeScene {
				sceneIDs = append(sceneIDs, target.ID)
			}

			oldTargetID := target.ID
			target.ID = 0
			target.ParentID = targetMemo[target.ParentID]
			target.PlanID = p.ID
			target.CreatedAt = time.Now()
			target.UpdatedAt = time.Now()
			if err := tx.Target.WithContext(ctx).Create(target); err != nil {
				return err
			}

			targetMemo[oldTargetID] = target.ID
		}

		v, err := tx.Variable.WithContext(ctx).Where(tx.Variable.SceneID.In(sceneIDs...)).Find()
		if err != nil {
			return err
		}

		for _, variable := range v {
			variable.ID = 0
			variable.SceneID = targetMemo[variable.SceneID]
			variable.CreatedAt = time.Now()
			variable.UpdatedAt = time.Now()
			if err := tx.Variable.WithContext(ctx).Create(variable); err != nil {
				return err
			}
		}

		vi, err := tx.VariableImport.WithContext(ctx).Where(tx.VariableImport.SceneID.In(sceneIDs...)).Find()
		if err != nil {
			return err
		}

		for _, variableImport := range vi {
			variableImport.ID = 0
			variableImport.SceneID = targetMemo[variableImport.SceneID]
			variableImport.CreatedAt = time.Now()
			variableImport.UpdatedAt = time.Now()
			if err := tx.VariableImport.WithContext(ctx).Create(variableImport); err != nil {
				return err
			}
		}

		var flows []*mao.Flow
		c1 := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectFlow)
		cur, err := c1.Find(ctx, bson.D{{"scene_id", bson.D{{"$in", sceneIDs}}}})
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &flows); err != nil {
			return err
		}

		for _, flow := range flows {
			flow.SceneID = targetMemo[flow.SceneID]
			if _, err := c1.InsertOne(ctx, flow); err != nil {
				return err
			}
		}

		var task mao.Task
		c2 := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
		err = c2.FindOne(ctx, bson.D{{"plan_id", planID}}).Decode(&task)
		if err != nil {
			return err
		}

		task.PlanID = p.ID
		if _, err := c2.InsertOne(ctx, task); err != nil {
			return err
		}

		return nil

	})

}
