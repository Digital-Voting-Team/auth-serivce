package middleware

import (
	"errors"
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
			ok, err := AuthJWT(r)
			if err != nil || !ok {
				if err == nil {
					err = errors.New("invalid credentials")
				}
				ape.Render(w, problems.BadRequest(err))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthJWT(r *http.Request) (bool, error) {
	auth := r.Header.Get("Authorization")
	spitted := strings.Fields(auth)
	if len(spitted) > 1 {
		if spitted[0] != "Bearer" || len(spitted) != 2 {
			return false, errors.New("invalid auth string")
		}
		auth = spitted[1]
	}
	if auth == "" {
		return false, errors.New("empty fields username (password)")
	}
	resultJwt, err := helpers.JWTsQ(r).FilterByJWT(auth).Get()
	if err != nil {
		return false, errors.New("failed to get jwt by jwt string")
	}
	resultUser, err := helpers.UsersQ(r).FilterByID(resultJwt.UserID).Get()
	if err != nil {
		return false, errors.New("failed to get user by id")
	}

	return jwt.ParseToken(auth, resultUser.CheckHash)
}
