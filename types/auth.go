package types

import "github.com/golang-jwt/jwt"

type TokenClaim struct {
	jwt.StandardClaims
	UserID int64 `json:"userid"`
}
