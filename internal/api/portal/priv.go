package portal

import (
	"context"
	"ekyc-app/gcloud"
	"ekyc-app/internal/auth"
	"ekyc-app/internal/fsdb"
	"ekyc-app/internal/model"
	"ekyc-app/internal/ws"
	"ekyc-app/library/ascii"
	"ekyc-app/package/socket"
	"ekyc-app/package/token"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

/* */
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
	log.Println(doc_id)
	// end
	return rendQRLoginResponse{Payload: rend_qr_code_login_resp{
		//Path:      fmt.Sprintf("/login/qr/download/%s/%s?action=%s", doc_id, uuid.NewString(), qrActionLoginWebPortal),
		JWT:       jwt_login.Token,
		IssuedAt:  jwt_login.IssuedAt,
		ExpiresIn: jwt_login.ExpiresAt - jwt_login.IssuedAt,
	}}, nil

}

/* WS */
func websocketConnection(c *gin.Context) {
	var (
		r   = c.Request
		w   = c.Writer
		ctx = c.Request.Context()

		get_token = func(r *http.Request) (string, error) {
			h := strings.TrimSpace(r.Header.Get("Sec-Websocket-Protocol"))
			if h == "" {
				return "", errors.New("subprotocols is invalid")
			}
			protocols := strings.Split(h, ",")
			for i := range protocols {
				protocols[i] = strings.TrimSpace(protocols[i])
			}
			if len(protocols) > 0 {
				return protocols[0], nil
			}
			return "", errors.New("subprotocols is invalid")
		}

		get_session_id = func(token string) string {
			// Try search on LoginSessionDBC
			if _, session_id, ok, err := fsdb.LoginSessionDBC.GetSessionId(ctx, token); err != nil {
				println("Websocket-Connecting:", err.Error(), "\nTry again search on Profile.")
			} else if ok && len(session_id) > 0 {
				return session_id
			}
			// Retry search on ProfileToken
			if _, session_id, ok, err := fsdb.PersonProfile.GetSessionID(ctx, token); err != nil {
				println("Websocket-Connecting:", err.Error())
				return ""
			} else if !ok {
				println("Websocket-Connecting: token does not exist")
				return ""
			} else {
				return session_id
			}
		}

		connection_id = "default"
		read          = make(chan []byte, 100)
		write         = make(chan []byte, 100)

		websocket_receiver = func(ctx context.Context, read <-chan []byte) {
			var failure_count int
			for {
				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != nil {
						println("Close connection by client:", err.Error())
					} else {
						println("Close connection by client.")
					}
					return
				case _, ok := <-read:
					if !ok {
						if failure_count > 3 {
							println("can not read payload data becase READ channel is closed.")
						}
						failure_count++
						time.Sleep(300 * time.Millisecond)
						continue
					}
				}
			}
		}
	)

	token, err := get_token(r)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	println("Websocket-Connecting: Protocol [", token, "]")
	connection_id = get_session_id(token)
	if len(connection_id) == 0 {
		println("Websocket-Connecting: disconnected")
		c.AbortWithError(http.StatusUnauthorized, errors.New("token does not exist"))
		return
	}
	println("Websocket-Connecting: ConnectionID=", connection_id)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		Subprotocols:    []string{token},
		CheckOrigin: func(r *http.Request) bool {
			// allow all connections by default
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	defer conn.Close()

	go websocket_receiver(ctx, read)

	// Add connection to Station
	ws.Station.AddClient(connection_id, write)

	// create read-only-chanel (task) & write-only-chanel (task) and put in the function Listen
	if err := socket.Listen(r, conn, read, write); err != nil {
		println("Error:", err.Error())
	}
}

