package handler

import (
	"time"

	"github.com/go-omnibus/omnibus"
	"github.com/go-omnibus/proof"

	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/mail"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/auth"

	"github.com/gin-gonic/gin"
)

// AuthSignup 注册
func AuthSignup(ctx *gin.Context) {
	var req rao.AuthSignupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	cnt, err := tx.WithContext(ctx).Where(tx.Email.Eq(req.Email)).Count()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}
	if cnt > 0 {
		response.ErrorWithMsg(ctx, errno.ErrYetRegister, "")
		return
	}

	u, err := auth.SignUp(ctx, req.Email, req.Password, req.Nickname)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	token, exp, err := jwt.GenerateToken(u.ID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrInvalidToken, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthSignupResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
	return
}

// AuthLogin 登录
func AuthLogin(ctx *gin.Context) {
	var req rao.AuthLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	u, err := auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrAuthFailed, err.Error())
		return
	}

	d := 7 * 24 * time.Hour
	if req.IsAutoLogin {
		d = 30 * 24 * time.Hour
	}

	token, exp, err := jwt.GenerateTokenByTime(u.ID, d)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrInvalidToken, err.Error())
		return
	}

	if err := auth.UpdateLoginTime(ctx, u.ID); err != nil {
		proof.Errorf("update login time err %s", err)
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
	return
}

// RefreshToken 续期
func RefreshToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	token, exp, err := jwt.RefreshToken(tokenString)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
	return
}

func UpdatePassword(ctx *gin.Context) {
	var req rao.UpdatePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if req.CurrentPassword == req.NewPassword {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}

	tx := dal.GetQuery().User
	u, err := tx.WithContext(ctx).Where(tx.ID.Eq(jwt.GetUserIDByCtx(ctx))).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}
	if err := omnibus.CompareBcryptHashAndPassword(u.Password, req.CurrentPassword); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}

	hashedPassword, err := omnibus.GenerateBcryptFromPassword(req.NewPassword)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}
	if _, err := tx.WithContext(ctx).Where(tx.ID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(tx.Password, hashedPassword); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}

	response.Success(ctx)
	return
}

func UpdateNickname(ctx *gin.Context) {
	var req rao.UpdateNicknameReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	if _, err := tx.WithContext(ctx).Where(tx.ID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(tx.Nickname, req.Nickname); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "nickname failed")
		return
	}

	response.Success(ctx)
	return
}

func UpdateAvatar(ctx *gin.Context) {
	var req rao.UpdateAvatarReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	if _, err := tx.WithContext(ctx).Where(tx.ID.Eq(jwt.GetUserIDByCtx(ctx))).UpdateColumn(tx.Avatar, req.AvatarURL); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, "password = new password")
		return
	}

	response.Success(ctx)
	return
}

// SetUserSettings 设置用户配置
func SetUserSettings(ctx *gin.Context) {
	var req rao.SetUserSettingsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	if err := auth.SetUserSettings(ctx, jwt.GetUserIDByCtx(ctx), &req.UserSettings); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

// GetUserSettings 获取用户配置
func GetUserSettings(ctx *gin.Context) {
	settings, err := auth.GetUserSettings(ctx, jwt.GetUserIDByCtx(ctx))
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.GetUserSettingsResp{
		UserSettings: settings,
	})
	return
}

func AuthForgetPassword(ctx *gin.Context) {
	var req rao.ForgetPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	u, err := tx.WithContext(ctx).Where(tx.Email.Eq(req.Email)).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	if err := mail.SendForgetPasswordEmail(ctx, req.Email, u.ID); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrRPCFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func AuthResetPassword(ctx *gin.Context) {
	var req rao.ResetPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	hashedPassword, err := omnibus.GenerateBcryptFromPassword(req.NewPassword)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	if _, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.U)).UpdateColumn(tx.Password, hashedPassword); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrMysqlFailed, err.Error())
		return
	}

	response.Success(ctx)
	return
}

func VerifyPassword(ctx *gin.Context) {
	var req rao.VerifyPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	tx := dal.GetQuery().User
	u, err := tx.WithContext(ctx).Where(tx.ID.Eq(jwt.GetUserIDByCtx(ctx))).First()
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}

	err = omnibus.CompareBcryptHashAndPassword(u.Password, req.Password)

	response.SuccessWithData(ctx, rao.VerifyPasswordResp{IsMatch: err == nil})
	return
}
