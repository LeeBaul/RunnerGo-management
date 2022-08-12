package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"

	"github.com/gin-gonic/gin"
)

func SaveTarget(ctx *gin.Context) {
	var req rao.CreateTargetReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

}

func ListTarget(ctx *gin.Context) {

}
