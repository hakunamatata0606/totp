package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"example/totp/cache"
	jwttoken "example/totp/jwt_token"
	"example/totp/otp"
	"example/totp/repository"
	"example/totp/service"
	"fmt"
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
			Password: "123",
			Secret:   "aloha",
		}, nil
	}
	return nil, errors.New("not found")
}

func TestTotpHandler(t *testing.T) {
	secret := []byte("this is secret")
	interval := 5

	r := gin.Default()

	om := otp.New(secret, int64(interval), 6)
	tm := jwttoken.New([]byte("aloha"), 10*time.Second)
	cm := cache.NewInMemCache(24 * time.Hour)

	claims := &jwttoken.Claims{
		"username": "bao",
	}
	token, err := tm.GenerateToken(claims)
	require.Nil(t, err)

	r.Use(service.AuthorizationMiddleware(tm))
	r.POST("/otp", service.TotpHandler(&mockRepo{}, om, cm))

	otp := om.GenerateOtp([]byte("aloha"))
	user := service.OtpPayload{
		Otp: otp,
	}

	payload, err := json.Marshal(user)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	require.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	time.Sleep(time.Duration(interval) * time.Second)

	w1 := httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/otp", strings.NewReader(string(payload)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	require.Nil(t, err)
	r.ServeHTTP(w1, req)
	require.Equal(t, http.StatusBadRequest, w1.Result().StatusCode)
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func TestLoginService(t *testing.T) {
	tm := jwttoken.New([]byte("aloha"), 1*time.Second)
	r := gin.Default()
	r.POST("/login", service.LoginHandler(&mockRepo{}, tm))
	user := &service.LoginPayload{
		Username: "bao",
		Password: "123",
	}
	payload, err := json.Marshal(user)
	require.Nil(t, err)
	req, err := http.NewRequest("POST", "/login", strings.NewReader(string(payload)))
	require.Nil(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	var respJson LoginResponse
	err = json.NewDecoder(w.Body).Decode(&respJson)
	require.Nil(t, err)
	claims, err := tm.VerifyToken(&respJson.AccessToken)
	require.Nil(t, err)
	require.Equal(t, "bao", (*claims)["username"])
	time.Sleep(2 * time.Second)
	_, err = tm.VerifyToken(&respJson.AccessToken)
	require.NotNil(t, err)
}

type SeedResponse struct {
	Seed string `json:"seed"`
}

func TestSeedHandler(t *testing.T) {
	r := gin.Default()

	username := "bao"
	tm := jwttoken.New([]byte("aloha"), 10*time.Second)
	cm := cache.NewInMemCache(24 * time.Hour)
	repo := &mockRepo{}
	claims := &jwttoken.Claims{
		"username": username,
	}
	token, err := tm.GenerateToken(claims)
	require.Nil(t, err)

	r.Use(service.AuthorizationMiddleware(tm))
	r.GET("/seed", service.SeedHandler(repo, cm))

	req, err := http.NewRequest("GET", "/seed", nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	user, err := repo.GetUser(context.Background(), &username)
	require.Nil(t, err)
	require.Equal(t, username, user.Username)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	var respJson SeedResponse
	err = json.NewDecoder(w.Body).Decode(&respJson)
	require.Nil(t, err)
	require.Equal(t, user.Secret, respJson.Seed)
}
