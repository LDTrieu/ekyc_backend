package token

import (
	"bytes"
	"time"
)

var (
	DefaultAccessSecret  = bytes.NewBufferString("f268192bee1e39664050f34b90f8b087ce5f0962857221d3f9616d0e592294bc")
	DefaultRefreshSecret = bytes.NewBufferString("c94e8f770a87d394476d3271684fe0115785945143d8052f0b27c8d9f428c244")
)

func GenerateJWT(secret []byte, longevity time.Duration, v interface{}) (
	id string, _ *JWTDetails, err error) {
	return
}
