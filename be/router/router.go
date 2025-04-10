package router

import (
	"example/totp/appstate"
	"example/totp/cache"
	jwttoken "example/totp/jwt_token"
	"example/totp/otp"
	"example/totp/repository"
	"example/totp/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	state := appstate.GetAppState()
	router := gin.Default()
	repo := repository.New(state.Db)
	otpConfig := &state.Config.Otp
	om := otp.New([]byte(otpConfig.Secret), int64(otpConfig.Interval), otpConfig.Digit)
	tokenManager := jwttoken.New([]byte(otpConfig.Secret), 1*time.Hour)
	cacheManager := cache.NewInMemCache(24 * time.Hour)
	router.POST("/login", service.LoginHandler(repo, tokenManager))
	router.POST("/otp_validate", service.AuthorizationMiddleware(tokenManager), service.TotpHandler(repo, om, cacheManager))
	router.GET("/seed", service.AuthorizationMiddleware(tokenManager), service.SeedHandler(repo, cacheManager))

	go func() {
		for {
			token := om.GenerateOtp([]byte("aloha"))
			log.Println("bao token: ", token)
			time.Sleep(1 * time.Second)
		}

	}()
	return router
}
