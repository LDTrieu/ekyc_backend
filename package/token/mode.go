package token

import "github.com/golang-jwt/jwt"

type Cliams struct {
	//GoEncoded []byte `json:"goEncoded"`
	jwt.StandardClaims
}
type JWTDetails struct {
	Token               string
	IssuedAt, ExpiresAt int64
}
