package handler

import (
	"github.com/gin-gonic/gin"

	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/scene"
)

// SaveScene 创建/修改场景
func SaveScene(ctx *gin.Context) {
	var req rao.SaveSceneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	targetID, err := scene.Save(ctx, &req, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveSceneResp{TargetID: targetID})
	return
}

// BatchGetScene 获取场景
func BatchGetScene(ctx *gin.Context) {
	var req rao.GetSceneReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	s, err := scene.BatchGetByTargetID(ctx, req.TeamID, req.TargetID, req.Source)
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
