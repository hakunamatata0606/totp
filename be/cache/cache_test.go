package cache_test

import (
	"errors"
	"example/totp/cache"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInMemCache(t *testing.T) {
	c := cache.NewInMemCache(2000 * time.Millisecond)
	key1 := "key1"
	value1 := 1
	key2 := "key2"
	value2 := "value2"

	c.Set(key1, value1)
	c.Set(key2, value2)

	v1, err := c.Get(key1)
	require.Nil(t, err)
	v1Int, ok := v1.(int)
	require.True(t, ok)
	require.Equal(t, value1, v1Int)

	v2, err := c.Get(key2)
	require.Nil(t, err)
	v2Str, ok := v2.(string)
	require.True(t, ok)
	require.Equal(t, value2, v2Str)

	_, err = c.Get("dummy")
	require.NotNil(t, err)
	require.Equal(t, errors.New("key not found"), err)

	time.Sleep(3000 * time.Millisecond)

	v1, err = c.Get(key1)
	require.Nil(t, v1)
	require.NotNil(t, err)
	require.Equal(t, errors.New("key expired"), err)
	v2, err = c.Get(key2)
	require.Nil(t, v2)
	require.NotNil(t, err)
	require.Equal(t, errors.New("key expired"), err)
}
