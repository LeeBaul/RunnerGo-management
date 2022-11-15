package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/plan"
	"kp-management/internal/pkg/logic/preinstall"
)

// SavePreinstall 保存预设设置
func SavePreinstall(ctx *gin.Context) {
	var req rao.SavePreinstallReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	errNum, err := preinstall.SavePreinstall(ctx, &req)
	if err != nil {
		fmt.Println(999)
		response.ErrorWithMsg(ctx, errNum, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// GetPreinstallDetail 获取预设设置
func GetPreinstallDetail(ctx *gin.Context) {
	var req rao.GetPreinstallReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	p, err := plan.GetPreinstall(ctx, req.PlanID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetPreinstallResp{Preinstall: p})
}
