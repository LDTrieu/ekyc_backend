package auth

import (
	"context"
	"ekyc-app/package/token"
	"ekyc-app/package/wlog"
	"time"
)

const (
	dur_time = 24 * time.Hour
)

func GenerateJWTLoginSession(ctx context.Context, session_id, account_id string) (
	id string, token_accesss *token.JWTDetails, err error) {
	jwtKey, err := loadPrePrivKey(ctx)
	if err != nil {
		return "", nil, err
	}
	return token.GenerateJWT(jwtKey, 24*time.Hour,
		&DataJWT{
			SessionID: session_id,
			AccountID: account_id,
		})
}

func ValidateLoginJWT(ctx context.Context, jwt_token string) (
	*DataJWT, token.Status, error) {
	jwtKey, err := loadPrePrivKey(ctx)
	if err != nil {
		return nil, token.EXCEPTION, err
	}
	var temp DataJWT
	status, err := token.ValidateJWT(jwtKey,
		jwt_token, &temp)
	if err != nil {
		return nil, token.EXCEPTION, err
	}
	wlog.Info(ctx, temp)
	return &temp, status, nil

}
