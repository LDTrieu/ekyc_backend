package faceauth

import (
	"context"
	"log"
	"testing"
)

func Test_AddFace(t *testing.T) {
	ctx := context.Background()
	var (
		request = &AddFaceRequest{
			Name:     "long4",
			FaceId:   "123456",
			VideoURL: "https://storage.googleapis.com/ekyc_image_bucket/n18dccn241/face_video_n18dccn241",
		}
	)

	response, err := AddFace(ctx, request)
	if err != nil {
		t.Fatal("ERR: ", err)
	}
	log.Println("OKE", response.Mesage)
}
