package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-omnibus/proof"
	"gorm.io/gorm"
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

	// 判断任务配置类型
	if req.TaskType == 1 { // 普通任务
		err := query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
			// 1、先去把定时任务数据删掉
			_, err := tx.TimedTaskConf.WithContext(ctx).
				Where(tx.TimedTaskConf.TeamID.Eq(req.TeamID)).
				Where(tx.TimedTaskConf.PlanID.Eq(req.PlanID)).
				Where(tx.TimedTaskConf.SenceID.Eq(req.SceneID)).Delete()
			if err != nil {
				proof.Infof("保存配置--不存在定时任务或删除mysql失败,err:", err)
			}
			// 2、去mg里面创建或更新配置数据
			err = collection.FindOne(ctx, bson.D{{"scene_id", req.SceneID}}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments { // 如果没查到
					// 如果没有当前场景配置，则创建
					_, err := collection.InsertOne(ctx, task)
					if err != nil {
						proof.Errorf("保存配置--在mg保存任务配置失败，err:", err)
						return err
					}

					// 记录操作日志
					err = record.InsertCreate(ctx, plan.TeamID, userID, record.OperationOperateCreatePlan, plan.Name)
					if err != nil {
						proof.Errorf("保存配置--保存操作日志失败，err:", err)
						return err
					}
				} else {
					proof.Errorf("保存配置--查找任务配置失败，err:", err)
					return err
				}
			} else { // 如果mg里面有数据的话
				_, err = collection.UpdateOne(ctx, bson.D{{"scene_id", req.SceneID}}, bson.M{"$set": task})
				if err != nil {
					proof.Errorf("保存配置--更新任务配置项失败，err:", err)
					return err
				}

				err := record.InsertUpdate(ctx, plan.TeamID, userID, record.OperationOperateUpdatePlan, plan.Name)
				if err != nil {
					proof.Errorf("保存配置--保存操作日志失败，err:", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			proof.Errorf("保存配置--保存普通任务配置失败，err:", err)
			return err
		}
	} else { // 定时任务
		// 1、先去把mg里面可能存在的普通任务配置给删掉
		_, err := collection.DeleteOne(ctx, bson.D{{"plan_id", req.PlanID}, {"scene_id", req.SceneID}})
		if err != nil {
			proof.Infof("保存配置--不存在普通任务或删除mg失败,err:", err)
		}

		// 把定时任务保存到数据库中
		tx := dal.GetQuery().TimedTaskConf
		// 查询当前定时任务是否存在
		_, err = tx.WithContext(ctx).
			Where(tx.TeamID.Eq(req.TeamID)).
			Where(tx.PlanID.Eq(req.PlanID)).
			Where(tx.SenceID.Eq(req.SceneID)).First()
		if err != nil && err != gorm.ErrRecordNotFound { // 查询出错
			proof.Infof("保存配置--查询定时任务数据失败，err:", req)
			return err
		} else if err == gorm.ErrRecordNotFound { // 数据不存在
			// 新增配置
			timingTaskConfig, err := packer.TransSaveTimingTaskConfigReqToModelData(req)
			if err != nil {
				proof.Infof("保存配置--压缩mode_conf为字符串时失败", err)
				return err
			}
			err = tx.WithContext(ctx).Create(timingTaskConfig)
			if err != nil {
				proof.Infof("保存配置--定时任务配置项保存失败，err：", err)
				return err
			}
		} else {
			// 把mode_conf压缩成字符串
			modeConfString, err := json.Marshal(req.ModeConf)
			if err != nil {
				proof.Infof("保存配置--压缩mode_conf为字符串时失败", err)
				return err
			}

			// 修改配置
			updateData := make(map[string]interface{}, 3)
			updateData["frequency"] = req.TimedTaskConf.Frequency
			updateData["task_exec_time"] = req.TimedTaskConf.TaskExecTime
			updateData["task_close_time"] = req.TimedTaskConf.TaskCloseTime
			updateData["task_mode"] = req.Mode
			updateData["mode_conf"] = modeConfString
			updateData["status"] = consts.TimedTaskWaitEnable
			_, err = tx.WithContext(ctx).Where(tx.TeamID.Eq(req.TeamID)).
				Where(tx.PlanID.Eq(req.PlanID)).
				Where(tx.SenceID.Eq(req.SceneID)).Updates(updateData)
			if err != nil {
				proof.Errorf("保存配置--更新定时任务配置失败，err:", err)
				return err
			}
		}

	}

	// 判断当前计划是否是混合任务
	cur, err := collection.Find(ctx, bson.D{{"plan_id", req.PlanID}})
	if err != nil {
		proof.Errorf("保存配置--查找当前计划下所有任务配置失败，plan_id:", req.PlanID, " err:", err)
		return err
	}
	var tasks []*mao.Task
	if err := cur.All(ctx, &tasks); err != nil {
		proof.Errorf("保存配置--解析当前计划下所有任务配置失败，plan_id:", req.PlanID, " err:", err)
		return err
	}

	// 查询当前计划下是否有定时任务
	tx := dal.GetQuery()
	timedTaskList, err := tx.TimedTaskConf.WithContext(ctx).Where(tx.TimedTaskConf.PlanID.Eq(req.PlanID)).Find()
	if err != nil && err != gorm.ErrRecordNotFound {
		proof.Infof("保存配置--查询当前计划下是否有定时任务时出错，err:", err)
		return err
	}

	var planType int32 = 0
	var planMode int32 = 0
	if len(tasks) > 0 && len(timedTaskList) > 0 { // 两种类型都有
		planType = consts.PlanTaskTypeMix
		planMode = tasks[0].TaskMode
		for i, t := range tasks {
			if i > 0 {
				if t.TaskMode != planMode && planMode != 0 {
					planMode = consts.PlanModeMix
					break
				}
			}
		}
		for _, timeTaskConf := range timedTaskList {
			if timeTaskConf.TaskMode != planMode && planMode != 0 {
				planMode = consts.PlanModeMix
				break
			}
		}
	} else if len(tasks) > 0 {
		planType = consts.PlanTaskTypeNormal
		// 模式
		planMode = tasks[0].TaskMode
		for i, t := range tasks {
			if i > 0 {
				if t.TaskMode != planMode && planMode != 0 {
					planMode = consts.PlanModeMix
					break
				}
			}
		}
	} else if len(timedTaskList) > 0 {
		planType = consts.PlanTaskTypeCronjob
		// 模式
		planMode = timedTaskList[0].TaskMode
		for _, timeTaskConf := range timedTaskList {
			if timeTaskConf.TaskMode != planMode && planMode != 0 {
				planMode = consts.PlanModeMix
				break
			}
		}
	}

	_, err = tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(req.PlanID)).UpdateSimple(tx.Plan.TaskType.Value(planType), tx.Plan.Mode.Value(planMode))
	if err != nil {
		proof.Errorf("保存配置--修改计划的任务类型和也测模式失败，err:", err)
		return err
	}

	// 最后的返回
	return nil
}

