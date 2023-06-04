package help

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

type TokenClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"userid"`
}

func ParseToken(accessToken string) (*jwt.Token, int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing methon: %v ", token.Header["alg"])
		}
		return []byte("My Key"), nil
	})
	if err != nil {
		return nil, -1, err
	}
	claims, ok := token.Claims.(*TokenClaim)
	if !ok {
		return nil, -1, err
	}
	return token, claims.UserID, err
}
