package helpers

import (
	"context"
	"github.com/Digital-Voting-Team/auth-serivce/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	userCtxKey
	jwtCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxUsersQ(entry data.UsersQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userCtxKey, entry)
	}
}

func UsersQ(r *http.Request) data.UsersQ {
	return r.Context().Value(userCtxKey).(data.UsersQ).New()
}

func CtxJWTsQ(entry data.JWTsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, jwtCtxKey, entry)
	}
}

func JWTsQ(r *http.Request) data.JWTsQ {
	return r.Context().Value(jwtCtxKey).(data.JWTsQ).New()
}
