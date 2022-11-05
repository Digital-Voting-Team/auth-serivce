package middleware

import (
	"context"
	"errors"
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	"github.com/Digital-Voting-Team/auth-serivce/jwt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strings"
)

func BasicAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok, err := AuthDataJWT(r)
			if err != nil || !ok {
				if err == nil {
					err = errors.New("invalid credentials")
				}
				ape.Render(w, problems.BadRequest(err))
				return
			}
			ctx := context.WithValue(r.Context(), "userId", token.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func baseJWTCheck(r *http.Request) (string, *data.JWT, *data.User, error) {
	auth := r.Header.Get("Authorization")
	spitted := strings.Fields(auth)
	if len(spitted) > 1 {
		if spitted[0] != "Bearer" || len(spitted) != 2 {
			return "", nil, nil, errors.New("invalid auth string")
		}
		auth = spitted[1]
	}
	if auth == "" {
		return "", nil, nil, errors.New("empty fields username (password)")
	}
	resultJwt, err := helpers.JWTsQ(r).FilterByJWT(auth).Get()
	if err != nil {
		return "", nil, nil, errors.New("failed to get jwt by jwt string")
	}
	resultUser, err := helpers.UsersQ(r).FilterByID(resultJwt.UserID).Get()
	if err != nil {
		return "", nil, nil, errors.New("failed to get user by id")
	}
	return auth, resultJwt, resultUser, nil
}

func AuthDataJWT(r *http.Request) (*data.JWT, bool, error) {
	auth, resultJwt, resultUser, err := baseJWTCheck(r)
	if err != nil {
		return nil, false, err
	}
	ok, _, err := jwt.ParseToken(auth, resultUser.CheckHash)

	return resultJwt, ok, err
}

func AuthUserJWT(r *http.Request) (int64, bool, error) {
	auth, _, resultUser, err := baseJWTCheck(r)
	if err != nil {
		return -1, false, err
	}
	ok, userId, err := jwt.ParseToken(auth, resultUser.CheckHash)

	return userId, ok, err
}
