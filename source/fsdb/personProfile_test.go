package fsdb

import (
	"context"
	"testing"
)

func Test_Add(t *testing.T) {
	ctx := context.Background()
	accountId := "abc"
	sessionId := "123"
	token := "aaa"
	id, err := PersonProfile.Add(ctx, accountId, sessionId, token)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	t.Fatal("OKE: ", id)
}

func Test_GetByEmail(t *testing.T) {
	ctx := context.Background()
	email := "letrieu106@gmai.com"
	id, info, ok, err := PersonProfile.GetByEmail(ctx, email)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	t.Fatal("OKE: ", id, "  ", info, "  ", ok)

}

// GetByPhone
func Test_GetByPhone(t *testing.T) {
	ctx := context.Background()
	phone := "0948518286"
	id, info, ok, err := PersonProfile.GetByPhone(ctx, phone)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	t.Fatal("OKE: ", id, "  ", info, "  ", ok)

}

func Test_GetSessionID(t *testing.T) {
	ctx := context.Background()
	// accountId := "abc"
	// sessionId := "123"
	token := "aaa"
	id, sessionId, ok, err := PersonProfile.GetSessionID(ctx, token)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	t.Fatal("OKE: ", id, "  ", sessionId, "  ", ok)
}