/* */
func __loginBasic(ctx context.Context,
	request *loginBasicRequest) (
	loginBasicResponse, error) {
	var (
		login_session_id = uuid.NewString()
	)
	if err := request.validate(); err != nil {
		return loginBasicResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}

	// Select from DB
	// check email exist, and check hashed password
	// get info
	docId, accountId, fullName, phoneNumber, birthday, role, unitId, isBlocked, err := fsdb.PersonProfile.CheckLogin(ctx, request.Email, request.Password)
	if err != nil {
		return loginBasicResponse{
			Code:    model.StatusForbidden,
			Message: err.Error()}, err
	}

	if isBlocked {
		return loginBasicResponse{
			Code:    model.StatusMethodNotAllowed,
			Message: "account is blocked"}, errors.New("account is blocked")
	}

	// gen and save token
	_, jwt_login, err := auth.GenerateJWTLoginSession(
		ctx, login_session_id, accountId)
	if err != nil {
		return loginBasicResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}
	// save to cache
	_, err = fsdb.LoginSessionDBC.Add(ctx, login_session_id, jwt_login.Token, jwt_login.ExpiresAt)
	if err != nil {
		return loginBasicResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}
	// save to profile
	err = fsdb.PersonProfile.SetToken(ctx, docId, login_session_id, jwt_login.Token)
	if err != nil {
		return loginBasicResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	payload := login_basic_data{
		AccountId:   accountId,
		FullName:    fullName,
		FirstName:   fullName,
		LastName:    fullName,
		Email:       request.Email,
		PhoneNumber: phoneNumber,
		Birthday:    birthday,
		Role:        role,
		UnitId:      unitId,
		Token:       jwt_login.Token,
		Avt:         "https://png.pngtree.com/png-clipart/20190924/original/pngtree-user-vector-avatar-png-image_4830521.jpg",
		Banner:      "https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/f/cebd17f1-b283-45e5-8600-6ec3edc558fd/dee2aqv-222532a7-8676-4788-b8e3-08d4f5be55e2.png/v1/fill/w_1264,h_632,q_70,strp/profile_banner_by_darkfigure4_dee2aqv-pre.jpg?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1cm46YXBwOjdlMGQxODg5ODIyNjQzNzNhNWYwZDQxNWVhMGQyNmUwIiwiaXNzIjoidXJuOmFwcDo3ZTBkMTg4OTgyMjY0MzczYTVmMGQ0MTVlYTBkMjZlMCIsIm9iaiI6W1t7ImhlaWdodCI6Ijw9NjQwIiwicGF0aCI6IlwvZlwvY2ViZDE3ZjEtYjI4My00NWU1LTg2MDAtNmVjM2VkYzU1OGZkXC9kZWUyYXF2LTIyMjUzMmE3LTg2NzYtNDc4OC1iOGUzLTA4ZDRmNWJlNTVlMi5wbmciLCJ3aWR0aCI6Ijw9MTI4MCJ9XV0sImF1ZCI6WyJ1cm46c2VydmljZTppbWFnZS5vcGVyYXRpb25zIl19.sdy7FtZ92V4tHXX-hTf0PupZmkD7CQoG-BkmOY0_mQg",
	}

	return loginBasicResponse{
		Payload: payload,
	}, nil
}

/* */
func __signupBasic(ctx context.Context,
	request *signupBasicRequest) (
	signupBasicResponse, error) {
	if err := request.validate(); err != nil {
		return signupBasicResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}

	log.Println("request", request)

	// check email exist
	email_already_exist, err := fsdb.PersonProfile.ValidateEmail(ctx, request.Email)
	if err != nil {
		return signupBasicResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if email_already_exist {
		return signupBasicResponse{Code: model.StatusEmailDuplicated, Message: "DATA_ALREADY_EXIST"},
			errors.New("email is duplicated")
	}
	// check phone number exist
	phone_number_already_exist, err := fsdb.PersonProfile.ValidatePhoneNumber(ctx, request.PhoneNumber)
	if err != nil {
		return signupBasicResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if phone_number_already_exist {
		return signupBasicResponse{Code: model.StatusPhoneNumberDuplicated, Message: "DATA_ALREADY_EXIST"},
			errors.New("phone number is duplicated")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password), 8)

	if err != nil {
		return signupBasicResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	// Insert to DB
	accountId := primitive.NewObjectID().Hex()
	sessionId := uuid.NewString()
	info, _, err := fsdb.PersonProfile.CreateSignupProfile(ctx, accountId, sessionId,
		request.Email, request.PhoneNumber, request.FullName, request.Role, request.UnitId, string(hashedPassword), request.Permit.AccountID)
	if err != nil {
		return signupBasicResponse{
			Code:    model.StatusInternalServerError,
			Message: err.Error()}, err
	}

	return signupBasicResponse{
		Payload: signup_basic_data{
			AccountId: info.AccountId,
			Email:     info.Email,
		},
	}, nil
}

/* */
func __filterListUser(
	ctx context.Context,
	request *filterListUserRequest) (
	filterListUserResponse, error) {
	db_users, err := fsdb.PersonProfile.GetAll(ctx)
	if err != nil {
		return filterListUserResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	var (
		list_user = make([]user_data, 0)
	)
	for _, user := range db_users {
		list_user = append(list_user, withUserModel(user))
	}
	return filterListUserResponse{
		Payload: list_user_resp{
			TotalUser: len(list_user),
			ListUser:  list_user,
		},
	}, nil

}

/* */

func __userDetail(ctx context.Context, request *userDetailRequest) (userDetailResponse, error) {
	_, db_user, ok, err := fsdb.PersonProfile.GetByAccountId(ctx, request.AccountId)
	if err != nil {
		return userDetailResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !ok {
		return userDetailResponse{Code: model.StatusNotFound, Message: "NOT FOUND"}, errors.New("student_id does not exist")
	}

	return userDetailResponse{
		Payload: user_detail_data{
			AccountId:   db_user.AccountId,
			FullName:    db_user.FullName,
			Email:       db_user.Email,
			PhoneNumber: db_user.PhoneNumber,
			Role:        db_user.Role,
			UnitId:      db_user.UnitId,
			IsBlocked:   db_user.IsBlocked,
			LastLoginAt: db_user.LastLoginAt,
			CreatedAt:   db_user.CreatedAt,
			CreatedBy:   db_user.CreatedBy,
		}}, nil
}

/* */
func __filterListStudent(
	ctx context.Context,
	request *filterListStudentRequest) (
	filterListStudentResponse, error) {
	db_students, err := fsdb.StudentProfile.GetAll(ctx)
	if err != nil {
		return filterListStudentResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	var (
		list_student = make([]student_data, 0)
	)
	for _, student := range db_students {
		list_student = append(list_student, withStudentModel(student))
	}
	return filterListStudentResponse{
		Payload: list_student_resp{
			TotalStudent: len(list_student),
			ListStudent:  list_student,
		},
	}, nil
}

/* */
func __submitStudentProfile(ctx context.Context, request *createStudentProfileRequest) (createStudentProfileResponse, error) {
	if err := request.Payload.validate(); err != nil {
		return createStudentProfileResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	// validate fields
	// student_id_already_exist, err := fsdb.StudentProfile.ValidateStudentId(ctx, request.Payload.StudentId)
	// if err != nil {
	// 	return createStudentProfileResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	// }
	// if student_id_already_exist {
	// 	return createStudentProfileResponse{Code: model.StatusStudentIdDuplicated, Message: "STUDENT_ID_ALREADY_EXIST"}, errors.New("student_id is duplicated")
	// }

	// national_id_already_exist, err := fsdb.StudentProfile.ValidateNationalId(ctx, request.Payload.NationalId)
	// if err != nil {
	// 	return createStudentProfileResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	// }
	// if national_id_already_exist {
	// 	return createStudentProfileResponse{Code: model.StatusNationalIdDuplicated, Message: "NATIONAL_ID_ALREADY_EXIST"}, errors.New("national_id is duplicated")
	// }

	if err := fsdb.StudentProfile.CreateStudentProfile(ctx,
		request.Payload.StudentId,
		request.Payload.Email,
		request.Payload.FirstName,
		request.Payload.LastName,
		request.Payload.PhoneNumber,
		request.Payload.NationalId,
		request.Payload.Birthday,
		request.Payload.Gender,
		request.Payload.Address,
		request.Payload.Hometown,
		request.Payload.UnitId,
		request.Payload.Image,
		request.Payload.ImageEkyc,
		request.Permit.AccountID); err != nil {
		return createStudentProfileResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	return createStudentProfileResponse{}, nil
}

/* */
func __updateStudentEkyc(ctx context.Context, request *updateStudentEkycRequest) (updateStudentEkycResponse, error) {
	if err := request.Payload.validate(); err != nil {
		return updateStudentEkycResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	// validate fields
	student_id_already_exist, err := fsdb.StudentProfile.ValidateStudentId(ctx, request.Payload.StudentId)
	if err != nil {
		return updateStudentEkycResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !student_id_already_exist {
		return updateStudentEkycResponse{Code: model.StatusDataNotFound, Message: "STUDENT_ID_NOT_EXIST"}, errors.New("student_id not exist")
	}
	// check CCCD - Face Id

	if err := fsdb.StudentProfile.SetEkyc(ctx,
		request.Payload.StudentId,
		request.Payload.NationalId,
		request.Payload.PersonId,
		request.Payload.FullName,
		request.Payload.Gender,
		request.Payload.FaceImageURL,
		request.Payload.NationalIdCardURL,
		request.Payload.Address,
		request.Payload.PlaceOfOrigin,
		request.Payload.Nationality,
		request.Permit.AccountID,
		request.Payload.DateOfBirth,
		request.Payload.DateOfExpiry); err != nil {
		return updateStudentEkycResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	return updateStudentEkycResponse{}, nil
}

/* */
// __updateStudentFaceVideo
func __uploadStudentFaceVideo(ctx context.Context, request *updateStudentFaceVideoRequest) (updateStudentFaceVideoResponse, error) {
	doc_id, _, _, exist, err := fsdb.StudentProfile.GetNationIdByStudentId(ctx, request.Payload.StudentId)
	if err != nil {
		return updateStudentFaceVideoResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !exist {
		return updateStudentFaceVideoResponse{Code: model.StatusNotFound, Message: "NOT_FOUND"}, errors.New("student_id does not exist")
	}
	// save face-video to DB
	uri, err := gcloud.SaveFaceVideoFile(ctx, request.Payload.StudentId, request.Payload.FileName, request.Payload.File)
	if err != nil {
		return updateStudentFaceVideoResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	// save face-thumbnail to DB
	uri_thumbnail, err := gcloud.SaveFaceThumnailFile(ctx, request.Payload.StudentId, request.Payload.FileName, request.Payload.FileThumbnail)
	if err != nil {
		return updateStudentFaceVideoResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	// update link photo to DB
	var (
		videoPath    = fmt.Sprintf("/ekyc_image_bucket/%s", uri)
		videoURL     = fmt.Sprintf("https://storage.googleapis.com/ekyc_image_bucket/%s", uri)
		thumbnailURL = fmt.Sprintf("https://storage.googleapis.com/ekyc_image_bucket/%s", uri_thumbnail)
	)
	log.Println("uri_thumbnail", uri_thumbnail)
	log.Println("thumbnailURL", thumbnailURL)
	if err := fsdb.StudentProfile.SetFaceVideoURL(ctx, doc_id, videoURL, thumbnailURL); err != nil {
		return updateStudentFaceVideoResponse{Code: 55, Message: err.Error()}, err
	}

	return updateStudentFaceVideoResponse{
		Payload: student_face_video_resp{
			URL:  videoURL,
			Path: videoPath,
		},
	}, nil
}

/* */
func __studentDetail(ctx context.Context, request *studentDetailRequest) (studentDetailResponse, error) {
	email, full_name, phone_number, national_id,
		birthday, sex, address, address_origin, unit_id, image, image_ekyc, modified_by, created_by,
		modified_at, created_at, exist, err := fsdb.StudentProfile.GetByStudentId(ctx, request.StudentId)
	if err != nil {
		return studentDetailResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !exist {
		return studentDetailResponse{Code: model.StatusNotFound, Message: "NOT FOUND"}, errors.New("student_id does not exist")
	}

	return studentDetailResponse{
		Payload: student_detail_data{
			StudentId:     request.StudentId,
			Email:         email,
			FullName:      full_name,
			PhoneNumber:   phone_number,
			NationalId:    national_id,
			Birthday:      birthday,
			Sex:           sex,
			Address:       address,
			AddressOrigin: address_origin,
			UnitId:        unit_id,
			Image:         image,
			ImageEkyc:     image_ekyc,
			ModifiedBy:    modified_by,
			ModifiedAt:    modified_at,
			CreatedBy:     created_by,
			CreatedAt:     created_at,
		}}, nil
}

/* */
func __updateStudent(ctx context.Context, request *updateStudentRequest) (updateStudentResponse, error) {
	// if err := request.Payload.validate(); err != nil {
	// 	return updateStudentResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }
	// validate fields
	student_id_already_exist, err := fsdb.StudentProfile.ValidateStudentId(ctx, request.Payload.StudentId)
	if err != nil {
		return updateStudentResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !student_id_already_exist {
		return updateStudentResponse{Code: model.StatusDataNotFound, Message: "STUDENT_ID_NOT_EXIST"}, errors.New("student_id not exist")
	}
	// check CCCD - Face Id

	if err := fsdb.StudentProfile.SetIsBlocked(ctx, request.Payload.StudentId, request.Payload.IsBlocked); err != nil {
		return updateStudentResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	return updateStudentResponse{}, nil
}

/* */
func __uploadFaceImage(ctx context.Context, request *uploadFaceImageRequest) (uploadFaceImageResponse, error) {
	if err := request.Payload.validate(); err != nil {
		return uploadFaceImageResponse{
			Code:    model.StatusBadRequest,
			Message: err.Error()}, err
	}
	doc_id, _, _, ok, err := fsdb.StudentProfile.GetNationIdByStudentId(ctx, request.Payload.StudentId)
	if err != nil {
		return uploadFaceImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !ok {
		return uploadFaceImageResponse{Code: model.StatusNotFound, Message: "NOT_FOUND"}, errors.New("student_id does not exist")
	}
	// save image to DB
	uri, err := gcloud.SaveFaceImageFile(ctx, request.Payload.StudentId, request.Payload.FileName, request.Payload.File)
	if err != nil {
		return uploadFaceImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	// update link photo to DB
	var (
		photoPath = fmt.Sprintf("/ekyc_image_bucket/%s", uri)
	)
	if err := fsdb.StudentProfile.SetFaceImageURL(ctx, doc_id, photoPath); err != nil {
		return uploadFaceImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	return uploadFaceImageResponse{
		Payload: face_image_resp{
			Path: photoPath,
		},
	}, nil
}

/* */
func __djangoNationalIdImage(ctx context.Context, request *uploadNationalIdImageRequest) (uploadNationalIdImageResponse, error) {
	// if err := request.Payload.validate(); err != nil {
	// 	return uploadNationalIdImageResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }

	return uploadNationalIdImageResponse{
		Payload: national_id_image_resp{
			FullName:          "Le Dinh Trieu",
			NationalId:        "123456789",
			DateOfBirth:       time.Now(),
			DateOfExpiry:      time.Now(),
			Gender:            "Nam",
			Address:           "TP HCM",
			PlaceOfOrigin:     "TP HCM",
			Nationality:       "Viet Nam",
			NationalIdCardURL: "https://www.consilium.europa.eu/prado/images/EST-JO-06001/334716.jpg",
		},
	}, nil
}

/* */
func __djangoFaceVideo(ctx context.Context, request *uploadNationalIdImageRequest) (uploadNationalIdImageResponse, error) {
	// if err := request.Payload.validate(); err != nil {
	// 	return uploadNationalIdImageResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }

	// var (
	// 	student_id = request.Payload.StudentId
	// 	face_id    = primitive.NewObjectID().Hex()
	// )

	// go func(){

	// }

	return uploadNationalIdImageResponse{
		Payload: national_id_image_resp{
			FullName:          "Le Dinh Trieu",
			NationalId:        "123456789",
			DateOfBirth:       time.Now(),
			DateOfExpiry:      time.Now(),
			Gender:            "Nam",
			Address:           "TP HCM",
			PlaceOfOrigin:     "TP HCM",
			Nationality:       "Viet Nam",
			NationalIdCardURL: "https://www.consilium.europa.eu/prado/images/EST-JO-06001/334716.jpg",
		},
	}, nil
}

/* */
func __uploadNationalIdImage(ctx context.Context, request *uploadNationalIdImageRequest) (uploadNationalIdImageResponse, error) {
	doc_id, _, _, exist, err := fsdb.StudentProfile.GetNationIdByStudentId(ctx, request.Payload.StudentId)
	if err != nil {
		return uploadNationalIdImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !exist {
		return uploadNationalIdImageResponse{Code: model.StatusNotFound, Message: "NOT_FOUND"}, errors.New("student_id does not exist")
	}
	// save image to DB
	uri, err := gcloud.SaveEkycImageFile(ctx, request.Payload.StudentId, request.Payload.FileName, request.Payload.File)
	if err != nil {
		return uploadNationalIdImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	// update link photo to DB
	var (
		//photoPath = fmt.Sprintf("/ekyc_image_bucket/%s", uri)
		photoURL = fmt.Sprintf("https://storage.googleapis.com/ekyc_image_bucket/%s", uri)
	)
	if err := fsdb.StudentProfile.SetNationalIdImageURL(ctx, doc_id, photoURL); err != nil {
		return uploadNationalIdImageResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}

	return uploadNationalIdImageResponse{
		Payload: national_id_image_resp{
			NationalIdCardURL: photoURL,
		},
	}, nil
}

/* */
func __uploadFaceRegImage(ctx context.Context, request *uploadFaceRegImageRequest) (uploadFaceRegImageResponse, error) {
	// if err := request.Payload.validate(); err != nil {
	// 	return uploadNationalIdImageResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }

	return uploadFaceRegImageResponse{
		Payload: face_reg_image_resp{
			FaceImageURL: "https://png.pngtree.com/png-clipart/20190924/original/pngtree-user-vector-avatar-png-image_4830521.jpg",
			PersonId:     "ab123456",
			StudentId:    "123",
			FullName:     "Le Dinh Trieu",
		},
	}, nil
}

/* */
func __filterListSession(ctx context.Context,
	request *filterListSessionRequest) (
	filterListSessionResponse, error) {
	db_sessions, err := fsdb.AuthSession.GetAll(ctx)
	if err != nil {
		return filterListSessionResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	var (
		list_session = make([]session_data, 0)
	)
	for _, session := range db_sessions {
		list_session = append(list_session, withAuthSessionModel(session))
	}
	return filterListSessionResponse{
		Payload: list_session_resp{
			TotalSession: len(list_session),
			ListSession:  list_session,
		},
	}, nil
}

/* */
func __filterListDevice(ctx context.Context,
	request *filterListDeviceRequest) (
	filterListDeviceResponse, error) {
	db_devices, err := fsdb.DeviceProfile.GetAll(ctx)
	if err != nil {
		return filterListDeviceResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	var (
		list_device = make([]device_data, 0)
	)
	for _, device := range db_devices {
		list_device = append(list_device, withDeviceModel(device))
	}
	return filterListDeviceResponse{
		Payload: list_device_resp{
			TotalDevice: len(list_device),
			ListDevice:  list_device,
		},
	}, nil
}

/* */
func __deviceDetail(ctx context.Context, request *deviceDetailRequest) (deviceDetailResponse, error) {
	terminal_name, avt, is_blocked, modified_by, created_by,
		last_login_at, modified_at, created_at, exist, err := fsdb.DeviceProfile.GetByTerminalId(ctx, request.TerminalId)

	if err != nil {
		return deviceDetailResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !exist {
		return deviceDetailResponse{Code: model.StatusNotFound, Message: "NOT FOUND"}, errors.New("terminal_id does not exist")
	}

	return deviceDetailResponse{
		Payload: device_detail_data{
			TerminalId:   request.TerminalId,
			TerminalName: terminal_name,
			Avatar:       avt,
			IsBlocked:    is_blocked,
			LastLoginAt:  last_login_at,
			ModifiedBy:   modified_by,
			ModifiedAt:   modified_at,
			CreatedBy:    created_by,
			CreatedAt:    created_at,
		}}, nil
}

/* */
func __updateDevice(ctx context.Context, request *updateDeviceRequest) (updateDeviceResponse, error) {
	// if err := request.Payload.validate(); err != nil {
	// 	return updateStudentResponse{
	// 		Code:    model.StatusBadRequest,
	// 		Message: err.Error()}, err
	// }
	// validate fields
	terminal_id_already_exist, err := fsdb.DeviceProfile.ValidateTerminalId(ctx, request.Payload.TerminalId)
	if err != nil {
		return updateDeviceResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	if !terminal_id_already_exist {
		return updateDeviceResponse{Code: model.StatusDataNotFound, Message: "TERMINAL_ID_NOT_EXIST"}, errors.New("terminal_id not exist")
	}

	if err := fsdb.DeviceProfile.SetIsBlocked(ctx, request.Payload.TerminalId, request.Payload.IsBlocked); err != nil {
		return updateDeviceResponse{
			Code:    model.StatusServiceUnavailable,
			Message: err.Error()}, err
	}

	return updateDeviceResponse{}, nil
}

/* */

func __mockAuthSession(ctx context.Context, request *mockAuthSessionRequest) (mockAuthSessionResponse, error) {

	var (
		session_id  = primitive.NewObjectID().Hex()
		student_id  = request.Payload.StudentId
		face_id     = "123456a"
		full_name   = "Le Dinh Trieu"
		terminal_id = "123456a"
		unit_id     = "d18cqcp02-n"
		image_url   = "https://png.pngtree.com/png-clipart/20190924/original/pngtree-user-vector-avatar-png-image_4830521.jpg"
		auth_at     = time.Date(2023, 04, request.Payload.Date, request.Payload.Hour, 30, 15, 00, time.Local)
	)
	_, err := fsdb.AuthSession.Add(ctx, session_id, student_id, face_id, terminal_id, full_name, unit_id, image_url, auth_at)
	if err != nil {
		return mockAuthSessionResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}

	return mockAuthSessionResponse{}, nil
}

func __reportAuthSession(ctx context.Context, request *reportAuthSessionRequest) (reportAuthSessionResponse, error) {
	var (
		student_id = request.Payload.StudentId
		month      = request.Payload.Month
		year       = 2023
	)
	log.Println("student_id", student_id)
	_, err := fsdb.AuthSession.ReportByMonth(ctx, student_id, month, year)
	if err != nil {
		return reportAuthSessionResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	return reportAuthSessionResponse{}, nil

}

func __filterReportAuthSession(ctx context.Context, request *filterReportAuthSessionRequest) (filterReportAuthSessionResponse, error) {

	var (
		student_id = request.StudentId
		month      = request.Month
		year       = request.Year
	)
	log.Println("student_id", student_id)
	_, err := fsdb.AuthSession.ReportByMonth(ctx, student_id, int(month), year)
	if err != nil {
		return filterReportAuthSessionResponse{Code: model.StatusServiceUnavailable, Message: err.Error()}, err
	}
	return filterReportAuthSessionResponse{}, nil

}
