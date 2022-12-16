package fsdb

import (
	"context"
	"log"
	"testing"
)

func Test_GetByStudentId(t *testing.T) {
	ctx := context.Background()
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, err := StudentProfile.GetByStudentId(ctx, "n18dccn229")
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Println(StudentProfile.GetByStudentId(ctx, "n18dccn229"))
	//t.Fatal("OK")
}

func Test_SetFaceImageURL(t *testing.T) {
	ctx := context.Background()
	if err := StudentProfile.SetFaceImageURL(ctx, "L3JMMGWlu8RmlH5KyTZu", "photoURL"); err != nil {
		t.Fatal("ERR: ", err)
	}
	//t.Fatal("OK")
}
