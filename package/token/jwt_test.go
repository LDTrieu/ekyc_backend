package token

import (
	"bytes"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_JWT_Full(t *testing.T) {
	var (
		secret     = bytes.NewBufferString("test_key")
		timeout    = 5 * time.Second
		input_data = struct {
			SessionID string `json:"sessionId"`
			AccountID string `json:"accountId"`
		}{
			SessionID: primitive.NewObjectID().Hex(),
			AccountID: "ACCOUNT_TEST",
		}
		output = struct {
			DocId     string `json:"docId"`
			SessionId string `json:"sessionId"`
			AccountID string `json:"accountId"`
			// LinkStatus int    `json:"linkStatus"`
		}{}
	)
	// Case 1:
	id, jwt, err := GenerateJWT(secret.Bytes(), timeout, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Case 1: GenerateJWT(NULL): %s\n-> %+v\n", id, jwt)
	status, err := ValidateJWT(secret.Bytes(), jwt.Token, &output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Case 1: ValidateJWT: %s\n-> %+v\n", status.String(), output)

	// Case 2:
	id, jwt, err = GenerateJWT(secret.Bytes(), timeout, input_data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("case 2: GenerateJWT: %s\n-> %+v\n", id, jwt)
	status, err = ValidateJWT(secret.Bytes(), jwt.Token, &output)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("case 2: ValidateJWT: %s\n-> %+v\n", status.String(), output)
	t.Fatal("OKE")
}
