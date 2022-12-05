package helpers

import (
	"github.com/Digital-Voting-Team/auth-service/internal/data"
	"github.com/Digital-Voting-Team/auth-service/jwt"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strings"
)

func AuthJWT(r *http.Request) (*data.JWT, bool, error) {
	auth := r.Header.Get("Authorization")
	spitted := strings.Fields(auth)
	if len(spitted) > 1 {
		if spitted[0] != "Bearer" || len(spitted) != 2 {
			return nil, false, errors.New("invalid auth string")
		}
		auth = spitted[1]
	}
	if auth == "" {
		return nil, false, errors.New("auth string is empty")
	}
	resultJwt, err := JWTsQ(r).FilterByJWT(auth).Get()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get jwt by jwt string")
	}
	if resultJwt == nil {
		return nil, false, errors.New("there is no resultJwt (nil) for such jwt string")
	}

	resultUser, err := UsersQ(r).FilterByID(resultJwt.UserID).Get()
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get user by id")
	}
	if resultUser == nil {
		return nil, false, errors.New("there is no resultUser (nil) for such user Id")
	}

	ok, _, err := jwt.ParseToken(auth, resultUser.CheckHash)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to Parse Token")
	}

	return resultJwt, ok, nil
}
