package repository

import (
	"context"
	db "example/totp/db/sqlc"
	"log"
)

type User struct {
	Username string
	Secret   []byte
}

type RepositoryIf interface {
	GetUser(ctx context.Context, username *string) (*User, error)
}

type repositoryImpl struct {
	queries db.Queries
}

func New(dbtx db.DBTX) RepositoryIf {
	impl := &repositoryImpl{
		queries: *db.New(dbtx),
	}
	return impl
}

func (repo *repositoryImpl) GetUser(ctx context.Context, username *string) (*User, error) {
	userDb, err := repo.queries.GetUser(ctx, *username)
	// handle error ?
	if err != nil {
		log.Println("Failed to get user: ", err)
		return nil, err
	}
	user := &User{
		Username: userDb.Username,
		Secret:   userDb.Secret,
	}
	return user, nil
}
