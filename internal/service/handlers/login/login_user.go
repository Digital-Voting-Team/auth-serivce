package handlers

import (
	"errors"
	"github.com/Digital-Voting-Team/auth-service/internal/data"
	"github.com/Digital-Voting-Team/auth-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-service/internal/service/requests/login"
	"github.com/Digital-Voting-Team/auth-service/jwt"
	"github.com/Digital-Voting-Team/auth-service/resources"
	"github.com/Digital-Voting-Team/auth-service/utils"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAuthUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to parse Login User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	checkHash := utils.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	user := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	foundUser, err := helpers.UsersQ(r).FilterByUsername(user.Username).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to find user by it's username")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if foundUser == nil {
		helpers.Log(r).WithError(err).Error("there is no such user with username: " + user.Username)
		ape.Render(w, problems.NotFound())
		return
	}

	if foundUser.PasswordHashHint != user.PasswordHashHint || foundUser.CheckHash != user.CheckHash {
		ape.Render(w, problems.BadRequest(errors.New("invalid credentials")))
		return
	}

	token, err := jwt.CreateToken(foundUser.CheckHash, foundUser.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create token")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	jwtSample := data.JWT{
		UserID: foundUser.ID,
		JWT:    token,
	}

	var resultToken data.JWT
	checkUser, err := helpers.JWTsQ(r).FilterByUserID(foundUser.ID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get jwt by the user Id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if checkUser != nil && checkUser.ID != 0 {
		resultToken, err = helpers.JWTsQ(r).FilterByUserID(foundUser.ID).Update(jwtSample)
	} else {
		resultToken, err = helpers.JWTsQ(r).Insert(jwtSample)
	}
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to insert/update new token")
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
