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
