package handlers

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/user"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	utils2 "github.com/Digital-Voting-Team/auth-serivce/utils"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to Parse Update User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	user, err := helpers.UsersQ(r).FilterByID(request.UserID).Get()
	if user == nil {
		helpers.Log(r).Error("user == nil (update)")
		ape.Render(w, problems.NotFound())
		return
	}

	checkHash := utils2.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	newUser := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils2.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	userId := r.Context().Value("userId").(int64)
	if userId != user.ID {
		helpers.Log(r).Info("invalid user")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	var resultUser data.User
	resultUser, err = helpers.UsersQ(r).FilterByID(user.ID).Update(newUser)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.UserResponse{
		Data: resources.User{
			Key: resources.NewKeyInt64(resultUser.ID, resources.USER),
			Attributes: resources.UserAttributes{
				Username: resultUser.Username,
			},
		},
	}

	ape.Render(w, result)
}
