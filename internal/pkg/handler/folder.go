package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/folder"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"

	"github.com/gin-gonic/gin"
)

func SaveFolder(ctx *gin.Context) {
	var req rao.SaveFolderReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	if err := folder.SaveFolder(ctx, &req); err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}
