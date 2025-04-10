package service

import (
	"context"
	"errors"
	"example/totp/cache"
	jwttoken "example/totp/jwt_token"
	"example/totp/otp"
	"example/totp/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpPayload struct {
	Otp uint64 `json:"otp"`
}

func getClientSecret(ctx context.Context, repo repository.RepositoryIf, cacheManager cache.CacheIf, username *string) (string, error) {
	s, err := cacheManager.Get(*username)
	if err == nil {
		userSession, ok := s.(*Session)
		if !ok {
			log.Fatal("getClientSecret(): should cast to session normally")
		}
		return userSession.ClientSecret, nil
	}
	user, err := repo.GetUser(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}
	cacheManager.Set(user.Username, &Session{ClientSecret: user.Secret})
	return user.Secret, nil
}

func TotpHandler(repo repository.RepositoryIf, otpManager otp.OtpManagerIf, cacheManager cache.CacheIf) gin.HandlerFunc {
	return HandleWithClaims(func(ctx *gin.Context, claims *jwttoken.Claims) {
		username, ok := (*claims)["username"].(string)
		if !ok {
			log.Println("TotpHandler(): could not get username")
			ctx.Status(http.StatusInternalServerError)
			return
		}
		var payload OtpPayload
		err := ctx.BindJSON(&payload)
		if err != nil {
			log.Println("TotpHandler(): could not bind json : ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		secret, err := getClientSecret(ctx, repo, cacheManager, &username)
		if err != nil {
			log.Println("TotpHandler(): could not get client secret : ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		otp := otpManager.GenerateOtp([]byte(secret))
		if otp != payload.Otp {
			log.Println("TotpHandler(): otp mismatch, expect: ", payload.Otp, ", generate: ", otp)
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.Status(http.StatusOK)
	})
}

func SeedHandler(repo repository.RepositoryIf, cacheManager cache.CacheIf) gin.HandlerFunc {
	return HandleWithClaims(func(ctx *gin.Context, claims *jwttoken.Claims) {
		username, ok := (*claims)["username"].(string)
		if !ok {
			log.Println("SeedHandler(): cannot get username")
			ctx.Status(http.StatusInternalServerError)
			return
		}
		secret, err := getClientSecret(ctx, repo, cacheManager, &username)
		if err != nil {
			log.Println("TotpHandler(): could not get client secret : ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"seed": secret,
		})
	})
}
