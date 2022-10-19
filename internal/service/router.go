package service

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/data/pg"
	login "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/login"
	reg "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/register"
	user "github.com/Digital-Voting-Team/auth-serivce/internal/service/handlers/user"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/middleware"

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
		r.Post("/login", reg.RegisterUser)
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
