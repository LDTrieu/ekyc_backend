package mobile

import "errors"

type traceField struct {
	RequestId string `json:"reqId"`
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
