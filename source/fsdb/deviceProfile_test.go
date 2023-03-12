package fsdb

import (
	"context"
	"log"
	"testing"
)

func Test_Create(t *testing.T) {
	ctx := context.Background()
	var (
		terminal_id     = "abc"
		hashed_password = "password_hash"
		avt             = "image1"
		create_by       = "admin1"
	)
	id, err := DeviceProfile.Add(ctx, terminal_id, avt, hashed_password, create_by)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Println("OKE: ", id)
}
