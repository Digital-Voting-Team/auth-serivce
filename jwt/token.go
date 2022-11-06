package jwt

import (
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"time"
)

type UserClaims struct {
	UserId int64
	jwt.StandardClaims
}

func CreateToken(password string, userId int64) (string, error) {
	claims := UserClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(password))
}

func ParseToken(tokenString, password string) (bool, int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})
	if err != nil {
		return false, -1, errors.Wrap(err, "failed to parse token with claims")
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return false, -1, errors.New("token expired")
		}
		return true, claims.UserId, nil
	} else {
		return false, -1, errors.Wrap(err, "failed to cast token claims and validate token")
	}
}
