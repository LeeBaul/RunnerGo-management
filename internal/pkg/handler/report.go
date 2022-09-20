package handler

import (
	"github.com/gin-gonic/gin"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/report"
)

// ListReports 测试报告列表
func ListReports(ctx *gin.Context) {
	var req rao.ListReportsReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	reports, total, err := report.ListByTeamID(ctx, req.TeamID, req.Size, (req.Page-1)*req.Size, req.Keyword, req.StartTimeSec, req.EndTimeSec)
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
	var req rao.GetReport
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

func GetReportTaskDetail(ctx *gin.Context) {
	var req rao.GetReport
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	err, result := report.GetTaskDetail(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, result)
}

func GetDebug(ctx *gin.Context) {
	var req rao.GetReport
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	err, result := report.GetReportDebugLog(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
	response.SuccessWithData(ctx, result)
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
