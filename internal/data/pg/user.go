package pg

import (
	"database/sql"
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
)

const usersTableName = "public.user"

func NewUsersQ(db *pgdb.DB) data.UsersQ {
	return &usersQ{
		db:        db.Clone(),
		sql:       sq.Select("public.user.*").From(usersTableName),
		sqlUpdate: sq.Update(usersTableName).Suffix("returning *"),
	}
}

type usersQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (q *usersQ) New() data.UsersQ {
	return NewUsersQ(q.db)
}

func (q *usersQ) Get() (*data.User, error) {
	var result data.User
	err := q.db.Get(&result, q.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *usersQ) Select() ([]data.User, error) {
	var result []data.User
	err := q.db.Select(&result, q.sql)
	return result, err
}

func (q *usersQ) Update(address data.User) (data.User, error) {
	var result data.User
	clauses := structs.Map(address)
	clauses["username"] = address.Username
	clauses["password_hash_hint"] = address.PasswordHashHint
	clauses["check_hash"] = address.CheckHash

	err := q.db.Get(&result, q.sqlUpdate.SetMap(clauses))

	return result, err
}

func (q *usersQ) Transaction(fn func(q data.UsersQ) error) error {
	return q.db.Transaction(func() error {
		return fn(q)
	})
}

func (q *usersQ) Insert(address data.User) (data.User, error) {
	clauses := structs.Map(address)
	clauses["username"] = address.Username
	clauses["password_hash_hint"] = address.PasswordHashHint
	clauses["check_hash"] = address.CheckHash

	var result data.User
	stmt := sq.Insert(usersTableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)

	return result, err
}

func (q *usersQ) Delete(id int64) error {
	stmt := sq.Delete(usersTableName).Where(sq.Eq{"id": id})
	err := q.db.Exec(stmt)
	return err
}

func (q *usersQ) Page(pageParams pgdb.OffsetPageParams) data.UsersQ {
	q.sql = pageParams.ApplyTo(q.sql, "id")
	return q
}

func (q *usersQ) FilterByID(ids ...int64) data.UsersQ {
	q.sql = q.sql.Where(sq.Eq{"id": ids})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"id": ids})
	return q
}

func (q *usersQ) FilterByUsername(usernames ...string) data.UsersQ {
	q.sql = q.sql.Where(sq.Eq{"username": usernames})
	return q
}
