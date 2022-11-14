package handler

import (
	"github.com/gin-gonic/gin"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/machine"
)

// GetMachineList 获取机器列表
func GetMachineList(ctx *gin.Context) {
	var req rao.GetMachineListParam
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	res, total, err := machine.GetMachineList(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.SuccessWithData(ctx, rao.GetMachineListResponse{
		MachineList: res,
		Total:       total,
	})
	return
}

// ChangeMachineOnOff 压力机启用或卸载
func ChangeMachineOnOff(ctx *gin.Context) {
	var req rao.ChangeMachineOnOff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := machine.ChangeMachineOnOff(ctx, req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	response.Success(ctx)
	return
}
