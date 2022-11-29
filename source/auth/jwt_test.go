package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func Test_JWT_Full(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	t.Logf("STEP_1: GenerateJWTLoginSession")
	id, login_access, err := GenerateJWTLoginSession(
		ctx, uuid.NewString(), "1234567890")
	if err != nil {
		t.Errorf("GenerateJWTLoginSession: %+v\n", err)
		return
	}
	data, status, err := ValidateLoginSessionJWT(ctx, login_access.Token)
	if err != nil {
		t.Errorf("ValidateLoginJWT: %+v\n", err)
		return
	} else {
		t.Logf("STEP_1: ValidateLoginJWT ID= %+v\n", id)
		t.Logf("STEP_1: ValidateLoginJWT STATUS= %+v\n", status)
		t.Logf("STEP_1: ValidateLoginJWT DATA= %+v\n", data)
	}
}
