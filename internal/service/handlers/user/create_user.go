package handlers

import (
	"github.com/Digital-Voting-Team/auth-service/internal/data"
	"github.com/Digital-Voting-Team/auth-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-service/internal/service/requests/user"
	"github.com/Digital-Voting-Team/auth-service/resources"
	utils2 "github.com/Digital-Voting-Team/auth-service/utils"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	if userId != 1 {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewCreateUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to parse Create User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	var resultUser data.User

	checkHash := utils2.HashString(request.Data.Attributes.Username + request.Data.Attributes.Password + "CSCA")
	user := data.User{
		Username:         request.Data.Attributes.Username,
		PasswordHashHint: utils2.Hint(request.Data.Attributes.Password, 4),
		CheckHash:        checkHash,
	}

	resultUser, err = helpers.UsersQ(r).Insert(user)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create user")
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
