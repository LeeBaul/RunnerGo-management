package plan

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransPlansToResp(ret, users), cnt, nil
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

	ret, cnt, err := tx.WithContext(ctx).Where(conditions...).
		Order(tx.UpdatedAt.Desc()).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	var userIDs []int64
	for _, r := range ret {
		userIDs = append(userIDs, r.RunUserID)
	}

	u := query.Use(dal.DB()).User
	users, err := u.WithContext(ctx).Where(u.ID.In(userIDs...)).Find()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransPlansToResp(ret, users), cnt, nil
}

func Save(ctx context.Context, req *rao.SavePlanReq, userID int64) error {
	plan := packer.TransSavePlanReqToModel(req, userID)
	task := packer.TransSavePlanReqToTask(req)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if plan.ID == 0 {
			if err := tx.Plan.WithContext(ctx).Create(plan); err != nil {
				return err
			}

			task.PlanID = plan.ID
			_, err := collection.InsertOne(ctx, task)

			record.InsertCreate(ctx, plan.TeamID, userID, fmt.Sprintf("创建计划 - %s", plan.Name))

			return err
		}

		if _, err := tx.Plan.WithContext(ctx).Omit(tx.Plan.CreateUserID).Updates(plan); err != nil {
			return err
		}

		_, err := collection.UpdateOne(ctx, bson.D{{"plan_id", plan.ID}}, bson.M{"$set": task})

		record.InsertUpdate(ctx, plan.TeamID, userID, fmt.Sprintf("修改计划 - %s", plan.Name))

		return err
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
	if err != nil {
		return nil, err
	}

	return packer.TransTaskToPlanResp(p, t), nil
}
