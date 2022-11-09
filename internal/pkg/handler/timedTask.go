package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/proof"
	"gorm.io/gen"
	"gorm.io/gorm"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"time"
)

func TimedTaskExec() {
	ctx := context.Background()
	tx := query.Use(dal.DB()).TimedTaskConf
	// 开启定时任务轮询
	for {
		// 组装查询条件
		conditions := make([]gen.Condition, 0)
		// 当前时间
		nowTime := time.Now().Unix()
		nextTime := time.Now().Unix() + 60
		conditions = append(conditions, tx.Status.Eq(0))

		// 从数据库当中，查出当前需要执行的定时任务
		timedTaskData, err := tx.WithContext(ctx).Where(conditions...).Find()
		if err == nil { // 查到了数据
			// 组装运行计划参数
			for _, timedTaskInfo := range timedTaskData {
				// 排除过期的定时任务
				if timedTaskInfo.TaskCloseTime < nowTime {
					//  todo 把当前定时任务状态变成已过期
					_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(timedTaskInfo.TeamID)).
						Where(tx.PlanID.Eq(timedTaskInfo.PlanID)).
						Where(tx.SenceID.Eq(timedTaskInfo.SenceID)).
						UpdateColumn(tx.Status, consts.TimedTaskTimeout)
					if err != nil {
						proof.Infof("定时任务过期状态修改失败，err：", err)
					}
					continue
				}

				// 获取定时任务的执行时间相关数据
				tm := time.Unix(timedTaskInfo.TaskExecTime, 0)
				taskHour := tm.Hour()
				taskMinute := tm.Minute()
				taskWeekday := tm.Weekday()
				taskDay := tm.Day()

				// 当前时间的 时，分
				nowTimeInfo := time.Unix(nowTime, 0)
				nowHour := nowTimeInfo.Hour()
				nowMinute := nowTimeInfo.Minute()
				nowWeekday := nowTimeInfo.Weekday()
				nowDay := nowTimeInfo.Day()

				// 根据不同的任务频次，进行不同的运行逻辑
				switch timedTaskInfo.Frequency {
				case 0: // 一次
					if timedTaskInfo.TaskExecTime < nowTime || timedTaskInfo.TaskExecTime > nextTime {
						continue
					}
				case 1: // 每天
					// 比较当前时间是否等于定时任务的时间
					if taskHour != nowHour || taskMinute != nowMinute {
						continue
					}

				case 2: // 每周
					// 比较当前周几是否等于定时任务的时间
					if taskWeekday != nowWeekday || taskHour != nowHour || taskMinute != nowMinute {
						continue
					}

				case 3: // 每月
					// 比较当前每月几号是否等于定时任务的时间
					if taskDay != nowDay || taskHour != nowHour || taskMinute != nowMinute {
						continue
					}
				}

				// 执行定时任务计划
				ctx2 := &gin.Context{}
				err := runTimedTask(ctx2, timedTaskInfo)
				if err != nil {
					proof.Infof("定时任务运行失败，任务信息：", timedTaskInfo, " err：", err)
				}

			}
		} else if err != gorm.ErrRecordNotFound {
			proof.Infof("定时任务查询数据库出错，err：", err)
			continue
		}

		// 睡眠一分钟，再循环执行
		time.Sleep(60 * time.Second)
	}
}

func runTimedTask(ctx *gin.Context, timedTaskInfo *model.TimedTaskConf) error {
	// 开始执行计划
	sceneIds := make([]int64, 0, 1)
	sceneIds = append(sceneIds, timedTaskInfo.SenceID)
	runStressParams := RunStressReq{
		PlanID:  timedTaskInfo.PlanID,
		TeamID:  timedTaskInfo.TeamID,
		SceneID: sceneIds,
		UserID:  jwt.GetUserIDByCtx(ctx),
	}
	// 进入执行计划方法
	_, runErr := RunStress(ctx, runStressParams)
	if runErr != nil {
		proof.Infof("定时任务执行失败，定时任务信息：", runStressParams, " err：", runErr)
	} else {
		tx := query.Use(dal.DB()).TimedTaskConf
		_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(timedTaskInfo.TeamID)).
			Where(tx.PlanID.Eq(timedTaskInfo.PlanID)).
			Where(tx.SenceID.Eq(timedTaskInfo.SenceID)).
			UpdateColumn(tx.Status, consts.PlanStatusUnderway)
		if err != nil {
			proof.Infof("定时任务状态修改失败，err：", err, " 参数为：", runStressParams)
		}
	}
	return nil
}
