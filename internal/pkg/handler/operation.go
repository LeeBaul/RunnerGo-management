package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/operation"

	"github.com/gin-gonic/gin"
)

func ListOperation(ctx *gin.Context) {
	var req rao.ListOperationReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	operations, total, err := operation.List(ctx, req.TeamID, req.Size, (req.Page-1)*req.Size)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListOperationResp{
		Operations: operations,
		Total:      total,
	})
	return
}
