package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type JWTsQ interface {
	New() JWTsQ

	Get() (*JWT, error)
	Select() ([]JWT, error)

	Transaction(fn func(q JWTsQ) error) error

	Insert(JWT) (JWT, error)
	Update(JWT) (JWT, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) JWTsQ

	FilterByID(ids ...int64) JWTsQ

	JoinUser() JWTsQ
}

type JWT struct {
	ID     int64  `db:"id" structs:"-"`
	UserID int64  `db:"user_id" structs:"user_id"`
	JWT    string `db:"jwt" structs:"jwt"`
}
