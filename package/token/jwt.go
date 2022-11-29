package token

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	DefaultAccessSecret  = bytes.NewBufferString("f268192bee1e39664050f34b90f8b087ce5f0962857221d3f9616d0e592294bc")
	DefaultRefreshSecret = bytes.NewBufferString("c94e8f770a87d394476d3271684fe0115785945143d8052f0b27c8d9f428c244")
)

func GenerateJWT(secret []byte, longevity time.Duration, v interface{}) (
	id string, _ *JWTDetails, err error) {
	var (
		new_claim = func(encoded_data []byte, le time.Duration) *Claims {
			var random string
			uuid, err := uuid.NewRandom()
			if err != nil {
				random = time.Now().GoString()
			} else {
				random = uuid.String()
			}
			return &Claims{
				GoEncoded: encoded_data,
				StandardClaims: jwt.StandardClaims{
					Id:        random + "-" + primitive.NewObjectID().Hex(),
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(le).Unix(),
				},
			}

		}
	)
	// Stand-in for a network connection
	var network bytes.Buffer
	if v != nil {
		if err := gob.NewEncoder(&network).Encode(v); err != nil {
			println(">> GenerateJWT: Encode GoData error:", err.Error())
		}
	}
	// Creating Cliams
	claim := new_claim(network.Bytes(), longevity)
	// Signature claims
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claim).SignedString(secret)
	if err != nil {
		return "", nil, err
	}
	return claim.Id, &JWTDetails{token, claim.IssuedAt, claim.ExpiresAt}, nil
}

func ValidateJWT(secret []byte, token string, v interface{}) (
	Status, error) {
	if len(token) == 0 {
		return INPUT_EMPTY, errors.New("token must be not empty")
	}
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(
		token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing methood: %v",
				token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if parsedToken == nil && err != nil {
		return ACCESS_TOKEN_INVALID, errors.New("access_token invalid")
	}
	claim, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return ACCESS_TOKEN_INVALID, errors.New("invalid access_token")
	}
	if claim.ExpiresAt < time.Now().Unix() {
		return ACCESS_TOKEN_EXPIRED, errors.New("access_token expired")
	}
	if len(claim.GoEncoded) > 0 {
		// Stand-in for a network connection
		var network bytes.Buffer
		if _, err := network.Write(claim.GoEncoded); err != nil {
			println(">> Write buffer error:", err.Error())
		}
		if err := gob.NewDecoder(&network).Decode(v); err != nil {
			println(">> ValidateJWT: Decode GoData error:", err.Error())
		}
	}
	return SUCCEED, nil
}
