package handlers

import (
	"errors"
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/login"
	"github.com/Digital-Voting-Team/auth-serivce/jwt"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	utils2 "github.com/Digital-Voting-Team/auth-serivce/utils"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewLoginUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	checkHash := utils2.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	user := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils2.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	foundUser, err := helpers.UsersQ(r).FilterByUsername(user.Username).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to Login user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if foundUser.PasswordHashHint != user.PasswordHashHint || foundUser.CheckHash != user.CheckHash {
		ape.Render(w, problems.BadRequest(errors.New("invalid credentials")))
		return
	}

	token, err := jwt.CreateToken(foundUser.Username, foundUser.CheckHash)
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

	var includes resources.Included
	includes.Add(&resources.User{
		Key: resources.NewKeyInt64(foundUser.ID, resources.USER),
		Attributes: resources.UserAttributes{
			Username: foundUser.Username,
		},
	})

	result := resources.JwtResponse{
		Data: resources.Jwt{
			Key: resources.NewKeyInt64(resultToken.ID, resources.JWT),
			Attributes: resources.JwtAttributes{
				Jwt: resultToken.JWT,
			},
			Relationships: resources.JwtRelationships{
				User: resources.Relation{
					Data: &resources.Key{
						ID:   strconv.FormatInt(foundUser.ID, 10),
						Type: resources.USER,
					},
				},
			},
		},
		Included: includes,
	}

	ape.Render(w, result)
}
