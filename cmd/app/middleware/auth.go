package middleware

import (
	"github.com/RakhimovAns/Shop/cmd/help"
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
			token, id, err := help.ParseToken(bearerToken[1])
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
