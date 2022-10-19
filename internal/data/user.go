package data

import "gitlab.com/distributed_lab/kit/pgdb"

type UsersQ interface {
	New() UsersQ

	Get() (*User, error)
	Select() ([]User, error)

	Transaction(fn func(q UsersQ) error) error

	Insert(User) (User, error)
	Update(User) (User, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) UsersQ

	FilterByID(ids ...int64) UsersQ
	FilterByUsername(usernames ...string) UsersQ
}

type User struct {
	ID               int64  `db:"id" structs:"-"`
	Username         string `db:"username" structs:"username"`
	PasswordHashHint string `db:"password_hash_hint" structs:"-"`
	CheckHash        string `db:"check_hash" structs:"-"`
}
