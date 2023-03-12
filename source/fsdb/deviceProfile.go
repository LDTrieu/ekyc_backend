package fsdb

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

type deviceProfileFs struct {
	coll                string
	fieldTerminalId     string
	fieldAvt            string
	fieldHashedPassword string
	fieldToken          string
	fieldIsBlocked      string
	fieldLastLoginAt    string
	fieldModifiedAt     string
	fieldCreatedAt      string
}

var DeviceProfile = &deviceProfileFs{
	coll:                "device_profile",
	fieldTerminalId:     "terminal_id",
	fieldAvt:            "avt",
	fieldHashedPassword: "hashed_password",
	fieldToken:          "token",
	fieldIsBlocked:      "isBlocked",
	fieldLastLoginAt:    "last_login_at",
	fieldModifiedAt:     "modified_at",
	fieldCreatedAt:      "created_at",
}

type DeviceProfileModel struct {
	TerminalId     string    `json:"terminalId" firestore:"terminal_id"`
	Avatar         string    `json:"avt" firestore:"avt"`
	Token          string    `json:"token" firestore:"token"`
	HashedPassword string    `json:"hashedPassword" firestore:"hashed_password"`
	CreatedBy      string    `json:"createdBy" firestore:"created_by"`
	IsBlocked      bool      `json:"isBlocked" firestore:"is_blocked"`
	LastLoginAt    time.Time `json:"lastLoginDate" firestore:"last_login_at"`
	ModifiedAt     time.Time `json:"modifiedDate" firestore:"modified_at"`
	CreatedAt      time.Time `json:"createdDate" firestore:"created_at"`
}

func (ins *deviceProfileFs) Add(
	ctx context.Context, terminal_id, avt, hashed_password,
	create_by string) (id string, err error) {
	init := DeviceProfileModel{
		TerminalId:     terminal_id,
		Avatar:         avt,
		HashedPassword: hashed_password,
		IsBlocked:      false,
		CreatedBy:      create_by,
		LastLoginAt:    time.Now(),
		ModifiedAt:     time.Now(),
		CreatedAt:      time.Now(),
	}
	return add(ctx, ins.coll, init)
}

// CheckLogin
func (ins *deviceProfileFs) CheckLogin(ctx context.Context, terminal_id, password string) (doc_id string, is_blocked bool, err error) {
	var (
		temp DeviceProfileModel
	)
	doc_id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldTerminalId, terminal_id)
	if err != nil {
		return "", true, errors.New("terminal_id or password invalid")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(temp.HashedPassword), []byte(password)); err != nil {
		return "", true, errors.New("terminal_id or password invalid")
	}
	if temp.IsBlocked {
		return doc_id, true, nil
	}
	return doc_id, false, nil
}
func (ins *deviceProfileFs) SetToken(ctx context.Context, docId, token string, linked_at time.Time) error {
	var update = map[string]interface{}{
		ins.fieldLastLoginAt: linked_at,
		ins.fieldToken:       token,
	}
	return updateFields(ctx, docId, ins.coll, update)
}

func (ins *deviceProfileFs) ValidateTerminalId(ctx context.Context,
	terminal_id string) (already_exist bool, err error) {

	var (
		count int
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.
			Where(ins.fieldTerminalId, Equal, terminal_id).
			Limit(1).
			Documents(ctx)
		for {
			_, err := dIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			count++
		}
		return nil
	}); err != nil {
		return false, err
	}
	return count > 0, nil
}
