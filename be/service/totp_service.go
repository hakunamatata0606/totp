package service

import (
	"example/totp/otp"
	"example/totp/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpPayload struct {
	Username string `json:"username"`
	Otp      uint64 `json:"otp"`
}

func TotpHandler(repo repository.RepositoryIf, otpManager otp.OtpManagerIf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload OtpPayload
		err := ctx.BindJSON(&payload)
		if err != nil {
			log.Println("TotpHandler(): could not bind json : ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		log.Println("TotpHandler(): got otp request: ", payload)
		user, err := repo.GetUser(ctx, &payload.Username)
		if err != nil {
			log.Println("TotpHandler(): could not find user in db: ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		if user.Username != payload.Username {
			log.Println("TotpHandler(): get wrong user???, expect: ", payload.Username, " , got: ", user.Username)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		otp := otpManager.GenerateOtp(user.Secret)
		if otp != payload.Otp {
			log.Println("TotpHandler(): otp mismatch, expect: ", payload.Otp, ", generate: ", otp)
			ctx.Status(http.StatusBadRequest)
			return
		}
		ctx.Status(http.StatusOK)
	}
}
