package service

import (
	jwttoken "example/totp/jwt_token"
	"example/totp/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(repo repository.RepositoryIf, tokenManager jwttoken.JwtTokenManagerIf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload LoginPayload
		err := ctx.BindJSON(&payload)
		if err != nil {
			log.Println("LoginHandler(): could not bind json: ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		log.Println("LoginHandler(): got login request: ", payload)
		user, err := repo.GetUser(ctx, &payload.Username)
		if err != nil {
			log.Println("LoginHandler(): could not find user in db: ", err)
			ctx.Status(http.StatusBadRequest)
			return
		}
		if user.Username != payload.Username {
			log.Println("LoginHandler(): get wrong user???, expect: ", payload.Username, " , got: ", user.Username)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		if user.Password != payload.Password {
			log.Println("LoginHandler(): wrong password")
			ctx.Status(http.StatusBadRequest)
			return
		}
		claims := &jwttoken.Claims{
			"username": user.Username,
		}
		token, err := tokenManager.GenerateToken(claims)
		if err != nil {
			log.Println("LoginHandler(): cannot generate token ???: ", err)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"access_token": token,
		})
	}
}
