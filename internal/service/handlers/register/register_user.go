package handlers

import (
	"github.com/Digital-Voting-Team/auth-service/internal/data"
	"github.com/Digital-Voting-Team/auth-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-service/internal/service/requests/register"
	"github.com/Digital-Voting-Team/auth-service/jwt"
	"github.com/Digital-Voting-Team/auth-service/resources"
	"github.com/Digital-Voting-Team/auth-service/utils"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRegisterUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to parse Register User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	var resultUser data.User

	checkHash := utils.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	user := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	findUser, err := helpers.UsersQ(r).FilterByUsername(user.Username).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to find user by it's username")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if findUser != nil {
		helpers.Log(r).Error("username already used")
		ape.Render(w, problems.Conflict())
		return
	}

	resultUser, err = helpers.UsersQ(r).Insert(user)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to insert user (register)")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	token, err := jwt.CreateToken(resultUser.CheckHash, resultUser.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	jwtSample := data.JWT{
		UserID: resultUser.ID,
		JWT:    token,
	}

	resultToken, err := helpers.JWTsQ(r).Insert(jwtSample)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to insert new token (register)")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.JwtResponse{
		Data: resources.Jwt{
			Key: resources.NewKeyInt64(resultToken.ID, resources.JWT),
			Attributes: resources.JwtAttributes{
				Jwt: resultToken.JWT,
			},
			Relationships: resources.JwtRelationships{
				User: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(resultToken.UserID, 10),
						Type: resources.USER,
					},
				},
			},
		},
	}

	ape.Render(w, result)
}
