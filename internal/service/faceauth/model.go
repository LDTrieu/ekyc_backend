package faceauth

import "time"

type MockModel struct {
	Name string `json:"name" bson:"name"`
	Code int    `json:"code" bson:"code"`
}

type MockData struct {
	LastSign time.Time `json:"lastSign"`
}

/* */
type FaceAuthSessionModel struct {
	FileName string `json:"filename"`
	File     []byte `json:"file"`
}

/* */
type FaceAuthSessionResponse struct {
	Name   string `json:"name"`
	FaceId string `json:"faceId"`
}
