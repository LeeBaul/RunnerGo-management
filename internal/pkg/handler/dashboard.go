package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardDefault(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
