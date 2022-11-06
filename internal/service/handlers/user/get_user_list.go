package handlers

import (
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/auth-serivce/internal/service/requests/user"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(int64)
	if userId != 1 {
		helpers.Log(r).Info("insufficient user permissions")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	request, err := requests.NewGetUserListRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("failed to Parse Get User List Request")
		ape.Render(w, problems.BadRequest(err))
		return
	}

	useresQ := helpers.UsersQ(r)
	applyFilters(useresQ, request)
	user, err := useresQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get user")
		ape.Render(w, problems.InternalError())
		return
	}

	response := resources.UserListResponse{
		Data:  newUseresList(user),
		Links: helpers.GetOffsetLinks(r, request.OffsetPageParams),
	}
	ape.Render(w, response)
}

func applyFilters(q data.UsersQ, request requests.GetUserListRequest) {
	q.Page(request.OffsetPageParams)

	if len(request.FilterUsername) > 0 {
		q.FilterByUsername(request.FilterUsername...)
	}
}

func newUseresList(useres []data.User) []resources.User {
	result := make([]resources.User, len(useres))
	for i, user := range useres {
		result[i] = resources.User{
			Key: resources.NewKeyInt64(user.ID, resources.USER),
			Attributes: resources.UserAttributes{
				Username: user.Username,
			},
		}
	}
	return result
}
