package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/target"

	"github.com/gin-gonic/gin"
)

func DashboardDefault(ctx *gin.Context) {
	var teamID int64

	apiCnt, err := target.APICount(ctx, teamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.DashboardDefaultResp{
		PlanNum:   0,
		SceneNum:  0,
		ReportNum: 0,
		APINum:    apiCnt,
	})
}
