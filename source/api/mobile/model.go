package mobile

import (
	"ekyc-app/source/auth"
	"errors"
	"time"
)

type traceField struct {
	RequestId string `json:"reqId"`
}

type Credentials struct {
	TerminalId string `json:"terminalId"`
	Password   string `json:"password"`
}

func (ins *Credentials) validate() error {
	if len(ins.TerminalId) < 1 {
		return errors.New("field terminal_id invalid")
	}
	if len(ins.Password) < 1 {
		return errors.New("field password invalid")
	}
	return nil
}

type submitQRLoginRequest struct {
	traceField
	Payload submit_qr_login_req `json:"payload"`
}

type submitQRLoginResponse struct {
	traceField
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Payload submit_qr_login_resp `json:"payload"`
}

type submit_qr_login_req struct {
	JWT       string `json:"jwt"`
	QrData    string `json:"qrData"`
	AccountID string `json:"accId"`
}

func (ins *submit_qr_login_req) validate() error {
	if len(ins.QrData) < 5 {
		ins.QrData = ins.JWT
	}
	if len(ins.QrData) < 5 {
		return errors.New("field jwt is required")
	}
	if len(ins.AccountID) < 3 {
		return errors.New("field accId is required")
	}
	return nil
}

type submit_qr_login_resp struct {
	AccountID string `json:"accId"`
}

/* */
type loginTerminalRequest struct {
	traceField
	Credentials
}

type loginTerminalResponse struct {
	traceField
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Payload login_terminal_data `json:"payload"`
}

type login_terminal_data struct {
	TerminalId string `json:"terminalId"`
	// SerialNumber string    `json:"serialNumber"`
	// Location     string    `json:"location"`
	Avt         string    `json:"avt"`
	LinkedAt    time.Time `json:"linkedAt"`
	Description string    `json:"description"`
	Token       string    `json:"token"`
}

/* */
type signupTerminalRequest struct {
	traceField
	Credentials
	Permit      *auth.DataJWT
	Avt         string `json:"avt"`
	Description string `json:"description"`
}

type signupTerminalResponse struct {
	traceField
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ins *signupTerminalRequest) validate() error {
	if err := ins.Credentials.validate(); err != nil {
		return err
	}

	return nil
}
