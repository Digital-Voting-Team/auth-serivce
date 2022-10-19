package pg

import (
	"database/sql"
	"fmt"
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const jwtsTableName = "public.jwt"

func NewJWTsQ(db *pgdb.DB) data.JWTsQ {
	return &jwtsQ{
		db:        db.Clone(),
		sql:       sq.Select("jwt.*").From(jwtsTableName),
		sqlUpdate: sq.Update(jwtsTableName).Suffix("returning *"),
	}
}

type jwtsQ struct {
	db        *pgdb.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func (c *jwtsQ) New() data.JWTsQ {
	return NewJWTsQ(c.db)
}

func (c *jwtsQ) Get() (*data.JWT, error) {
	var result data.JWT
	err := c.db.Get(&result, c.sql)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (c *jwtsQ) Select() ([]data.JWT, error) {
	var result []data.JWT
	err := c.db.Select(&result, c.sql)
	return result, err
}

func (c *jwtsQ) Transaction(fn func(q data.JWTsQ) error) error {
	return c.db.Transaction(func() error {
		return fn(c)
	})
}

func (c *jwtsQ) Insert(jwt data.JWT) (data.JWT, error) {
	clauses := structs.Map(jwt)
	clauses["person_id"] = jwt.UserID
	clauses["jwt"] = jwt.JWT
	clauses["registration_date"] = jwt.ExpirationTime

	var result data.JWT
	stmt := sq.Insert(jwtsTableName).SetMap(clauses).Suffix("returning *")
	err := c.db.Get(&result, stmt)

	return result, err
}

func (c *jwtsQ) Update(jwt data.JWT) (data.JWT, error) {
	var result data.JWT
	clauses := structs.Map(jwt)
	clauses["person_id"] = jwt.UserID
	clauses["jwt"] = jwt.JWT
	clauses["registration_date"] = jwt.ExpirationTime

	err := c.db.Get(&result, c.sqlUpdate.SetMap(clauses))
	return result, err
}

func (c *jwtsQ) Delete(id int64) error {
	stmt := sq.Delete(jwtsTableName).Where(sq.Eq{"id": id})
	err := c.db.Exec(stmt)
	return err
}

func (c *jwtsQ) Page(pageParams pgdb.OffsetPageParams) data.JWTsQ {
	c.sql = pageParams.ApplyTo(c.sql, "id")
	return c
}

func (c *jwtsQ) FilterByID(ids ...int64) data.JWTsQ {
	c.sql = c.sql.Where(sq.Eq{"id": ids})
	c.sqlUpdate = c.sqlUpdate.Where(sq.Eq{"id": ids})
	return c
}

func (c *jwtsQ) JoinUser() data.JWTsQ {
	stmt := fmt.Sprintf("%s as jwt on public.user.id = jwt.user_id",
		jwtsTableName)
	c.sql = c.sql.Join(stmt)
	return c
}
