package portal

import (
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"errors"
	"time"
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
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Birthday    time.Time `json:"birthday"`
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
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
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
		return errors.New("field phone number invalid")
	}
	if len(ins.FullName) < 1 {
		return errors.New("field full name invalid")
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
	Image       string `json:"image"`
	PhoneNumber string `json:"phoneNumber"`
	IsBlocked   bool   `json:"isBlocked"`
}

func withUserModel(um *fsdb.PersonProfileModel) user_data {
	return user_data{
		FullName:    um.FullName,
		PhoneNumber: um.PhoneNumber,
		// Image     : ,
		IsBlocked: um.IsBlocked,
	}
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
	FullName  string `json:"fullName"`
	Image     string `json:"image"`
	StudentId string `json:"studentId"`
	IsBlocked bool   `json:"isBlocked"`
}

func withStudentModel(um *fsdb.StudentProfileModel) student_data {
	return student_data{
		FullName:  um.FullName,
		StudentId: um.StudentId,
		// Image     : ,
		IsBlocked: um.IsBlocked,
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
	FullName  string `json:"fullName"`
	Image     string `json:"image"`
	StudentId string `json:"studentId"`
	IsBlocked bool   `json:"isBlocked"`
}
