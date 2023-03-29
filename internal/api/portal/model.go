package portal

import (
	"ekyc-app/internal/auth"
	"ekyc-app/internal/fsdb"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	qrActionLoginWebPortal = "login-web-portal"
	// qrActionLinkDevice     = "link-business-device"
)

type traceField struct {
	RequestId string `json:"reqId"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* */
type rendQRLoginRequest struct {
	traceField
}

type rendQRLoginResponse struct {
	traceField
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Payload rend_qr_code_login_resp `json:"payload"`
}
type rend_qr_code_login_resp struct {
	Path      string `json:"path"`
	JWT       string `json:"jwt"`
	ExpiresIn int64  `json:"expiresIn"`
	IssuedAt  int64  `json:"issuedAt"`
}

/* */
type loginBasicRequest struct {
	traceField
	Credentials
}

type loginBasicResponse struct {
	traceField
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Payload login_basic_data `json:"payload"`
}

type login_basic_data struct {
	AccountId   string    `json:"accountId"`
	FullName    string    `json:"fullName"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	UnitId      string    `json:"unitId"`
	Birthday    time.Time `json:"birthday"`
	Avt         string    `json:"avt"`
	Banner      string    `json:"banner"`
	Token       string    `json:"token"`
}

func (ins *Credentials) validate() error {
	if len(ins.Email) < 1 {
		return errors.New("field email invalid")
	}
	if len(ins.Password) < 1 {
		return errors.New("field password invalid")
	}

	return nil
}

/* */
type signupBasicRequest struct {
	traceField
	Credentials
	Permit      *auth.DataJWT
	FullName    string    `json:"fullName"`
	UnitId      string    `json:"unitId"`
	PhoneNumber string    `json:"phoneNumber"`
	DateOfBirth time.Time `json:"dateOfBirth"`

	// some fields
}

type signupBasicResponse struct {
	traceField
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Payload signup_basic_data `json:"payload"`
}

type signup_basic_data struct {
	AccountId string `json:"accountId"`
	Email     string `json:"email"`
}

func (ins *signupBasicRequest) validate() error {
	if err := ins.Credentials.validate(); err != nil {
		return err
	}
	if len(ins.PhoneNumber) < 1 {
		return errors.New("phone_number invalid")
	}
	if len(ins.FullName) < 1 {
		return errors.New("full_name invalid")
	}
	if len(ins.UnitId) < 1 {
		return errors.New("unit_id invalid")
	}
	return nil
}

/* */
type filterListUserRequest struct {
	traceField
	Permit *auth.DataJWT
	//	Payload list_user_req `json:"payload"`
}

type filterListUserResponse struct {
	traceField
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Payload list_user_resp `json:"payload"`
}

type list_user_resp struct {
	TotalUser int         `json:"totalUser"`
	ListUser  []user_data `json:"listUser"`
}

type user_data struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	AccountId   string `json:"accountId"`
	UnitId      string `json:"unitId"`
	PhoneNumber string `json:"phoneNumber"`
	IsBlocked   bool   `json:"isBlocked"`
	CreatedBy   string `json:"createdBy"`
}

func withUserModel(um *fsdb.PersonProfileModel) user_data {
	return user_data{
		FullName:    um.FullName,
		PhoneNumber: um.PhoneNumber,
		Email:       um.Email,
		AccountId:   um.AccountId,
		UnitId:      "um.UnitId",
		IsBlocked:   um.IsBlocked,
		CreatedBy:   "um.CreatedBy",
	}
}

/* */
type userDetailRequest struct {
	traceField
	Permit    *auth.DataJWT
	AccountId string
	// Payload   student_profile_data `json:"payload"`
}

type userDetailResponse struct {
	traceField
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Payload user_detail_data `json:"payload"`
}

