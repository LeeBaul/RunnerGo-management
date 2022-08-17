package handler

import (
	"kp-management/internal/pkg/biz/errno"
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
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	token, err := auth.GenerateJWT(u)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthSignupResp{
		User:  nil,
		Token: token,
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
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	token, err := auth.GenerateJWT(u)
	if err != nil {
		response.ErrorWithMsg(ctx, errno.MysqlOperFailed, err.Error())
		return
	}

	response.SuccessWithData(ctx, rao.AuthLoginResp{
		User:  nil,
		Token: token,
	})
}
