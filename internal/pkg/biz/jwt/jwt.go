package jwt

import (
	"fmt"
	"kp-management/internal/pkg/conf"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID int64) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(24 * time.Hour * 365)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iss":     conf.Conf.JWT.Issuer,
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(conf.Conf.JWT.Secret))

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

			if i, ok := userID.(float64); ok {
				return int64(i), nil
			}

		}
	}

	return 0, jwt.ErrTokenInvalidClaims
}

func RefreshToken(tokenString string) (string, time.Time, error) {
	userID, err := ParseToken(tokenString)
	if err != nil {
		return "", time.Now(), err
	}

	return GenerateToken(userID)
}
