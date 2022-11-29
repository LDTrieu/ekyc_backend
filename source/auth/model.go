package auth

import (
	"context"
	"ekyc-app/package/token"
)

func GenerateJWTLoginSession(ctx context.Context, session_id, account_id string) (
	id string, token_accesss *token.JWTDetails, err error) {
	if err != nil {
		return "", nil, err
	}
	return
}
