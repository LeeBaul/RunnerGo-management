package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-omnibus/omnibus"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/bson"
	"kp-management/internal/pkg/biz/consts"

	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/dal/runner"
	"kp-management/internal/pkg/logic/report"
)

// ListReports 测试报告列表
func ListReports(ctx *gin.Context) {
	var req rao.ListReportsReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	reports, total, err := report.ListByTeamID2(ctx, req.TeamID, req.Size, (req.Page-1)*req.Size, req.Keyword, req.StartTimeSec, req.EndTimeSec)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListReportsResp{
		Reports: reports,
		Total:   total,
	})
	return
}

// DeleteReport 删除报告
func DeleteReport(ctx *gin.Context) {
	var req rao.DeleteReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := report.DeleteReport(ctx, req.TeamID, req.ReportID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// ReportDetail 报告详情
func ReportDetail(ctx *gin.Context) {
	var req rao.GetReportReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	err, result := report.GetReportDetail(ctx, req, conf.Conf.ES.Host, conf.Conf.ES.Username, conf.Conf.ES.Password)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, result)
	return
}

// GetReportTaskDetail 获取报告任务详情
func GetReportTaskDetail(ctx *gin.Context) {
	var req rao.GetReportReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	ret, err := report.GetTaskDetail(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetReportResp{Report: ret})
	return
}

// DebugDetail 查询报告debug状态
func DebugDetail(ctx *gin.Context) {
	var req rao.GetReportReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	result := report.GetReportDebugStatus(ctx, req)
	response.SuccessWithData(ctx, result)
}

// GetDebug 获取debug日志
func GetDebug(ctx *gin.Context) {
	var req rao.GetReportReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	err, result := report.GetReportDebugLog(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, result)
}

// DebugSetting 开启debug模式
func DebugSetting(ctx *gin.Context) {
	var req rao.DebugSettingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	collection := dal.GetMongo().Database(dal.MongoDB()).Collection(consts.CollectStressDebug)
	filter := bson.D{{"report_id", req.ReportID}}
	singleResult := collection.FindOne(ctx, filter)
	result, err := singleResult.DecodeBytes()
	if err != nil {
		debug := bson.D{{"report_id", req.ReportID}, {"debug", req.Setting}}
		_, err = collection.InsertOne(ctx, debug)
		if err != nil {
			response.ErrorWithMsg(ctx, errno.ErrMongoFailed, err.Error())
			return
		}
	} else {
		_, err = result.Elements()
		if err != nil {
			debug := bson.D{{"report_id", req.ReportID}, {"debug", req.Setting}}
			_, err = collection.InsertOne(ctx, debug)
			if err != nil {
				response.ErrorWithMsg(ctx, errno.ErrMongoFailed, err.Error())
				return
			}
		} else {
			debug := bson.D{{"report_id", req.ReportID}, {"debug", req.Setting}}
			_, err = collection.UpdateOne(ctx, filter, debug)
			if err != nil {
				response.ErrorWithMsg(ctx, errno.ErrMongoFailed, err.Error())
				return
			}
		}

	}

	//reportID:debug
	//_, err = dal.GetRDB().Set(ctx, fmt.Sprintf("%d:debug", req.ReportID), req.Setting, 10*time.Minute).Result()
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrRedisFailed, err.Error())
	//	return
	//}

	response.Success(ctx)
	return
}

// ListMachines 施压机列表
func ListMachines(ctx *gin.Context) {
	var req rao.ListMachineReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	resp, err := report.ListMachines(ctx, req.ReportID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, resp)
	return
}

func StopReport(ctx *gin.Context) {
	var req rao.StopReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	_, err := resty.New().R().
		SetBody(runner.StopRunnerReq{ReportIds: omnibus.Int64sToStrings(req.ReportIDs)}).
		Post(conf.Conf.Clients.Runner.StopPlan)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrHttpFailed, err.Error())
		return
	}

	//tx := dal.GetQuery().Report
	//r, err := tx.WithContext(ctx).Where(tx.ID.In(req.ReportIDs...)).First()
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//	return
	//}

	//_, err = tx.WithContext(ctx).Where(tx.ID.In(req.ReportIDs...)).UpdateColumn(tx.Status, consts.ReportStatusFinish)
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//	return
	//}
	//
	//reportCnt, err := tx.WithContext(ctx).Where(tx.PlanID.Eq(r.PlanID)).Count()
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//	return
	//}
	//finishReportCnt, err := tx.WithContext(ctx).Where(tx.PlanID.Eq(r.PlanID), tx.Status.Eq(consts.ReportStatusFinish)).Count()
	//if err != nil {
	//	response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//	return
	//}
	//if finishReportCnt >= reportCnt {
	//	px := query.Use(dal.DB()).Plan
	//	_, err := px.WithContext(ctx).Where(px.ID.Eq(r.PlanID)).UpdateColumn(px.Status, consts.PlanStatusUnderway)
	//	if err != nil {
	//		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
	//		return
	//	}
	//}

	response.Success(ctx)
	return
}
