package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/report"

	"github.com/gin-gonic/gin"
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
