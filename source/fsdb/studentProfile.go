package fsdb

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type studentProfileFs struct {
	coll             string
	fieldStudentId   string
	fieldEmail       string
	fieldFullName    string
	fieldPhoneNumber string
	fieldNationalId  string
	fieldBirthday    string
	fieldAddress     string
	fieldUnitId      string
	fieldImage       string
	filedImageEkyc   string
	fieldIsBlocked   string
	fieldModifiedBy  string
	fieldModifiedAt  string
	CreatedBy        string
	fieldCreatedAt   string
}

var StudentProfile = &studentProfileFs{
	coll:             "student_profile",
	fieldStudentId:   "student_id",
	fieldEmail:       "email",
	fieldFullName:    "full_name",
	fieldPhoneNumber: "phone_number",
	fieldNationalId:  "national_id",
	fieldBirthday:    "birthday",
	fieldAddress:     "address",
	fieldUnitId:      "unit_id",
	fieldImage:       "image",
	filedImageEkyc:   "image_ekyc",
	fieldIsBlocked:   "is_blocked",
	fieldModifiedBy:  "modified_by",
	fieldModifiedAt:  "modified_at",
	CreatedBy:        "created_by",
	fieldCreatedAt:   "created_at",
}

type StudentProfileModel struct {
	StudentId   string    `json:"studentId" firestore:"student_id"`
	Email       string    `json:"email" firestore:"email"`
	FullName    string    `json:"fullName" firestore:"full_name"`
	PhoneNumber string    `json:"phoneNumber" firestore:"phone_number"`
	NationalId  string    `json:"nationalId" firestore:"national_id"`
	Birthday    time.Time `json:"birthday" firestore:"birthday"`
	Address     string    `json:"address" firestore:"address"`
	UnitId      string    `json:"unitId" firestore:"unit_id"`
	Image       string    `json:"image" firestore:"image"`
	ImageEkyc   string    `json:"imageEkyc" firestore:"image_ekyc"`
	IsBlocked   bool      `json:"isBlocked" firestore:"is_blocked"`
	ModifiedBy  string    `json:"modifiedBy" firestore:"modified_by"`
	ModifiedAt  time.Time `json:"modifiedAt" firestore:"modified_at"`
	CreatedBy   string    `json:"createdBy" firestore:"created_by"`
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
}

func (ins *studentProfileFs) AddStudentProfile(ctx context.Context,
	studentProfile *StudentProfileModel) error {
	_, err := add(ctx, ins.coll, studentProfile)
	if err != nil {
		return err
	}
	return nil
}

// func (ins *studentProfileFs) CreateStudentProfile(ctx context.Context,
// 	student_id, email, full_name, phone_number, national_id, birthday, address, unitId,
// 	image, image_ekyc string) (*StudentProfileModel, bool, error) {

// }

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
