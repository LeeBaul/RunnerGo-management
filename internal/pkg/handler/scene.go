package handler

import (
	"github.com/gin-gonic/gin"

	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/dal/runner"
	"kp-management/internal/pkg/logic/scene"
	"kp-management/internal/pkg/logic/target"
)

// SendScene 调试场景
func SendScene(ctx *gin.Context) {
	var req rao.SendSceneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	retID, err := target.SendScene(ctx, req.TeamID, req.SceneID, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SendSceneResp{RetID: retID})
	return
}

// StopScene 停止调试场景
func StopScene(ctx *gin.Context) {
	var req rao.StopSceneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err := runner.StopScene(ctx, &req)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrHttpFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// SendSceneAPI 场景调试接口
func SendSceneAPI(ctx *gin.Context) {
	var req rao.SendSceneAPIReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	retID, err := target.SendSceneAPI(ctx, req.TeamID, req.SceneID, req.NodeID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SendSceneAPIResp{RetID: retID})
	return
}

// GetSendSceneResult 获取调试场景结果
func GetSendSceneResult(ctx *gin.Context) {
	var req rao.GetSendSceneResultReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	s, err := target.GetSendSceneResult(ctx, req.RetID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetSendSceneResultResp{Scenes: s})
	return
}

// SaveScene 创建/修改场景
func SaveScene(ctx *gin.Context) {
	var req rao.SaveSceneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	targetID, targetName, err := scene.Save(ctx, &req, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveSceneResp{
		TargetID:   targetID,
		TargetName: targetName,
	})
	return
}

// BatchGetScene 获取场景
func BatchGetScene(ctx *gin.Context) {
	var req rao.GetSceneReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	s, err := scene.BatchGetByTargetID(ctx, req.TeamID, req.TargetID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetSceneResp{Scenes: s})
	return
}

// SaveFlow 保存场景流程
func SaveFlow(ctx *gin.Context) {
	var req rao.SaveFlowReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := scene.SaveFlow(ctx, &req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// GetFlow 获取场景流程
func GetFlow(ctx *gin.Context) {
	var req rao.GetFlowReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	resp, err := scene.GetFlow(ctx, req.SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, resp)
	return
}

// BatchGetFlow 批量获取场景流程
func BatchGetFlow(ctx *gin.Context) {
	var req rao.BatchGetFlowReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	flows, err := scene.BatchGetFlow(ctx, req.SceneID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.BatchGetFlowResp{Flows: flows})
	return
}

// DeleteScene 删除计划下的场景
func DeleteScene(ctx *gin.Context) {
	var req rao.DeleteSceneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := scene.DeleteScene(ctx, req.TargetID, jwt.GetUserIDByCtx(ctx)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}
