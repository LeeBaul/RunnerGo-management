package handler

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"kp-management/internal/pkg/biz/response"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/logic/auth"

	"github.com/gin-gonic/gin"
)

func AuthSignup(ctx *gin.Context) {
	var req rao.AuthSignupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	u, err := auth.SignUp(ctx, req.Email, req.Password, req.Nickname)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.InvalidToken, err.Error())
		return
	}

	token, exp, err := jwt.GenerateToken(u.ID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.InvalidToken, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthSignupResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
}

func AuthLogin(ctx *gin.Context) {
	var req rao.AuthLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(ctx, errno.ParamError, err.Error())
		return
	}

	u, err := auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.InvalidToken, err.Error())
		return
	}

	token, exp, err := jwt.GenerateToken(u.ID)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.InvalidToken, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
}

func AuthRefresh(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	token, exp, err := jwt.RefreshToken(tokenString)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		Token:         token,
		ExpireTimeSec: exp.Unix(),
	})
}
