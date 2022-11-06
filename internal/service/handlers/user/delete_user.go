package handlers

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/user"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	if userId != 1 {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewDeleteUserRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to Parse Delete User Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	user, err := helpers.UsersQ(r).FilterByID(request.UserID).Get()
	if user == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	err = helpers.UsersQ(r).Delete(request.UserID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
