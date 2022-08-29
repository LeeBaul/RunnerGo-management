package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/team"

	"github.com/gin-gonic/gin"
)

// ListTeam 团队列表
func ListTeam(ctx *gin.Context) {
	teams, err := team.ListByUserID(ctx, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListTeamResp{Teams: teams})
	return
}

// TeamMembers 团队成员列表
func TeamMembers(ctx *gin.Context) {
	var req rao.ListMembersReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	members, err := team.ListMembersByTeamID(ctx, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.ListMembersResp{
		Members: members,
	})
	return
}

// InviteMember 邀请成员
func InviteMember(ctx *gin.Context) {
	var req rao.InviteMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.InviteMember(ctx, req.TeamID, req.MemberEmail); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// RemoveMember 移除成员
func RemoveMember(ctx *gin.Context) {
	var req rao.RemoveMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.RemoveMember(ctx, req.TeamID, req.MemberID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}
