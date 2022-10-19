package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/auth-serivce/internal/service/helpers"
	"github.com/Digital-Voting-Team/auth-serivce/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateUserRequest struct {
	UserID int64 `url:"-" json:"-"`
	Data   resources.User
}

func NewUpdateUserRequest(r *http.Request) (UpdateUserRequest, error) {
	request := UpdateUserRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UserID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateUserRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/username": validation.Validate(&r.Data.Attributes.Username, validation.Required,
			validation.Length(3, 30)),
		"/data/attributes/password": validation.Validate(&r.Data.Attributes.Password, validation.Required,
			validation.Length(3, 120)),
	}).Filter()
}
