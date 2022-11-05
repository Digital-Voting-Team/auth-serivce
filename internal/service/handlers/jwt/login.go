package jwt

import (
	"github.com/Digital-Voting-Team/auth-serivce/middleware"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func LoginJWT(w http.ResponseWriter, r *http.Request) {
	token, ok, err := middleware.AuthDataJWT(r)
	if !ok || err != nil {
		ape.Render(w, problems.BadRequest(err))
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
