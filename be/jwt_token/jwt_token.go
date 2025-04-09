package jwttoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims map[string]interface{}

type JwtTokenManagerIf interface {
	GenerateToken(claims *Claims) (string, error)
	VerifyToken(token *string) (*Claims, error)
}

type jwtTokenManagerImpl struct {
	secret []byte
	expire time.Duration
}

func New(secret []byte, expire time.Duration) JwtTokenManagerIf {
	return &jwtTokenManagerImpl{
		secret: secret,
		expire: expire,
	}
}

func (tm *jwtTokenManagerImpl) GenerateToken(claims *Claims) (string, error) {
	now := time.Now()
	jwtClaims := (*jwt.MapClaims)(claims)
	(*jwtClaims)["iat"] = now.Unix()
	(*jwtClaims)["exp"] = now.Add(tm.expire).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := token.SignedString(tm.secret)
	return signedToken, err
}

func (tm *jwtTokenManagerImpl) VerifyToken(token *string) (*Claims, error) {
	jwtToken, err := jwt.Parse(*token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return tm.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return (*Claims)(&claims), nil
}
