package mobile

import (
	"context"
	"ekyc-app/package/token"
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"ekyc-app/source/model"
	"ekyc-app/source/ws"
	"errors"
)

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
