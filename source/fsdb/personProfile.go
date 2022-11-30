package fsdb

import (
	"context"
	"ekyc-app/source/model"
)

type personProfileFs struct {
	coll           string
	fieldAccountId string
	fieldSessionId string
	fieldToken     string
}

var PersonProfile = &personProfileFs{
	coll:           "person_profile",
	fieldAccountId: "account_id",
	fieldSessionId: "session_id",
	fieldToken:     "token",
}

type PersonProfileModel struct {
	AccountId string `json:"accountId" firestore:"account_id"`
	//DisplayName       string           `json:"displayName" firestore:"display_name"`
	//PhoneNumber       string           `json:"phoneNumber" firestore:"phone_number"`
	SessionId string `json:"sessionId" firestore:"session_id"`
	Token     string `json:"token" firestore:"token"`
}

func (ins *personProfileFs) GetSessionID(ctx context.Context, token string) (
	id string, session_id string, ok bool, err error) {
	var (
		temp PersonProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldToken, token)
	if err == model.ErrDocNotFound {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, err
	}
	return id, temp.SessionId, true, nil
}
