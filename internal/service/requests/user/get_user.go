package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetUserRequest struct {
	UserID int64 `url:"-"`
}

func NewGetUserRequest(r *http.Request) (GetUserRequest, error) {
	request := GetUserRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UserID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
