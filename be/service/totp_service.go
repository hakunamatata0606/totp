package service

import (
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

func TotpHandler(repo repository.RepositoryIf, otpManager otp.OtpManagerIf) gin.HandlerFunc {
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
		user, err := repo.GetUser(ctx, &username)
		if err != nil {
			log.Println("LoginHandler(): could not find user in db: ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		if user.Username != username {
			log.Println("LoginHandler(): get wrong user???, expect: ", username, " , got: ", user.Username)
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
	})
}
