package portal

import (
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
