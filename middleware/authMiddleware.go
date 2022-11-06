package middleware

import (
	"context"
	"errors"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func BasicAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok, err := helpers.AuthJWT(r)
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
