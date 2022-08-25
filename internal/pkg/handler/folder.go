package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/folder"

	"github.com/gin-gonic/gin"
)

func SaveFolder(ctx *gin.Context) {
	var req rao.SaveFolderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := folder.Save(ctx, jwt.GetUserIDByCtx(ctx), &req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlOperFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}
