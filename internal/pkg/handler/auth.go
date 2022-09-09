package handler

import (
	"time"

	"github.com/go-omnibus/proof"

	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
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

	u, err := auth.SignUp(ctx, req.Email, req.Password, req.Nickname)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrInvalidToken, err.Error())
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

	d := 2 * time.Hour
	if req.IsAutoLogin {
		d = 30 * 24 * time.Hour
	}

	token, exp, err := jwt.GenerateTokenByTime(u.ID, d)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.ErrInvalidToken, err.Error())
		return
	}

	// first login
	var isAPIPostUser bool
	if u.LastLoginAt.IsZero() {
		i, err := auth.IsAPIPostUser(ctx, u.Email)
		if err != nil {
			proof.Errorf("is apipost user err %s", err)
		}
		isAPIPostUser = i
	}

	if err := auth.UpdateLoginTime(ctx, u.ID); err != nil {
		proof.Errorf("update login time err %s", err)
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
		IsAPIPostUser: isAPIPostUser,
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

func AuthSendMailVerify(ctx *gin.Context) {
	var req rao.AuthSendMailVerifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ErrParam, err.Error())
		return
	}
}

func AuthResetPassword(ctx *gin.Context) {
}
