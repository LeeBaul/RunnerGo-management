package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/operation"
	"kp-management/internal/pkg/logic/plan"
	"kp-management/internal/pkg/logic/report"
	"kp-management/internal/pkg/logic/target"
	"kp-management/internal/pkg/logic/user"

	"github.com/gin-gonic/gin"
)

// DashboardDefault 首页控制台
func DashboardDefault(ctx *gin.Context) {
	var req rao.DashboardDefaultReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	u, err := user.FirstByUserID(ctx, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	operations, _, err := operation.List(ctx, req.TeamID, 5, 0)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	apiCnt, err := target.APICountByTeamID(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	sceneCnt, err := target.SceneCountByTeamID(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	planCnt, err := plan.CountByTeamID(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	reportCnt, err := report.CountByTeamID(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.DashboardDefaultResp{
		User:       u,
		Operations: operations,
		PlanNum:    planCnt,
		SceneNum:   sceneCnt,
		ReportNum:  reportCnt,
		APINum:     apiCnt,
	})
	return
}
