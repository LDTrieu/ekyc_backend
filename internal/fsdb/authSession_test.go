package fsdb

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_AddAuthSession(t *testing.T) {
	ctx := context.Background()
	var (
		session_id  = uuid.NewString()
		student_id  = "n18dccn229"
		face_id     = "123"
		terminal_id = "1234"
		full_name   = "Le Dinh Trieu"
		unit_id     = "123"
		image_url   = "https://tuk-cdn.s3.amazonaws.com/assets/components/advance_tables/at_1.png"
		auth_at     = time.Now()
	)
	id, err := AuthSession.Add(ctx, session_id, student_id, face_id, terminal_id, full_name, unit_id, image_url, auth_at)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Fatal("OKE: ", id)
}
func Test_GetAllAuthSession(t *testing.T) {
	ctx := context.Background()

	list, err := AuthSession.GetAll(ctx)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Fatal("OKE: ", len(list))
}
