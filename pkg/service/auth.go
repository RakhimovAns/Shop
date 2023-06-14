package service

import (
	"fmt"
	"github.com/RakhimovAns/Shop/types"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func Auth(channel chan *int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			auth := request.Header.Get("Authorization")
			if auth == "" {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			bearerToken := strings.Split(auth, " ")
			if bearerToken[0] != "Bearer" {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if len(bearerToken) != 2 {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token, id, err := ParseToken(bearerToken[1])
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if token.Valid {
				channel <- &id
				next.ServeHTTP(writer, request)
				return
			}
		})
	}
}

func ParseToken(accessToken string) (*jwt.Token, int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &types.TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing methon: %v ", token.Header["alg"])
		}
		return []byte("My Key"), nil
	})
	if err != nil {
		return nil, -1, err
	}
	claims, ok := token.Claims.(*types.TokenClaim)
	if !ok {
		return nil, -1, err
	}
	return token, claims.UserID, err
}
