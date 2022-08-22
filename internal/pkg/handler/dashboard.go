package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/operation"
	"kp-management/internal/pkg/logic/target"
	"kp-management/internal/pkg/logic/user"

	"github.com/gin-gonic/gin"
)

func DashboardDefault(ctx *gin.Context) {
	var req rao.DashboardDefaultReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	u, err := user.FirstByUserID(ctx, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	operations, _, err := operation.List(ctx, req.TeamID, 5, 0)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	apiCnt, err := target.APICount(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.DashboardDefaultResp{
		User:       u,
		Operations: operations,
		PlanNum:    0,
		SceneNum:   0,
		ReportNum:  0,
		APINum:     apiCnt,
	})
}
