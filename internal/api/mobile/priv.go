package mobile

import (
	"context"
	"ekyc-app/internal/auth"
	"ekyc-app/internal/fsdb"
	"ekyc-app/internal/model"
	"ekyc-app/internal/service/faceauth"
	"ekyc-app/internal/ws"
	"ekyc-app/library/ascii"
	"ekyc-app/package/token"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			_, terminal_id, ok, err := fsdb.DeviceProfile.GetTerminalIdByToken(ctx, bearer_token)
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
				if jwt_data.AccountID != terminal_id {
					return http.StatusForbidden, bearer_token, jwt_data, errors.New("terminal_id has been changed")
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
	_, err = fsdb.DeviceProfile.Add(ctx, request.TerminalId, request.TerminalName, " request.Avt", string(hashedPassword), request.Permit.AccountID)
	if err != nil {
		return signupTerminalResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}

	return signupTerminalResponse{}, nil
}

/* */
// pingThirdPartyRequest
func __pingThirdParty(ctx context.Context,
	request *pingThirdPartyRequest) (
	pingThirdPartyResponse, error) {
	var (
		mock_data = faceauth.MockModel{
			Name: "name_service_1",
			Code: 123,
		}
	)
	go func() {
		// Call to auth service
		code_resp, data_resp, err := faceauth.RequestSession(ctx,
			&mock_data)
		if err != nil {
			log.Println("err: ", err)
			// return pingThirdPartyResponse{
			// 	Code:    model.StatusServiceUnavailable,
			// 	Message: err.Error()}, err
		}
		log.Println("code_resp: ", code_resp)
		log.Println("data_resp: ", data_resp)
		//log.Println("time.Now() ", time.Now())
	}()
	//log.Println("time.Now() ", time.Now())
	return pingThirdPartyResponse{}, nil
}

/* */
func __faceAuthSession(ctx context.Context, request *faceAuthSessionRequest) (faceAuthSessionResponse, error) {

	// var (
	// 	face_data = options.FormFile{
	// 		Filename: request.Payload.FileName,
	// 		File:     request.Payload.File,
	// 	}
	// )
	//log.Println("time.Now() 1: ", time.Now())
	// code_resp, data_resp, err := faceauth.FaceAuthSession(ctx, &face_data)
	// if err != nil {
	// 	log.Println("ERR: ", err)
	// 	return faceAuthSessionResponse{
	// 		Code:    model.StatusServiceUnavailable,
	// 		Message: err.Error()}, err
	// }
	// log.Println("code_resp: ", code_resp)
	// log.Println("data_resp: ", data_resp)
	// log.Println("time.Now() 2: ", time.Now())
	var (
		student_id  = "n18dccn229"
		session_id  = primitive.NewObjectID().Hex()
		terminal_id = request.Payload.TerminalId
		auth_at     = time.Now()
	)

	// get name, face_id, student_id
	_, full_name, _, _, _, _, _, _, unit_id, _, _, _, _, _, _, exist, err := fsdb.StudentProfile.GetByStudentId(ctx, student_id)
	if err != nil {
		return faceAuthSessionResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}
	if !exist {
		return faceAuthSessionResponse{Code: model.StatusNotFound, Message: "NOT FOUND"}, errors.New("student_id does not exist")
	}

	// log session to db
	_, err = fsdb.AuthSession.Add(ctx, session_id, student_id, "face_id", terminal_id, full_name, unit_id, "image_url",
		auth_at)
	if err != nil {
		return faceAuthSessionResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	// return to terminal
	return faceAuthSessionResponse{
		Payload: face_image_resp{
			FullName:  full_name,
			FaceId:    "data_resp.FaceId",
			StudentId: student_id,
			Avt:       "https://fff5fe774f075efc5ff8a6afe4708446203de4fb1e5932769fb01e1-apidata.googleusercontent.com/download/storage/v1/b/ekyc_image_bucket/o/n18dccn229%2Fimage_file.PNG?jk=Ac_6HjJEn8lHrno2jAS0_lp-N_p_0pFGuSBx069Voo19aEJ7KXnsotoBFTI0mVlYvtUZvs_NhFvJVVoDX2da1M0yqoXfav1Fp6eZimhIjDGrcfI2-r6FzB4PCv4CIxAaFeTfvZXwQXSj-RZacP8BDoc_rz6sF87pcUEG8_KZlDYY3eOZnDdMN7sWGVRQJcBKNHayYWKtuC2dWkGjrXwtEZIubKuuVqCeb1e-GFOC1-LRZq0HaHb2fdh8dfqQsnjsXEPVfKdYCqUAzdb1EWZvPpvL4k8-Yy6xvPEtBoEnwf-vBHD2fREJ-1x5TYcuqkPYUgnzm8WGTfPXMFaIQBQ8mquuWLTC8LuY6yd0TVz1MU63iY0RhkHBj8SEm08YPWIwERIYQ_RY-BdN2eGCGhvis52qDVvYlIUB3XougeCCgbmgJMyx8sUErzjG4hNyXu8Hl7ZcFXdfJkij3mtH9BjQH7oOihA5nM2Y433WV2b3LHnIHHsBIf4eCK8WHiE6Ch3YzOgr9win2ChKi2nhQMm6--nCht0_c5W40javgyfU5Z1it_iKuWt3uh0f182h1QZvJ3MF2hFZIvprfk62p3izw1QDo1eISHE01SVTd9k16d61R2wDLJH6G5FJS56X7gWyTQu7yWNEmCo0UEF05p_zrNNrMkwgh7tmH_2eFD0_IxfBQBPw1PXDoCeSje8x4k1AK1S2k5evknDNrxOuYfP-CysNpLCIMV_NEoFFHA96fLVQyd5OBcd13rMnamfsrOM4Rxe6ERfKFQHTeHLLqz4wQtyMe_3AKJiubCSSNwk_1un5YW_BXNrX4TBB7fv-gEqATfQYZVKv2yEMMYdS034IwjqA14dRcm0e2XO_jE0bCiUu_rWzDr59P8rEU-cNlt32CHCLwCgMPOvZb19IJUcdVBTsxA5u7NmuZIIV3fOcRWsWJ67w5zrA8WY_Pd4fuKTDFxrn7qHS1UiJPQXPUAZCFRbIxB84yHZ3kClrF-vmiAknRSsmMcy5-RiKV7XE34dcyRajO3KjFmKeJaB7yTScg7c_AIBRKrGKembQYghqkI_F64VFV5UL7qg22ITaHtWiJtWhDaz-sGBlrG2Y3yTEs8cqK4pmQ0hL8Mc9_zjO6gU2gUwn4tc9LRCLWHDGtlyl3Tbj2LCl22waGowIJ0rD2Q31tmGMaSp_R3FSshe4aY-hD5Sml9OGIGbGD10tHD5wGxJvXZ0T0lFaSa-A4wBlALpiZZ3XJYPmZ8KMdMga_Y-wzZqQv_tfQx2JTMGe&isca=1",
			UnitId:    "data_resp.UnitId",
			AuthAt:    auth_at,
		},
	}, nil
}
