package token

import "github.com/golang-jwt/jwt"

const (
	SUCCEED Status = iota
	INPUT_EMPTY
	ACCESS_TOKEN_INVALID
	ACCESS_TOKEN_EXPIRED
	//REFRESH_TOKEN_INVALID
	//REFRESH_TOKEN_EXPIRED

	EXCEPTION
)

type Status uint32

func (ins Status) String() string {
	switch ins {
	case SUCCEED:
		return "SUCCEED"
	case INPUT_EMPTY:
		return "INPUT_EMPTY"
	case ACCESS_TOKEN_INVALID:
		return "ACCESS_TOKEN_INVALID"
	case ACCESS_TOKEN_EXPIRED:
		return "ACCESS_TOKEN_EXPIRED"
	// case REFRESH_TOKEN_INVALID:
	// 	return "REFRESH_TOKEN_INVALID"
	// case REFRESH_TOKEN_EXPIRED:
	// 	return "REFRESH_TOKEN_EXPIRED"
	default:
		return "EXCEPTION"
	}

}

type Claims struct {
	GoEncoded []byte `json:"goEncoded"`
	jwt.StandardClaims
}
type JWTDetails struct {
	Token               string
	IssuedAt, ExpiresAt int64
}
