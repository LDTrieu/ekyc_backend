package portal

import "errors"

const (
	qrActionLoginWebPortal = "login-web-portal"
	// qrActionLinkDevice     = "link-business-device"
)

type traceField struct {
	RequestId string `json:"reqId"`
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
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
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
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Token    string `json:"token"`
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
