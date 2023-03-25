package fsdb

import (
	"context"
	"log"
	"testing"
	"time"
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

func Test_SetEkyc(t *testing.T) {
	var (
		ctx                  = context.Background()
		student_id           = "n18dccn242"
		national_id          = "123456"
		person_id            = "123456"
		full_name_ekyc       = "123456"
		gender               = "nam"
		face_image_url       = "123456"
		national_id_card_url = "123456"
		address_ekyc         = "123456"
		place_of_origin      = "123456"
		natinality           = "123456"
		modified_by          = "123456"
		date_of_birth        = time.Now()
		date_of_expiry       = time.Now()
	)
	if err := StudentProfile.SetEkyc(ctx, student_id, national_id, person_id, full_name_ekyc, face_image_url,
		gender, national_id_card_url, address_ekyc, place_of_origin, natinality, modified_by,
		date_of_birth, date_of_expiry); err != nil {
		log.Fatal(err)
	}
}
