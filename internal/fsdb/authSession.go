package fsdb

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type authSessionFs struct {
	coll            string
	fieldSessionId  string
	fieldStudentId  string
	fieldFaceId     string
	fieldTerminalId string
	fieldFullName   string
	fieldUnitId     string
	fieldImageUrl   string
	fieldAuthAt     string
}

var AuthSession = &authSessionFs{
	coll:            "auth_session",
	fieldSessionId:  "session_id",
	fieldStudentId:  "student_id",
	fieldFaceId:     "face_id",
	fieldTerminalId: "terminal_id",
	fieldFullName:   "full_name",
	fieldUnitId:     "unit_id",
	fieldImageUrl:   "image_url",
	fieldAuthAt:     "auth_at",
}

type AuthSessionModel struct {
	SessionId  string    `json:"sessionId" firestore:"session_id"`
	StudentId  string    `json:"studentId" firestore:"student_id"`
	FaceId     string    `json:"faceId" firestore:"face_id"`
	TerminalId string    `json:"terminalId" firestore:"terminal_id"`
	FullName   string    `json:"fullName" firestore:"full_name"`
	UnitId     string    `json:"unitId" firestore:"unit_id"`
	ImageUrl   string    `json:"imageUrl" firestore:"image_url"`
	AuthAt     time.Time `json:"authAt" firestore:"auth_at"`
}

func (ins *authSessionFs) Add(ctx context.Context,
	session_id, student_id, face_id, terminal_id, full_name, unit_id, image_url string, auth_at time.Time) (id string, err error) {
	var item = AuthSessionModel{
		SessionId:  session_id,
		StudentId:  student_id,
		FaceId:     face_id,
		TerminalId: terminal_id,
		FullName:   full_name,
		UnitId:     unit_id,
		ImageUrl:   image_url,
		AuthAt:     auth_at,
	}
	return add(ctx, ins.coll, item)
}

func (ins *authSessionFs) GetAll(ctx context.Context) (
	[]*AuthSessionModel, error) {
	var (
		list = make([]*AuthSessionModel, 0)
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.Documents(ctx)
		for {
			doc, err := dIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var temp AuthSessionModel
			if err := doc.DataTo(&temp); err != nil {
				return err
			}
			list = append(list, &temp)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}
