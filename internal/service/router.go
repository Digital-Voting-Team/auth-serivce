package service

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/data/pg"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/jwt"
	login "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/login"
	reg "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/register"
	user "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/user"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	"github.com/Digital-Voting-Team/auth-serivce/middleware"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	log := s.log.WithFields(map[string]interface{}{
		"service": "auth-service-api",
	})

	r.Use(
		ape.RecoverMiddleware(log),
		ape.LoganMiddleware(log),
		ape.CtxMiddleware(
			helpers.CtxLog(log),
			helpers.CtxUsersQ(pg.NewUsersQ(s.db)),
			helpers.CtxJWTsQ(pg.NewJWTsQ(s.db)),
		),
	)

	r.Group(func(r chi.Router) {
		r.Get("/login", login.LoginUser)
		r.Get("/jwt/login", jwt.LoginJWT)
		r.Post("/register", reg.RegisterUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.BasicAuth())
		r.Route("/integrations/auth-service", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/", user.CreateUser)
				r.Get("/", user.GetUserList)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", user.GetUser)
					r.Put("/", user.UpdateUser)
					r.Delete("/", user.DeleteUser)
				})
			})
		})
	})

	return r
}
