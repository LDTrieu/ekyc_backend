package fsdb

import (
	"context"
	"ekyc-app/source/model"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type personProfileFs struct {
	coll           string
	fieldAccountId string
	fieldSessionId string
	fieldToken     string
	fieldFullName  string
	// fieldStudentId string
	fieldEmail          string
	fieldPhoneNumber    string
	fieldHashedPassword string
	fieldBirthday       string
	fieldModifiedAt     string
	fieldCreatedAt      string
}

var PersonProfile = &personProfileFs{
	coll:           "person_profile",
	fieldAccountId: "account_id",
	fieldSessionId: "session_id",
	fieldToken:     "token",
	fieldFullName:  "full_name",
	// fieldStudentId:      "student_id",
	fieldEmail:          "email",
	fieldPhoneNumber:    "phone_number",
	fieldHashedPassword: "hashed_password",
	fieldBirthday:       "birthday",
	fieldModifiedAt:     "modified_at",
	fieldCreatedAt:      "created_at",
}

type PersonProfileModel struct {
	AccountId string `json:"accountId" firestore:"account_id"`
	SessionId string `json:"sessionId" firestore:"session_id"`
	Token     string `json:"token" firestore:"token"`
	FullName  string `json:"fullname" firestore:"full_name"`
	// StudentId      string    `json:"studentId" firestore:"student_id"`
	Email          string    `json:"email" firestore:"email"`
	PhoneNumber    string    `json:"phoneNumber" firestore:"phone_number"`
	HashedPassword string    `json:"hashedPassword" firestore:"hashed_password"`
	Birthday       time.Time `json:"birthday" firestore:"birthday"`
	ModifiedAt     time.Time `json:"modifiedDate" firestore:"modified_at"`
	CreatedAt      time.Time `json:"createdDate" firestore:"created_at"`
}

func (me *personProfileFs) Add(ctx context.Context, accountId,
	sessionId, token string) (id string, err error) {
	var item = PersonProfileModel{
		AccountId: accountId,
		SessionId: sessionId,
		Token:     token,
	}
	return add(ctx, me.coll, item)
}
func (ins *personProfileFs) AddPersonProfile(ctx context.Context,
	personProfile *PersonProfileModel) error {
	_, err := add(ctx, ins.coll, personProfile)
	if err != nil {
		return err
	}
	return nil
}

func (ins *personProfileFs) CreateSignupProfile(ctx context.Context, account_id, session_id, email, phone_number, full_name, hashed_password string) (*PersonProfileModel, bool, error) {
	// make new data and insert db
	person_new := PersonProfileModel{
		AccountId:      account_id,
		SessionId:      session_id,
		Email:          email,
		PhoneNumber:    phone_number,
		FullName:       full_name,
		HashedPassword: hashed_password,
		ModifiedAt:     time.Now(),
		CreatedAt:      time.Now(),
	}
	if err := ins.AddPersonProfile(ctx, &person_new); err != nil {
		return nil, false, err
	}
	return &person_new, false, nil
}
func (ins *personProfileFs) CreateIfNotExist(ctx context.Context, account_id, session_id, token string) (*PersonProfileModel, bool, error) {
	id, inf, ok, err := ins.GetByAccountId(ctx, account_id)
	if err != nil {
		return nil, false, err
	}
	if ok {
		// update token
		if err := ins.SetToken(ctx, id, session_id, token); err != nil {
			return nil, false, err
		} else {
			inf.SessionId = session_id
			inf.Token = token
		}
		return inf, true, nil
	}
	// make new data and insert db
	person_new := PersonProfileModel{
		AccountId: account_id,
		SessionId: session_id,
		Token:     token,
	}
	if err := ins.AddPersonProfile(ctx, &person_new); err != nil {
		return nil, false, err
	}
	return &person_new, false, nil
}

func (ins *personProfileFs) CheckLogin(ctx context.Context, email, password string) (
	id, account_id, full_name, phone_number string, birthday time.Time, err error) {
	var (
		temp PersonProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldEmail, email)
	if err != nil {
		return "", "", "", "", time.Time{}, errors.New("email or password invalid")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(temp.HashedPassword), []byte(password)); err != nil {
		return "", "", "", "", time.Time{}, errors.New("email or password invalid")
	}

	return id, temp.AccountId, temp.FullName, temp.PhoneNumber, temp.Birthday, nil
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

func (ins *personProfileFs) GetByAccountId(ctx context.Context, account_id string) (
	id string, info *PersonProfileModel, ok bool, err error) {
	var (
		temp PersonProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldAccountId, account_id)
	if err == model.ErrDocNotFound {
		return "", nil, false, nil
	}
	if err != nil {
		return "", nil, false, err
	}
	return id, &temp, true, nil
}

func (ins *personProfileFs) GetByEmail(ctx context.Context, email string) (
	id string, info *PersonProfileModel, ok bool, err error) {
	var (
		temp PersonProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldEmail, email)
	if err == model.ErrDocNotFound {
		return "", nil, false, nil
	}
	if err != nil {
		return "", nil, false, err
	}
	return id, &temp, true, nil

}
func (ins *personProfileFs) GetByPhone(ctx context.Context, numberPhone string) (
	id string, info *PersonProfileModel, ok bool, err error) {
	var (
		temp PersonProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldPhoneNumber, numberPhone)
	if err == model.ErrDocNotFound {
		return "", nil, false, nil
	}
	if err != nil {
		return "", nil, false, err
	}
	return id, &temp, true, nil

}

func (ins *personProfileFs) SetToken(ctx context.Context, docId string, session_id, token string) error {
	var update = map[string]interface{}{
		ins.fieldModifiedAt: time.Now(),
		ins.fieldSessionId:  session_id,
		ins.fieldToken:      token,
	}
	return updateFields(ctx, docId, ins.coll, update)
}
