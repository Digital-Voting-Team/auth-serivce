package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type UserClaims struct {
	Username string
	jwt.StandardClaims
}

func CreateToken(username, password string) (string, error) {
	// Create the Claims
	claims := UserClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("\n\n\nPASSWORD: %v", password)
	return token.SignedString([]byte(password))
}

func ParseToken(tokenString, password string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(password), nil
	})
	fmt.Printf("\n\n\nPASSWORD: %v", password)

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			return false, errors.New("token expired")
		}
		return true, nil
	} else {
		return false, err
	}
}
