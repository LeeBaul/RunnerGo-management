package jwt

import (
	"fmt"
	"kp-management/internal/pkg/conf"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID int64) (string, time.Time, error) {
	now := time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC)
	exp := now.Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iss":     conf.Conf.JWT.Issuer,
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     exp.Unix(),
	})
	tokenString, err := token.SignedString(conf.Conf.JWT.Secret)

	return tokenString, exp, err
}

func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(conf.Conf.JWT.Secret), nil
	})

	if err != nil {
		return 0, jwt.ErrHashUnavailable
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"]; ok {
			if id, ok := userID.(int64); ok {
				return id, nil
			}
		}
	}

	return 0, jwt.ErrTokenInvalidClaims
}
