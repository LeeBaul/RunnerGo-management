package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/api"

	"github.com/gin-gonic/gin"
)

func SaveTarget(ctx *gin.Context) {
	var req rao.CreateTargetReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	if err := api.SaveAPI(ctx, &req); err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func ListTarget(ctx *gin.Context) {

}