type user_detail_data struct {
	AccountId   string `json:"accountId"`
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	UnitId      string `json:"unitId"`
	PhoneNumber string `json:"phoneNumber"`
	IsBlocked   bool   `json:"isBlocked"`
	// ModifiedBy  string    `json:"modifiedBy"`
	LastLoginAt time.Time `json:"lastLoginAt"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

/* */
type filterListStudentRequest struct {
	traceField
	Permit *auth.DataJWT
	//	Payload list_user_req `json:"payload"`
}

type filterListStudentResponse struct {
	traceField
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Payload list_student_resp `json:"payload"`
}

type list_student_resp struct {
	Stash        int            `json:"stash"`
	TotalStudent int            `json:"totalStudent"`
	ListStudent  []student_data `json:"listStudent"`
}

type student_data struct {
	//student_detail
	FullName    string    `json:"fullName"`
	Image       string    `json:"image"`
	StudentId   string    `json:"studentId"`
	PhoneNumber string    `json:"phoneNumber"`
	UnitId      string    `json:"unitId" `
	Birthday    time.Time `json:"birthday"`
	Gender      string    `json:"gender"`

	//student_ekyc

	//student_censorship
	IsBlocked bool `json:"isBlocked"`
}

func withStudentModel(sm *fsdb.StudentProfileModel) student_data {
	return student_data{
		//student_detail
		FullName:    sm.FullName,
		StudentId:   sm.StudentId,
		Image:       "https://tuk-cdn.s3.amazonaws.com/assets/components/advance_tables/at_1.png",
		PhoneNumber: sm.PhoneNumber,
		Birthday:    sm.Birthday,
		UnitId:      sm.UnitId,
		Gender:      sm.Sex,
		//student_ekyc

		//student_censorship
		IsBlocked: sm.IsBlocked,
	}
}

/* */
type createStudentProfileRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload student_profile_data `json:"payload"`
}

type createStudentProfileResponse struct {
	traceField
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Payload list_student_resp `json:"payload"`
}

type student_profile_data struct {
	StudentId   string    `json:"studentId"`
	Email       string    `json:"email"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	FullName    string    `json:"fullName"`
	PhoneNumber string    `json:"phoneNumber"`
	UnitId      string    `json:"unitId"`
	Birthday    time.Time `json:"dateOfBirth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	NationalId  string    `json:"nationalId"`
	Hometown    string    `json:"hometown"`
	Image       string    `json:"image" `
	ImageEkyc   string    `json:"imageEkyc"`
	IsBlocked   bool      `json:"isBlocked"`
}

func (ins *student_profile_data) validate() error {
	if len(ins.StudentId) < 1 {
		return errors.New("field student_id is invalid")
	}
	if len(ins.Email) < 1 {
		return errors.New("field Email is invalid")
	}
	if len(ins.FirstName) < 1 {
		return errors.New("field first_name is invalid")
	}
	if len(ins.LastName) < 1 {
		return errors.New("field last_name is invalid")
	}
	if len(ins.UnitId) < 1 {
		return errors.New("field unit_id is invalid")
	}
	// if len(ins.NationalId) < 1 {
	// 	return errors.New("field national_id is invalid")
	// }
	if len(ins.Address) < 1 {
		return errors.New("field address is invalid")
	}
	// if len(ins.Image) < 1 {
	// 	return errors.New("field image is invalid")
	// }
	// if len(ins.ImageEkyc) < 1 {
	// 	return errors.New("field image_ekyc is invalid")
	// }

	return nil
}

/* */
type studentDetailRequest struct {
	traceField
	Permit    *auth.DataJWT
	StudentId string
	// Payload   student_profile_data `json:"payload"`
}

type studentDetailResponse struct {
	traceField
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Payload student_detail_data `json:"payload"`
}

type student_detail_data struct {
	StudentId     string    `json:"studentId"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	PhoneNumber   string    `json:"phoneNumber"`
	UnitId        string    `json:"unitId"`
	NationalId    string    `json:"nationalId"`
	Birthday      time.Time `json:"birthday"`
	Sex           string    `json:"sex"`
	Address       string    `json:"address"`
	AddressOrigin string    `json:"addressOrigin"`
	Image         string    `json:"image" `
	ImageEkyc     string    `json:"imageEkyc"`
	IsBlocked     bool      `json:"isBlocked"`
	ModifiedBy    string    `json:"modifiedBy"`
	ModifiedAt    time.Time `json:"modifiedAt"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

/* */
type uploadFaceImageRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload face_image_req `json:"payload"`
}

type uploadFaceImageResponse struct {
	traceField
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Payload face_image_resp `json:"payload"`
}

type face_image_req struct {
	StudentId string `json:"studentId"`
	FileName  string `json:"fileName"`
	File      []byte `json:"file"`
}

type face_image_resp struct {
	URL  string `json:"url"`
	Path string `json:"photoPath"`
}

func (ins *face_image_req) validate() error {
	if len(ins.StudentId) < 1 {
		return errors.New("student_id invalids")
	}
	if len(ins.FileName) < 3 {
		ins.FileName = fmt.Sprintf("%s_%s.bin", ins.StudentId, primitive.NewObjectID().Hex())
	}
	// if len(ins.File) < 128 {
	// 	return errors.New("file is too small")
	// }
	return nil
}

/* */
type uploadNationalIdImageRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload national_id_image_req `json:"payload"`
}

type uploadNationalIdImageResponse struct {
	traceField
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Payload national_id_image_resp `json:"payload"`
}
type national_id_image_req struct {
	StudentId string `json:"studentId"`
	FileName  string `json:"fileName"`
	File      []byte `json:"file"`
}

type national_id_image_resp struct {
	FullName      string    `json:"fullName"`
	NationalId    string    `json:"nationalId"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	DateOfExpiry  time.Time `json:"dateOfExpiry"`
	Gender        string    `json:"gender"`
	Address       string    `json:"address"`
	PlaceOfOrigin string    `json:"placeOfOrigin"`
	Nationality   string    `json:"nationality"`
	//FaceImage string `json:"faceImage"`
	NationalIdCardURL string `json:"nationalIdCardURL"`
	// Path string `json:"photoPath"`
}

/* */
type uploadFaceRegImageRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload face_reg_image_req `json:"payload"`
}

type uploadFaceRegImageResponse struct {
	traceField
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Payload face_reg_image_resp `json:"payload"`
}
type face_reg_image_req struct {
	StudentId string `json:"studentId"`
	FileName  string `json:"fileName"`
	File      []byte `json:"file"`
}

