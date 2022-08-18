package middleware

import (
	"kp-management/internal/pkg/biz/errno"
	"kp-management/internal/pkg/biz/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"code": errno.MustLogin, "message": "must login"})
			c.Abort()
			return
		}

		userID, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": errno.MustLogin, "message": "must login"})
			c.Abort()
			return
		} else if userID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": errno.MustLogin, "message": "must login"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		//c.Set(UserName, claims.Username)

		c.Next()
	}
}

//func GetUid(c *gin.Context) int64 {
//	return c.GetInt64(UserId)
//}

//func GetUserName(c *gin.Context) string {
//	return c.GetString(UserName)
//}