func GetPlanTask(ctx context.Context, planID, sceneID int64) (*rao.PlanTask, error) {
	// 初始化返回值
	planTaskConf := new(rao.PlanTask)

	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
	t := new(mao.Task)
	err := collection.FindOne(ctx, bson.D{{"scene_id", sceneID}, {"plan_id", planID}}).Decode(&t)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			proof.Errorf("获取任务配置详情--从mg查询失败，err:", err)
			return nil, err
		} else { // 定时任务
			// 查询定时任务信息
			var timingTaskConfigInfo *model.TimedTaskConf
			tx := query.Use(dal.DB()).TimedTaskConf
			timingTaskConfigInfo, err = tx.WithContext(ctx).Where(tx.PlanID.Eq(planID), tx.SenceID.Eq(sceneID)).First()
			if err != nil {
				return nil, err
			}
			if timingTaskConfigInfo != nil {
				var modeConf rao.ModeConf
				err := json.Unmarshal([]byte(timingTaskConfigInfo.ModeConf), &modeConf)
				if err != nil {
					proof.Errorf("获取任务配置详情--解析定时任务详细配置失败，err:", err)
					return nil, err
				}

				planTaskConf = &rao.PlanTask{
					PlanID:   timingTaskConfigInfo.PlanID,
					SceneID:  timingTaskConfigInfo.SenceID,
					TaskType: timingTaskConfigInfo.TaskType,
					Mode:     timingTaskConfigInfo.TaskMode,
					ModeConf: &modeConf,
					TimedTaskConf: &rao.TimedTaskConf{
						Frequency:     timingTaskConfigInfo.Frequency,
						TaskExecTime:  timingTaskConfigInfo.TaskExecTime,
						TaskCloseTime: timingTaskConfigInfo.TaskCloseTime,
					},
				}
			}
		}
	} else { // 普通任务
		planTaskConf = &rao.PlanTask{
			PlanID:   t.PlanID,
			SceneID:  t.SceneID,
			TaskType: t.TaskType,
			Mode:     t.TaskMode,
			ModeConf: &rao.ModeConf{
				ReheatTime:       t.ModeConf.ReheatTime,
				RoundNum:         t.ModeConf.RoundNum,
				Concurrency:      t.ModeConf.Concurrency,
				ThresholdValue:   t.ModeConf.ThresholdValue,
				StartConcurrency: t.ModeConf.StartConcurrency,
				Step:             t.ModeConf.Step,
				StepRunTime:      t.ModeConf.StepRunTime,
				MaxConcurrency:   t.ModeConf.MaxConcurrency,
				Duration:         t.ModeConf.Duration,
			},
		}
	}

	return planTaskConf, nil
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

func DeleteByPlanID(ctx context.Context, teamID, planID, userID int64) error {
	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		p, err := tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(planID)).First()
		if err != nil {
			return err
		}

		if _, err := tx.Plan.WithContext(ctx).Where(tx.Plan.TeamID.Eq(teamID), tx.Plan.ID.Eq(planID)).Delete(); err != nil {
			return err
		}

		if _, err = tx.Target.WithContext(ctx).Where(tx.Target.TeamID.Eq(teamID), tx.Target.PlanID.Eq(planID)).Delete(); err != nil {
			return err
		}

		return record.InsertDelete(ctx, teamID, userID, record.OperationOperateDeletePlan, p.Name)
	})
}

func ClonePlan(ctx context.Context, planID, teamID, userID int64) error {

	return dal.GetQuery().Transaction(func(tx *query.Query) error {
		//克隆计划
		p, err := tx.Plan.WithContext(ctx).Where(tx.Plan.ID.Eq(planID)).First()
		if err != nil {
			return err
		}

		cnt, err := tx.Plan.WithContext(ctx).Unscoped().Where(tx.Plan.TeamID.Eq(teamID)).Count()
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
		p.Rank = cnt + 1
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
		var tasks []*mao.Task
		c2 := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectTask)
		cur, err = c2.Find(ctx, bson.D{{"plan_id", planID}})
		if err != nil {
			return err
		}
		if err := cur.All(ctx, &tasks); err != nil {
			return err
		}

		for _, task := range tasks {
			task.PlanID = p.ID
			task.SceneID = targetMemo[task.SceneID]
			if _, err := c2.InsertOne(ctx, task); err != nil {
				return err
			}
		}

		return record.InsertCreate(ctx, p.TeamID, userID, record.OperationOperateClonePlan, p.Name)
	})

}