type face_reg_image_resp struct {
	PersonId     string `json:"personId"`
	StudentId    string `json:"studentId"`
	FullName     string `json:"fullName"`
	FaceImageURL string `json:"faceImageURL"`

	//FaceImage string `json:"faceImage"`
}

/* */
type updateStudentEkycRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload update_ekyc_data `json:"payload"`
}
type updateStudentEkycResponse struct {
	traceField
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type update_ekyc_data struct {
	StudentId         string    `json:"studentId"`
	PersonId          string    `json:"personId"`
	NationalId        string    `json:"nationalId"`
	FullName          string    `json:"fullName"`
	Gender            string    `json:"gender"`
	Address           string    `json:"address"`
	PlaceOfOrigin     string    `json:"placeOfOrigin"`
	Nationality       string    `json:"nationality"`
	DateOfBirth       time.Time `json:"dateOfBirth"`
	DateOfExpiry      time.Time `json:"dateOfExpiry"`
	NationalIdCardURL string    `json:"nationalIdCardURL"`
	FaceImageURL      string    `json:"faceImageURL"`
}
type update_ekyc_resp struct {
}

func (ins *update_ekyc_data) validate() error {
	if len(ins.StudentId) < 1 {
		return errors.New("field student_id is invalid")
	}

	if len(ins.NationalId) < 1 {
		return errors.New("field national_id is invalid")
	}
	if len(ins.FullName) < 1 {
		return errors.New("field full_name is invalid")
	}
	if len(ins.Gender) < 1 {
		return errors.New("field gender is invalid")
	}
	if len(ins.Address) < 1 {
		return errors.New("field address is invalid")
	}
	if len(ins.PlaceOfOrigin) < 1 {
		return errors.New("field place_of_origin is invalid")
	}
	if len(ins.NationalIdCardURL) < 1 {
		return errors.New("field national_id_card_url is invalid")
	}
	if len(ins.FaceImageURL) < 1 {
		return errors.New("field face_image_url is invalid")
	}

	return nil
}

/* */
//updateStudentRequest
type updateStudentRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload update_student_data `json:"payload"`
}
type updateStudentResponse struct {
	traceField
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type update_student_data struct {
	StudentId string `json:"studentId"`
	IsBlocked bool   `json:"isBlocked"`
}

/* */
type filterListSessionRequest struct {
	traceField
	Permit *auth.DataJWT
	//Payload list_session_req `json:"payload"`
}

type filterListSessionResponse struct {
	traceField
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Payload list_session_resp `json:"payload"`
}
type list_session_req struct {
	// Stash        int            `json:"stash"`
	// TotalSession int            `json:"totalSession"`
	// ListSession  []session_data `json:"listSession"`
}

type list_session_resp struct {
	Stash        int            `json:"stash"`
	TotalSession int            `json:"totalSession"`
	ListSession  []session_data `json:"listSession"`
}

type session_data struct {
	SessionId  string    `json:"sessionId"`
	StudentId  string    `json:"studentId"`
	FullName   string    `json:"fullName"`
	FaceId     string    `json:"faceId"`
	TerminalId string    `json:"terminalId"`
	UnitId     string    `json:"unitId"`
	ImageUrl   string    `json:"imageUrl"`
	AuthAt     time.Time `json:"authAt"`
}

func withAuthSessionModel(sm *fsdb.AuthSessionModel) session_data {
	return session_data{
		SessionId:  sm.SessionId,
		StudentId:  sm.StudentId,
		FullName:   sm.FullName,
		FaceId:     sm.FaceId,
		TerminalId: sm.TerminalId,
		UnitId:     sm.UnitId,
		ImageUrl:   "https://tuk-cdn.s3.amazonaws.com/assets/components/advance_tables/at_1.png",
		AuthAt:     sm.AuthAt,
	}
}

/* */
type filterListDeviceRequest struct {
	traceField
	Permit *auth.DataJWT
	//Payload list_session_req `json:"payload"`
}

type filterListDeviceResponse struct {
	traceField
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Payload list_device_resp `json:"payload"`
}

type list_device_resp struct {
	Stash       int           `json:"stash"`
	TotalDevice int           `json:"totalDevice"`
	ListDevice  []device_data `json:"listDevice"`
}
type device_data struct {
	TerminalId   string    `json:"terminalId"`
	TerminalName string    `json:"terminalName"`
	Avatar       string    `json:"avt"`
	IsBlocked    bool      `json:"isBlocked"`
	ModifiedBy   string    `json:"modifiedBy"`
	LastLoginAt  time.Time `json:"lastLoginDate"`
	ModifiedAt   time.Time `json:"modifiedDate"`
}

func withDeviceModel(dm *fsdb.DeviceProfileModel) device_data {
	return device_data{
		TerminalId:   dm.TerminalId,
		TerminalName: dm.TerminalName,
		Avatar:       dm.Avatar,
		IsBlocked:    dm.IsBlocked,
		LastLoginAt:  dm.LastLoginAt,
		ModifiedBy:   dm.ModifiedBy,
		ModifiedAt:   dm.ModifiedAt,
	}
}
