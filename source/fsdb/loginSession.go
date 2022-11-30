package fsdb

import (
	"context"
	"ekyc-app/source/model"
	"time"
)

type loginSession struct {
	coll string `firebase:"-"`
}

type LoginSessionModel struct {
	SessionId string `firestore:"session_id"`
	JWT       string `firestore:"jwt"`
	ExpiresAt int64  `firestore:"expires_at"`
	CreatedAt int64  `firestore:"created_at"`
}

var LoginSessionDBC = &loginSession{
	coll: "cache_login_session",
}

func (me *loginSession) Add(ctx context.Context,
	sessionId, token string, expiresAt int64) (id string, err error) {
	var item = LoginSessionModel{
		SessionId: sessionId,
		JWT:       token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().Unix(),
	}
	return add(ctx, me.coll, item)
}

func (me *loginSession) GetToken(
	ctx context.Context,
	docID string) (
	token string, ok bool, err error) {
	var (
		temp LoginSessionModel
	)
	err = getById(ctx, &temp, me.coll, docID)
	if err == model.ErrDocNotFound {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return temp.JWT, true, nil
}
