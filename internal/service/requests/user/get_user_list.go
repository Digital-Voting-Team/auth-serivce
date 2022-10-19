package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetUserListRequest struct {
	pgdb.OffsetPageParams
	FilterUsername []string `filter:"username"`
}

func NewGetUserListRequest(r *http.Request) (GetUserListRequest, error) {
	var request GetUserListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
