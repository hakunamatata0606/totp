package repository_test

import (
	"context"
	"example/totp/appstate"
	"example/totp/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository(t *testing.T) {
	state := appstate.GetAppState()
	repo := repository.New(state.Db)
	username1 := "bao"
	username2 := "dummy"

	user, err := repo.GetUser(context.Background(), &username1)
	require.Nil(t, err)
	require.NotNil(t, user)
	require.Equal(t, username1, user.Username)

	user, err = repo.GetUser(context.Background(), &username2)
	require.NotNil(t, err)
	require.Nil(t, user)
}
