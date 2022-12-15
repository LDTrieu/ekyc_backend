package fsdb

import (
	"context"
	"ekyc-app/source/model"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type studentProfileFs struct {
	coll               string
	fieldStudentId     string
	fieldEmail         string
	fieldFullName      string
	fieldPhoneNumber   string
	fieldNationalId    string
	fieldBirthday      string
	fieldSex           string
	fieldAddress       string
	fieldAddressOrigin string
	fieldUnitId        string
	fieldImage         string
	filedImageEkyc     string
	fieldIsBlocked     string
	fieldModifiedBy    string
	fieldModifiedAt    string
	fieldCreatedBy     string
	fieldCreatedAt     string
}

var StudentProfile = &studentProfileFs{
	coll:               "student_profile",
	fieldStudentId:     "student_id",
	fieldEmail:         "email",
	fieldFullName:      "full_name",
	fieldPhoneNumber:   "phone_number",
	fieldNationalId:    "national_id",
	fieldBirthday:      "birthday",
	fieldSex:           "sex",
	fieldAddress:       "address",
	fieldAddressOrigin: "address_origin",
	fieldUnitId:        "unit_id",
	fieldImage:         "image",
	filedImageEkyc:     "image_ekyc",
	fieldIsBlocked:     "is_blocked",
	fieldModifiedBy:    "modified_by",
	fieldModifiedAt:    "modified_at",
	fieldCreatedBy:     "created_by",
	fieldCreatedAt:     "created_at",
}

type StudentProfileModel struct {
	StudentId     string    `json:"studentId" firestore:"student_id"`
	Email         string    `json:"email" firestore:"email"`
	FullName      string    `json:"fullName" firestore:"full_name"`
	PhoneNumber   string    `json:"phoneNumber" firestore:"phone_number"`
	UnitId        string    `json:"unitId" firestore:"unit_id"`
	NationalId    string    `json:"nationalId" firestore:"national_id"`
	Birthday      time.Time `json:"birthday" firestore:"birthday"`
	Sex           string    `json:"sex"  firestore:"sex"`
	Address       string    `json:"address" firestore:"address"`
	AddressOrigin string    `json:"addressOrigin" firestore:"address_origin"`

	Image      string    `json:"image" firestore:"image"`
	ImageEkyc  string    `json:"imageEkyc" firestore:"image_ekyc"`
	IsBlocked  bool      `json:"isBlocked" firestore:"is_blocked"`
	ModifiedBy string    `json:"modifiedBy" firestore:"modified_by"`
	ModifiedAt time.Time `json:"modifiedAt" firestore:"modified_at"`
	CreatedBy  string    `json:"createdBy" firestore:"created_by"`
	CreatedAt  time.Time `json:"createdAt" firestore:"created_at"`
}

func (ins *studentProfileFs) AddStudentProfile(ctx context.Context,
	studentProfile *StudentProfileModel) error {
	_, err := add(ctx, ins.coll, studentProfile)
	if err != nil {
		return err
	}
	return nil
}

func (ins *studentProfileFs) CreateStudentProfile(
	ctx context.Context, student_id, email,
	full_name, phone_number, national_id string, birthday time.Time, sex, address, address_origin,
	unitId, image, image_ekyc, create_by string) error {

	temp := StudentProfileModel{
		StudentId:     student_id,
		Email:         email,
		FullName:      full_name,
		PhoneNumber:   phone_number,
		NationalId:    national_id,
		Birthday:      birthday,
		Sex:           sex,
		Address:       address,
		AddressOrigin: address_origin,
		UnitId:        unitId,
		Image:         image,
		ImageEkyc:     image_ekyc,
		CreatedBy:     create_by,
		CreatedAt:     time.Now(),
	}
	_, err := add(ctx, ins.coll, &temp)
	if err != nil {
		return err
	}

	return nil

}

func (ins *studentProfileFs) GetAll(ctx context.Context) (
	[]*StudentProfileModel, error) {
	var (
		list = make([]*StudentProfileModel, 0)
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
			var temp StudentProfileModel
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

func (ins *studentProfileFs) ValidateEmail(ctx context.Context,
	email string) (already_exist bool, err error) {
	var (
		count int
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.
			Where(ins.fieldEmail, Equal, email).
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

func (ins *studentProfileFs) ValidatePhoneNumber(ctx context.Context,
	phone_number string) (already_exist bool, err error) {
	var (
		count int
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.
			Where(ins.fieldPhoneNumber, Equal, phone_number).
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

func (ins *studentProfileFs) ValidateNationalId(ctx context.Context,
	national_id string) (already_exist bool, err error) {
	var (
		count int
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.
			Where(ins.fieldNationalId, Equal, national_id).
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

func (ins *studentProfileFs) ValidateStudentId(ctx context.Context,
	student_id string) (already_exist bool, err error) {
	var (
		count int
	)
	if err := run(ctx, ins.coll, func(collectionRef *firestore.CollectionRef) error {
		dIter := collectionRef.
			Where(ins.fieldStudentId, Equal, student_id).
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

func (ins *studentProfileFs) GetModelByStudentId(
	ctx context.Context, student_id string) (
	id string, info *StudentProfileModel, ok bool, err error) {
	var (
		temp StudentProfileModel
	)
	id, err = getOneEqual(ctx, &temp, ins.coll, ins.fieldStudentId, student_id)
	if err == model.ErrDocNotFound {
		return "", nil, false, nil
	}
	if err != nil {
		return "", nil, false, err
	}
	return id, &temp, true, nil
}

func (ins *studentProfileFs) GetByStudentId(
	ctx context.Context, student_id string) (
	email, full_name, phone_number, national_id string,
	birthday time.Time, sex, address, address_origin,
	unit_id, image, image_ekyc, modified_by, created_by string, modified_at, created_at time.Time, ok bool, err error) {
	_, info, _, err := ins.GetModelByStudentId(ctx, student_id)
	if err == model.ErrDocNotFound {
		return "", "", "", "", time.Time{}, "", "", "", "", "", "", "", "", time.Time{}, time.Time{}, false, nil
	}
	if err != nil {
		return "", "", "", "", time.Time{}, "", "", "", "", "", "", "", "", time.Time{}, time.Time{}, false, err
	}

	return info.Email, info.FullName, info.PhoneNumber, info.NationalId, info.Birthday, info.Sex, info.Address, info.AddressOrigin,
		info.UnitId, info.Image, info.ImageEkyc, info.ModifiedBy, info.CreatedBy, info.ModifiedAt, info.CreatedAt, true, nil

}
