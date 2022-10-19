package handlers

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/user"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := helpers.UsersQ(r).FilterByID(request.UserID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get user from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if user == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.UserResponse{
		Data: resources.User{
			Key: resources.NewKeyInt64(user.ID, resources.USER),
			Attributes: resources.UserAttributes{
				Username: user.Username,
			},
		},
	}

	ape.Render(w, result)
}
