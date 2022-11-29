package auth

type DataJWT struct {
	SessionID string `json:"sessionId"`
	AccountID string `json:"accountId"`
}
