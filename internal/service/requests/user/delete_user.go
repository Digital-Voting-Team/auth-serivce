package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteUserRequest struct {
	UserID int64 `url:"-"`
}

func NewDeleteUserRequest(r *http.Request) (DeleteUserRequest, error) {
	request := DeleteUserRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UserID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
