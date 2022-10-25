package handler

import (
	"fmt"
	"time"

	"github.com/go-omnibus/omnibus"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/team"

	"github.com/gin-gonic/gin"
)

func SaveTeam(ctx *gin.Context) {
	var req rao.SaveTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	teamID, err := team.SaveTeam(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx), req.Name)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.SaveTeamResp{
		TeamID: teamID,
	})
	return
}

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

func GetTeamRole(ctx *gin.Context) {
	var req rao.GetTeamRoleReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().UserTeam
	ut, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(req.TeamID), tx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetTeamRoleResp{
		RoleID: ut.RoleID,
	})
	return
}

func GetInviteMemberURL(ctx *gin.Context) {
	var req rao.GetInviteMemberURLReq
	if err := ctx.ShouldBind(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().UserTeam
	_, err := tx.WithContext(ctx).Where(tx.UserID.Eq(jwt.GetUserIDByCtx(ctx)), tx.RoleID.In(consts.RoleTypeAdmin, consts.RoleTypeOwner)).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	k := fmt.Sprintf("invite:url:%d:%d", req.TeamID, req.RoleID)
	_, err = dal.GetRDB().Set(ctx, k, fmt.Sprintf("%d", jwt.GetUserIDByCtx(ctx)), 5*time.Minute).Result()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrRedisFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, &rao.GetInviteMemberURLResp{
		URL:     fmt.Sprintf("%s#/login?team_id=%d&role_id=%d", conf.Conf.Base.Domain, req.TeamID, req.RoleID),
		Expired: time.Now().Add(time.Hour * 24).Unix(),
	})
	return
}

func CheckInviteMemberURL(ctx *gin.Context) {
	var req rao.CheckInviteMemberURLReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	k := fmt.Sprintf("invite:url:%d:%d", req.TeamID, req.RoleID)
	inviteUserID, err := dal.GetRDB().Get(ctx, k).Result()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrURLExpired, err.Error())
		return
	}
	if inviteUserID == "" {
		response.ErrorWithMsg(ctx, errno.ErrURLExpired, "")
		return
	}

	tx := dal.GetQuery().UserTeam
	cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(req.TeamID), tx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).Count()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	if cnt > 0 {
		response.ErrorWithMsg(ctx, errno.ErrExistsTeam, "")
		return
	}

	err = tx.WithContext(ctx).Create(&model.UserTeam{
		UserID:       jwt.GetUserIDByCtx(ctx),
		TeamID:       req.TeamID,
		RoleID:       req.RoleID,
		InviteUserID: omnibus.DefiniteInt64(inviteUserID),
	})
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	sx := dal.GetQuery().Setting
	_, err = sx.WithContext(ctx).Where(sx.UserID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(sx.TeamID, req.TeamID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// InviteMember 邀请成员
func InviteMember(ctx *gin.Context) {
	var req rao.InviteMemberReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	resp, err := team.InviteMember(ctx, jwt.GetUserIDByCtx(ctx), req.TeamID, req.Members)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, resp)
	return
}

func RoleUser(ctx *gin.Context) {
	var req rao.RoleUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.RoleUser(ctx, req.TeamID, req.UserID, req.RoleID); err != nil {
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

	if err := team.RemoveMember(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx), req.MemberID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func QuitTeam(ctx *gin.Context) {
	var req rao.QuitTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.QuitTeam(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func DisbandTeam(ctx *gin.Context) {
	var req rao.DisbandTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := team.DisbandTeam(ctx, req.TeamID, jwt.GetUserIDByCtx(ctx)); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func TransferTeam(ctx *gin.Context) {
	var req rao.TransferTeamReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
}
