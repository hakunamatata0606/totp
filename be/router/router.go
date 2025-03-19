package router

import (
	"example/totp/appstate"
	"example/totp/otp"
	"example/totp/repository"
	"example/totp/service"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	state := appstate.GetAppState()
	router := gin.Default()
	repo := repository.New(state.Db)
	otpConfig := &state.Config.Otp
	om := otp.New([]byte(otpConfig.Secret), int64(otpConfig.Interval), otpConfig.Digit)

	router.POST("/otp_validate", service.TotpHandler(repo, om))

	return router
}
