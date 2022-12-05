package handlers

import (
	"errors"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/login"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAuthUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to parse Login User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}
	username := request.Data.Attributes.Username

	foundUser, err := helpers.UsersQ(r).FilterByUsername(username).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to find user by it's username")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if foundUser == nil {
		helpers.Log(r).WithError(err).Error("there is no such user with username: " + username)
		ape.Render(w, problems.NotFound())
		return
	}

	jwt, err := helpers.JWTsQ(r).FilterByUserID(foundUser.ID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get jwt by the user Id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if jwt == nil {
		helpers.Log(r).WithError(err).Error("there is no such user with username: " + username)
		ape.Render(w, problems.BadRequest(validation.Errors{"User": errors.New("already signed out")}))
		return
	}

	err = helpers.JWTsQ(r).Delete(jwt.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete jwt")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
