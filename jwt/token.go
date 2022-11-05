package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type UserClaims struct {
	UserId int64
	jwt.StandardClaims
}

func CreateToken(username, password string, userId int64) (string, error) {
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

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return false, -1, errors.New("token expired")
		}
		return true, claims.UserId, nil
	} else {
		return false, -1, err
	}
}
