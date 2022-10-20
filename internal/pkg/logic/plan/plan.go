package plan

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gen"
	"gorm.io/gen/field"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/record"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
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

func ListByTeamID(ctx context.Context, teamID int64, limit, offset int, keyword string, startTimeSec, endTimeSec int64, taskType, taskMode, status, sortTag int32) ([]*rao.Plan, int64, error) {
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

	if taskType > 0 {
		conditions = append(conditions, tx.TaskType.Eq(taskType))
	}

	if taskMode > 0 {
		conditions = append(conditions, tx.Mode.Eq(taskMode))
	}

	if status > 0 {
		conditions = append(conditions, tx.Status.Eq(status))
	}

	sort := make([]field.Expr, 0)
	if sortTag == 0 { // 默认排序
		sort = append(sort, tx.Rank.Desc())
		sort = append(sort, tx.ID.Desc())
	}
	if sortTag == 1 { // 创建时间倒序
		sort = append(sort, tx.CreatedAt.Desc())
	}
	if sortTag == 2 { // 创建时间正序
		sort = append(sort, tx.CreatedAt)
	}
	if sortTag == 3 { // 修改时间倒序
		sort = append(sort, tx.UpdatedAt.Desc())
	}
	if sortTag == 4 { // 修改时间正序
		sort = append(sort, tx.UpdatedAt)
	}

	conditions = append(conditions, tx.Status.In(consts.PlanStatusNormal, consts.PlanStatusUnderway))
	ret, cnt, err := tx.WithContext(ctx).Where(conditions...).Order(sort...).FindByPage(offset, limit)
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
	cnt, err := tx.WithContext(ctx).Unscoped().Where(tx.TeamID.Eq(req.TeamID)).Count()
	if err != nil {
		return 0, err
	}

	p := model.Plan{
		ID:           req.PlanID,
		TeamID:       req.TeamID,
		Name:         req.Name,
		Status:       consts.PlanStatusNormal,
		CreateUserID: userID,
		Remark:       req.Remark,
		Rank:         cnt + 1,
	}

	if req.PlanID == 0 {
		err := tx.WithContext(ctx).Create(&p)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.WithContext(ctx).Where(tx.ID.Eq(p.ID)).Omit(tx.Rank, tx.CreateUserID, tx.Status).Save(&p); err != nil {
		return 0, err
	}

	return p.ID, err
}

func SaveTask(ctx context.Context, req *rao.SavePlanConfReq, userID int64) error {
	plan := packer.TransSavePlanReqToPlanModel(req, userID)
	task := packer.TransSavePlanReqToMaoTask(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {

		err := collection.FindOne(ctx, bson.D{{"scene_id", req.SceneID}}).Err()
		if err == mongo.ErrNoDocuments {
			//if _, err := tx.Plan.WithContext(ctx).Omit(tx.Plan.CreateUserID).Updates(plan); err != nil {
			//	return err
			//}

			_, err := collection.InsertOne(ctx, task)
			if err != nil {
				return err
			}

			return record.InsertCreate(ctx, plan.TeamID, userID, fmt.Sprintf("创建计划 - %s", plan.Name))
		}

		if err == nil {
			//if _, err := tx.Plan.WithContext(ctx).Omit(tx.Plan.CreateUserID).Updates(plan); err != nil {
			//	return err
			//}

			_, err = collection.UpdateOne(ctx, bson.D{{"scene_id", req.SceneID}}, bson.M{"$set": task})
			if err != nil {
				return err
			}

			return record.InsertUpdate(ctx, plan.TeamID, userID, fmt.Sprintf("修改计划 - %s", plan.Name))
		}

		cur, err := collection.Find(ctx, bson.D{{"plan_id", req.PlanID}})
		if err != nil {
			return err
		}
		var tasks []*mao.Task
		if err := cur.All(ctx, &tasks); err != nil {
			return err
		}

		if len(tasks) > 0 {
			planType := tasks[0].TaskType
			planMode := tasks[0].TaskMode
			for i, t := range tasks {
				if i > 0 {
					if t.TaskType != planType {
						planType = consts.PlanTaskTypeMix
					}
					if t.TaskMode != planMode {
						planMode = consts.PlanModeMix
					}
				}
			}

			_, err := tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(req.PlanID)).UpdateSimple(tx.Plan.TaskType.Value(planType), tx.Plan.Mode.Value(planMode))
			if err != nil {
				return err
			}
		}

		return nil
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
	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		_, err := tx.Plan.WithContext(ctx).Where(tx.Plan.TeamID.Eq(teamID), tx.Plan.ID.Eq(planID)).Delete()
		if err != nil {
			return err
		}

		_, err = tx.Target.WithContext(ctx).Where(tx.Target.TeamID.Eq(teamID), tx.Target.PlanID.Eq(planID)).Delete()
		if err != nil {
			return err
		}

		return nil
	})
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
	err := collection.FindOne(ctx, bson.D{{"team_id", teamID}}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return packer.TransMaoPreinstallToRaoPreinstall(&p), nil

}

func ClonePlan(ctx context.Context, planID, userID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		//克隆计划
		p, err := tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(planID)).First()
		if err != nil {
			return err
		}

		p.ID = 0
		p.Name = fmt.Sprintf("%s - copy", p.Name)
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
		p.Status = consts.PlanStatusNormal
		p.CreateUserID = userID
		p.RunUserID = userID
		if err := tx.Plan.WithContext(ctx).Create(p); err != nil {
			return err
		}

		// 克隆场景，分组
		targets, err := tx.Target.WithContext(ctx).Where(tx.Target.PlanID.Eq(planID), tx.Target.Status.Eq(consts.TargetStatusNormal)).Order(tx.Target.ParentID).Find()
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

		// 克隆场景变量
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

		// 克隆导入变量
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

		// 克隆流程
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

		// 克隆任务配置
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

		return record.InsertCreate(ctx, p.TeamID, userID, fmt.Sprintf("克隆计划 - %s", p.Name))
	})

}
