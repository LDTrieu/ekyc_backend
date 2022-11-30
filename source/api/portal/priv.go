package portal

import (
	"context"
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"ekyc-app/source/model"
	"fmt"

	"github.com/google/uuid"
)

func guestRendQRLogin(ctx context.Context,
	request *rendQRLoginRequest) (
	rendQRLoginResponse, error) {
	var (
		login_session_id = uuid.NewString()
	)
	// make new JWT authen with account_id equal to empty
	_, jwt_login, err := auth.GenerateJWTLoginSession(
		ctx, login_session_id, "")
	if err != nil {
		return rendQRLoginResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}
	//save to cache
	doc_id, err := fsdb.LoginSessionDBC.
		Add(ctx, login_session_id,
			jwt_login.Token, jwt_login.ExpiresAt)
	if err != nil {
		return rendQRLoginResponse{Code: model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}
	// end
	return rendQRLoginResponse{Payload: rend_qr_code_login_resp{
		Path:      fmt.Sprintf("/login/qr/download/%s/%s?action=%s", doc_id, uuid.NewString(), qrActionLoginWebPortal),
		JWT:       jwt_login.Token,
		IssuedAt:  jwt_login.IssuedAt,
		ExpiresIn: jwt_login.ExpiresAt - jwt_login.IssuedAt,
	}}, nil

}
