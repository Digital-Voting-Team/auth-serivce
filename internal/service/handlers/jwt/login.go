package jwt

import (
	"errors"
	"github.com/Digital-Voting-Team/auth-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/auth-service/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func LoginJWT(w http.ResponseWriter, r *http.Request) {
	token, ok, err := helpers.AuthJWT(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to Auth JWT in handler")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	if !ok {
		helpers.Log(r).Error("JWT is invalid in handler")
		ape.Render(w, problems.BadRequest(errors.New("invalid JWT")))
		return
	}

	result := resources.JwtResponse{
		Data: resources.Jwt{
			Key: resources.NewKeyInt64(token.ID, resources.JWT),
			Attributes: resources.JwtAttributes{
				Jwt: token.JWT,
			},
			Relationships: resources.JwtRelationships{
				User: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(token.UserID, 10),
						Type: resources.USER,
					},
				},
			},
		},
	}

	ape.Render(w, result)
}
