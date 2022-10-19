package handlers

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/register"
	"github.com/Digital-Voting-Team/auth-serivce/jwt"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	utils2 "github.com/Digital-Voting-Team/auth-serivce/utils"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRegisterUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultUser data.User

	checkHash := utils2.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	user := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils2.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	if findUsername, err := helpers.UsersQ(r).FilterByUsername(user.Username).Get(); findUsername != nil {
		helpers.Log(r).WithError(err).Error("username already used")
		ape.Render(w, problems.Conflict())
		return
	}

	resultUser, err = helpers.UsersQ(r).Insert(user)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to register user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	token, err := jwt.CreateToken(resultUser.Username, resultUser.CheckHash)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.JwtResponse{
		Data: resources.Jwt{
			Key: resources.NewKeyInt64(resultUser.ID, resources.JWT),
			Attributes: resources.JwtAttributes{
				Jwt: token,
			},
			Relationships: resources.JwtRelationships{
				User: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultUser.ID, 10),
						Type: resources.USER,
					},
				},
			},
		},
	}
	ape.Render(w, result)
}
