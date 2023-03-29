package fsdb

import (
	"context"
	"log"
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
	log.Println("OKE: ", id)
}

func Test_GetByEmail(t *testing.T) {
	ctx := context.Background()
	email := "letrieu106@gmai.com"
	id, info, ok, err := PersonProfile.GetByEmail(ctx, email)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Println("OKE: ", id, "  ", info, "  ", ok)

}

// GetByPhone
func Test_GetByPhone(t *testing.T) {
	ctx := context.Background()
	phone := "0948518286"
	id, info, ok, err := PersonProfile.GetByPhone(ctx, phone)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Println("OKE: ", id, "  ", info, "  ", ok)

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
	log.Println("OKE: ", id, "  ", sessionId, "  ", ok)
}

func Test_CheckLogin(t *testing.T) {
	ctx := context.Background()
	email := "letrieu106@gmail.com"
	//hashed_password := "$2a$08$uWs19KsmPgZ8LGZEcsDz6O7wt/AgScVrQ27qx3lUE0sf5kAPZcIuW"
	password := ""
	id, account_id, full_name, phone_number, birthday, unit_id, isBlocked, err := PersonProfile.CheckLogin(ctx, email, password)
	if err != nil {
		t.Fatal("ERR1: ", err)
	}
	log.Println("OKE: ", id, "  ", account_id, "  ", full_name, "  ", phone_number, "  ", birthday, unit_id, isBlocked)

}
