package portal

const (
	qrActionLoginWebPortal = "login-web-portal"
	// qrActionLinkDevice     = "link-business-device"
)

type traceField struct {
	RequestId string `json:"reqId"`
}
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
