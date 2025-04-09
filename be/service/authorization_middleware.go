package service

import (
	jwttoken "example/totp/jwt_token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	tokenClaimsKey = "tokenClaims"
)

func AuthorizationMiddleware(tokenManager jwttoken.JwtTokenManagerIf) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auths, ok := ctx.Request.Header["Authorization"]
		if !ok || len(auths) == 0 {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		auth := strings.Split(auths[0], " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		token := auth[1]
		claims, err := tokenManager.VerifyToken(&token)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		ctx.Set(tokenClaimsKey, claims)
		ctx.Next()
	}
}

func HandleWithClaims(handle func(*gin.Context, *jwttoken.Claims)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, ok := ctx.Get(tokenClaimsKey)
		if !ok {
			log.Println("HandleWithClaims(): Failed to get tokenClaims")
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		claims, ok := c.(*jwttoken.Claims)
		if !ok {
			log.Println("HandleWithClaims(): Failed to cast claims")
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		handle(ctx, claims)
	}
}
