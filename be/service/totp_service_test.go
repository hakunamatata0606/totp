package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"example/totp/otp"
	"example/totp/repository"
	"example/totp/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mockRepo struct{}

func (repo *mockRepo) GetUser(ctx context.Context, username *string) (*repository.User, error) {
	if *username == "bao" {
		return &repository.User{
			Username: "bao",
			Secret:   []byte("aloha"),
		}, nil
	}
	return nil, errors.New("not found")
}

func TestTotpHandler(t *testing.T) {
	secret := []byte("this is secret")
	interval := 5

	r := gin.Default()

	om := otp.New(secret, int64(interval), 6)

	r.POST("/otp", service.TotpHandler(&mockRepo{}, om))

	otp := om.GenerateOtp([]byte("aloha"))
	user := service.OtpPayload{
		Username: "bao",
		Otp:      otp,
	}

	payload, err := json.Marshal(user)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	require.Nil(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	time.Sleep(time.Duration(interval) * time.Second)

	w1 := httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	require.Nil(t, err)
	r.ServeHTTP(w1, req)
	require.Equal(t, http.StatusBadRequest, w1.Result().StatusCode)
}

func TestTotpLibcHandler(t *testing.T) {
	secret := []byte("this is secret")
	interval := 5

	r := gin.Default()

	om := otp.NewLibc(secret, int64(interval), 6)

	r.POST("/otp", service.TotpHandler(&mockRepo{}, om))

	otp := om.GenerateOtp([]byte("aloha"))
	user := service.OtpPayload{
		Username: "bao",
		Otp:      otp,
	}

	payload, err := json.Marshal(user)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	require.Nil(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	time.Sleep(time.Duration(interval) * time.Second)

	w1 := httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	require.Nil(t, err)
	r.ServeHTTP(w1, req)
	require.Equal(t, http.StatusBadRequest, w1.Result().StatusCode)
}
