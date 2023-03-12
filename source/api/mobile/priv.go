package mobile

import (
	"context"
	"ekyc-app/library/ascii"
	"ekyc-app/package/token"
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"ekyc-app/source/model"
	"ekyc-app/source/ws"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func validateBearer(ctx context.Context,
	r *http.Request) (
	int, string, *auth.DataJWT, error) {
	var (
		excute = func(ctx context.Context, r *http.Request) (int, string, *auth.DataJWT, error) {
			var (
				// parseBearerAuth parses an HTTP Bearer Authentication string.
				// "Bearer QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns QWxhZGRpbjpvcGVuIHNlc2FtZQ.
				parseBearerAuth = func(auth string) (token string, ok bool) {
					const prefix = "Bearer "
					// Case insensitive prefix match. See Issue 22736.
					if len(auth) < len(prefix) || !ascii.EqualFold(auth[:len(prefix)], prefix) {
						return "", false
					}
					return auth[len(prefix):], true
				}
			)
			headerAuth := r.Header.Get("Authorization")
			if len(headerAuth) <= 0 {
				return http.StatusBadRequest, "", &auth.DataJWT{}, errors.New("authorization is empty")
			}
			bearer_token, ok := parseBearerAuth(headerAuth)
			if !ok {
				return http.StatusBadRequest, "", &auth.DataJWT{}, errors.New("authorization is invalid")
			}

			// get from cache DB
			_, account_id, ok, err := fsdb.PersonProfile.GetAccountIdByToken(ctx, bearer_token)
			if err != nil {
				return http.StatusForbidden, bearer_token, &auth.DataJWT{}, err
			}
			if !ok {
				return http.StatusForbidden, bearer_token, &auth.DataJWT{}, errors.New("token no login")
			}
			jwt_data, status, err := auth.ValidateLoginJWT(ctx, bearer_token)
			if err != nil {
				println("ValidateLoginJWT:", err.Error())
			}

			switch status {
			case token.INPUT_EMPTY:
				return http.StatusForbidden, bearer_token, jwt_data, errors.New("token is empty")
			case token.ACCESS_TOKEN_INVALID:
				return http.StatusForbidden, bearer_token, jwt_data, errors.New("token is invalid")
			case token.ACCESS_TOKEN_EXPIRED:
				return http.StatusForbidden, bearer_token, jwt_data, errors.New("token is expired")
			case token.SUCCEED:
				if jwt_data.AccountID != account_id {
					return http.StatusForbidden, bearer_token, jwt_data, errors.New("account_id has been changed")
				}
				// auth pass
				return http.StatusOK, bearer_token, jwt_data, nil
			default:
				return http.StatusForbidden, bearer_token, jwt_data, errors.New("validate token exception")
			}
		}
	)
	status, token, data, err := excute(ctx, r)
	if err != nil {
		println("[AUTH] ", r.RequestURI, "| Error:", err.Error())
	}
	println("[AUTH] ", r.RequestURI, "| Status:", status)
	println("[AUTH] ", r.RequestURI, "| Token:", token)
	println("[AUTH] ", r.RequestURI, "| Access Rights:", fmt.Sprintf("%+v", data))
	return status, token, data, err
}

func guestSubmitQRLogin(ctx context.Context, request *submitQRLoginRequest) (submitQRLoginResponse, error) {
	if err := request.Payload.validate(); err != nil {
		return submitQRLoginResponse{Code: model.StatusBadRequest, Message: err.Error()}, err
	}
	// validate JWT

	req_data, status, err := auth.ValidateLoginJWT(ctx, request.Payload.QrData)
	if status != token.SUCCEED {
		return submitQRLoginResponse{Code: model.StatusForbidden, Message: status.String()}, err
	}
	if err != nil {
		return submitQRLoginResponse{Code: model.StatusForbidden, Message: err.Error()}, err
	}
	// validate session by JWT (on DB) and get sessionID
	login_id, session_id, ok, err := fsdb.LoginSessionDBC.GetSessionId(ctx, request.Payload.QrData)
	if err != nil {
		return submitQRLoginResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !ok {
		return submitQRLoginResponse{Code: model.StatusForbidden, Message: "NOT_FOUND"}, errors.New("token does not exits")
	}
	if session_id != req_data.SessionID {
		return submitQRLoginResponse{Code: model.StatusForbidden, Message: "DATA_CHANGED"}, errors.New("session_id not matched")
	}
	// gen new JWT and replace old token
	_, jwt, err := auth.GenerateJWTLoginSession(ctx, session_id, request.Payload.AccountID)
	if err != nil {
		return submitQRLoginResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	// push notify to
	var (
		push = false
		noti = ws.WsRequestModel{
			Command:  ws.CmdLoginRedirect,
			JWT:      jwt.Token,
			Redirect: ws.RedirectStoreRegister,
		}
	)
	defer func() {
		if !push {
			return
		}
		ws.Station.PushSender(session_id, &noti)
		//
		fsdb.LoginSessionDBC.Revoke(ctx, login_id)
	}()

	// // get & validate user profile
	// resp_code, user, err := facepayapp.RequestUserProfile(ctx, request.Payload.AccountID)
	// if err != nil {
	// 	return submitQRLoginResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	// }
	// if resp_code != 0 {
	// 	return submitQRLoginResponse{Code: model.StatusBadRequest, Message: "accountId is removed"}, errors.New("account_id does not exist")
	// }
	// if !user.IsEnableFacepayWallet {
	// 	push = true
	// 	noti.Redirect = ws.RedirectMissWallet
	// 	return submitQRLoginResponse{Payload: submit_qr_login_resp{AccountID: request.Payload.AccountID}}, nil
	// }
	// case user is blocked ?

	// make Person_Profile (update token if exist else create new profile)
	_, exist, err := fsdb.PersonProfile.CreateIfNotExist(ctx, request.Payload.AccountID, session_id, jwt.Token)
	if err != nil {
		return submitQRLoginResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !exist {
		noti.Redirect = ws.RedirectStoreRegister
	} else {
		// switch person.CensorshipStatus {
		// case fsdb.StatusAccepted:
		// 	noti.Redirect = ws.RedirectLogin
		// case fsdb.StatusWaiting:
		// 	noti.Redirect = ws.RedirectStoreRegister
		// default:
		noti.Redirect = ws.RedirectStoreSubmited
		//}
	}
	push = true
	noti.PersonProfile.AccountId = request.Payload.AccountID
	//noti.PersonProfile.Fullname = person.FullName
	//noti.PersonProfile.IsBlocked = person.IsBlock
	return submitQRLoginResponse{Payload: submit_qr_login_resp{AccountID: request.Payload.AccountID}}, nil
}

// /* */
func __loginTerminal(ctx context.Context,
	request *loginTerminalRequest) (
	loginTerminalResponse, error) {
	var (
		login_session_id = uuid.NewString()
	)
	// if err := request.validate(); err != nil {
	// 	return loginTerminalResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }

	// Select from DB
	// check email exist, and check hashed password
	//get info
	docId, isBlocked, err := fsdb.DeviceProfile.CheckLogin(ctx, request.TerminalId, request.Password)
	if err != nil {
		return loginTerminalResponse{
			Code:    model.StatusForbidden,
			Message: err.Error()}, err
	}

	if isBlocked {
		return loginTerminalResponse{
			Code:    model.StatusMethodNotAllowed,
			Message: "terminal is blocked"}, errors.New("terminal is blocked")
	}
	// gen and save token
	_, jwt_login, err := auth.GenerateJWTLoginSession(
		ctx, login_session_id, request.TerminalId)
	if err != nil {
		return loginTerminalResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}
	// save to cache
	// _, err = fsdb.LoginSessionDBC.Add(ctx, login_session_id, jwt_login.Token, jwt_login.ExpiresAt)
	// if err != nil {
	// 	return loginTerminalResponse{
	// 		Code:    model.StatusServiceUnavailable,
	// 		Message: err.Error()}, err
	// }
	// save to profile
	var (
		linked_at = time.Now()
	)
	err = fsdb.DeviceProfile.SetToken(ctx, docId, jwt_login.Token, linked_at)
	if err != nil {
		return loginTerminalResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}
	payload := login_terminal_data{
		TerminalId: request.TerminalId,
		LinkedAt:   linked_at,
		Token:      jwt_login.Token,
		Avt:        "https://www.telpo.com.cn/uploadfiles/Products/C9/Telpo-C9-02.png",
	}

	return loginTerminalResponse{
		Payload: payload,
	}, nil
}

/* */
func __sigupTerminal(ctx context.Context,
	request *signupTerminalRequest) (
	signupTerminalResponse, error) {
	if err := request.validate(); err != nil {
		return signupTerminalResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	// check terminal_id exist
	email_already_exist, err := fsdb.DeviceProfile.ValidateTerminalId(ctx, request.TerminalId)
	if err != nil {
		return signupTerminalResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if email_already_exist {
		return signupTerminalResponse{Code: model.StatusEmailDuplicated, Message: "DATA_ALREADY_EXIST"},
			errors.New("terminal_id is duplicated")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password), 8)

	if err != nil {
		return signupTerminalResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	// Insert to DB
	_, err = fsdb.DeviceProfile.Add(ctx, request.TerminalId, request.Avt, string(hashedPassword), request.Permit.AccountID)
	if err != nil {
		return signupTerminalResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}

	return signupTerminalResponse{}, nil
}
