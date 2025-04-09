package router

import (
	"example/totp/appstate"
	jwttoken "example/totp/jwt_token"
	"example/totp/otp"
	"example/totp/repository"
	"example/totp/service"
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
	router.POST("/login", service.LoginHandler(repo, tokenManager))
	router.POST("/otp_validate", service.AuthorizationMiddleware(tokenManager), service.TotpHandler(repo, om))

	return router
}
