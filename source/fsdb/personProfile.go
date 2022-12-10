package fsdb

import (
	"context"
	"ekyc-app/source/model"
	"time"
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
	fieldCreatedDate    string
	fieldModifiedDate   string
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
	fieldCreatedDate:    "created_date",
	fieldModifiedDate:   "modified_date",
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
	CreatedDate    time.Time `json:"createdDate" firestore:"created_date"`
	ModifiedDate   time.Time `json:"modifiedDate" firestore:"modified_date"`
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

// func (ins *personProfileFs) CreateProfileSignUpBasic(ctx context.Context,
// 	personProfile *PersonProfileModel) error {
// 	// check email exist
// 	ins.GetByAccountId()

//		return nil
//	}
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
		ins.fieldSessionId: session_id,
		ins.fieldToken:     token,
	}
	return updateFields(ctx, docId, ins.coll, update)
}
