package jwttoken_test

import (
	jwttoken "example/totp/jwt_token"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJwtTokenManager(t *testing.T) {
	tm := jwttoken.New([]byte("aloha"), 3*time.Second)
	claims := &jwttoken.Claims{
		"key1": "value1",
		"key2": "value2",
	}

	token, err := tm.GenerateToken(claims)
	require.Nil(t, err)
	time.Sleep(1 * time.Second)

	claims1, err := tm.VerifyToken(&token)
	require.Nil(t, err)
	require.Equal(t, (*claims)["key1"], (*claims1)["key1"])
	require.Equal(t, (*claims)["key2"], (*claims1)["key2"])

	time.Sleep(3 * time.Second)
	_, err = tm.VerifyToken(&token)
	require.NotNil(t, err)
}
