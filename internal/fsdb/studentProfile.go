package fsdb

import (
	"context"
	"ekyc-app/internal/model"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type studentProfileFs struct {
	coll               string
	fieldStudentId     string
	fieldPersonId      string
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

	fieldFullNameEkyc      string
	fieldFaceImageURL      string
	fieldNationalIdCardURL string
	fieldFaceVideoURL      string
	fieldFaceThumbnailURL  string
	fieldAddressEkyc       string
	fieldPlaceOfOrigin     string
	fieldNationality       string
	fieldDateOfBirth       string
	fieldDateOfExpiry      string
}

var StudentProfile = &studentProfileFs{
	coll:             "student_profile",
	fieldStudentId:   "student_id",
	fieldEmail:       "email",
	fieldFullName:    "full_name",
	fieldPhoneNumber: "phone_number",

	fieldBirthday: "birthday",
	fieldSex:      "sex",
	fieldAddress:  "address",

	fieldUnitId:    "unit_id",
	fieldImage:     "image",
	filedImageEkyc: "image_ekyc",

	fieldPersonId:          "person_id",
	fieldNationalId:        "national_id",
	fieldFullNameEkyc:      "full_name_ekyc",
	fieldFaceImageURL:      "face_image_url",
	fieldNationalIdCardURL: "national_id_card_url",
	fieldFaceVideoURL:      "face_video_url",
	fieldFaceThumbnailURL:  "face_thumbnail_url",
	fieldModifiedBy:        "modified_by",
	fieldModifiedAt:        "modified_at",
	fieldCreatedBy:         "created_by",
	fieldCreatedAt:         "created_at",
	fieldAddressEkyc:       "address_ekyc",
	fieldPlaceOfOrigin:     "place_of_origin",
	fieldNationality:       "nationality",
	fieldDateOfBirth:       "date_of_birth",
	fieldDateOfExpiry:      "date_of_expiry",

	fieldIsBlocked: "is_blocked",
}

type StudentProfileModel struct {
	StudentId   string    `json:"studentId" firestore:"student_id"`
	Email       string    `json:"email" firestore:"email"`
	FirstName   string    `json:"firstName" firestore:"first_name"`
	LastName    string    `json:"lastName" firestore:"last_name"`
	FullName    string    `json:"fullName" firestore:"full_name"`
	PhoneNumber string    `json:"phoneNumber" firestore:"phone_number"`
	UnitId      string    `json:"unitId" firestore:"unit_id"`
	Birthday    time.Time `json:"birthday" firestore:"birthday"`
	Sex         string    `json:"sex"  firestore:"sex"`
	Address     string    `json:"address" firestore:"address"`

	PersonId          string    `json:"personId" firestore:"person_id"`
	FullNameEkyc      string    `json:"fullNameEkyc" firestore:"full_name_ekyc"`
	NationalId        string    `json:"nationalId" firestore:"national_id"`
	FaceImageURL      string    `json:"faceImageURL" firestore:"face_image_url"`
	NationalIdCardURL string    `json:"nationalIdCardURL" firestore:"national_id_card_url"`
	FaceVideoURL      string    `json:"faceVideoURL" firestore:"face_video_url"`
	FaceThumbnailURL  string    `json:"faceThumbnailURL" firestore:"face_thumbnail_url"`
	ModifiedBy        string    `json:"modifiedBy" firestore:"modified_by"`
	ModifiedAt        time.Time `json:"modifiedAt" firestore:"modified_at"`
	CreatedBy         string    `json:"createdBy" firestore:"created_by"`
	CreatedAt         time.Time `json:"createdAt" firestore:"created_at"`
	AddressEkyc       string    `json:"addressEkyc" firestore:"address_ekyc"`
	PlaceOfOrigin     string    `json:"placeOfOrigin" firestore:"place_of_origin"`
	Nationality       string    `json:"nationality" firestore:"nationality"`
	DateOfBirth       time.Time `json:"dateOfBirth" firestore:"date_of_birth"`
	DateOfExpiry      time.Time `json:"dateOfExpiry" firestore:"date_of_expiry"`

	IsBlocked bool `json:"isBlocked" firestore:"is_blocked"`
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
	first_name, last_name, phone_number, national_id string, birthday time.Time, sex, address, address_origin,
	unitId, image, image_ekyc, create_by string) error {

	temp := StudentProfileModel{
		StudentId:         student_id,
		Email:             email,
		FirstName:         first_name,
		LastName:          last_name,
		FullName:          fmt.Sprintf("%s%s%s", first_name, " ", last_name),
		PhoneNumber:       phone_number,
		NationalId:        national_id,
		Birthday:          birthday,
		Sex:               sex,
		Address:           address,
		PlaceOfOrigin:     address_origin,
		UnitId:            unitId,
		FullNameEkyc:      image,
		NationalIdCardURL: image_ekyc,
		CreatedBy:         create_by,
		CreatedAt:         time.Now(),
		ModifiedBy:        create_by,
		ModifiedAt:        time.Now(),
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
	unit_id, image, image_ekyc, modified_by, created_by string,
	modified_at, created_at time.Time, ok bool, err error) {
	_, info, _, err := ins.GetModelByStudentId(ctx, student_id)
	if err == model.ErrDocNotFound {
		return "", "", "", "", time.Time{}, "", "", "", "", "", "", "", "", time.Time{}, time.Time{}, false, nil
	}
	if err != nil {
		return "", "", "", "", time.Time{}, "", "", "", "", "", "", "", "", time.Time{}, time.Time{}, false, err
	}
	return info.Email, info.FullName, info.PhoneNumber, info.NationalId, info.Birthday, info.Sex, info.Address, info.PlaceOfOrigin,
		info.UnitId, info.FaceImageURL, info.NationalIdCardURL, info.ModifiedBy, info.CreatedBy, info.ModifiedAt, info.CreatedAt, true, nil
}

func (ins *studentProfileFs) GetNationIdByStudentId(
	ctx context.Context, student_id string) (
	doc_id, full_name, national_id string, ok bool, err error) {
	doc_id, info, _, err := ins.GetModelByStudentId(ctx, student_id)
	if err == model.ErrDocNotFound {
		return "", "", "", false, nil
	}
	if err != nil {
		return "", "", "", false, err
	}
	return doc_id, info.FullName, info.NationalId, true, nil
}

func (ins *studentProfileFs) SetFaceImageURL(
	ctx context.Context, docId, photoURL string) error {
	var update = map[string]interface{}{
		ins.fieldImage: photoURL,
	}

	return updateFields(ctx, docId, ins.coll, update)
}

func (ins *studentProfileFs) SetNationalIdImageURL(
	ctx context.Context, docId, photoURL string) error {
	var update = map[string]interface{}{
		ins.fieldNationalIdCardURL: photoURL,
	}

	return updateFields(ctx, docId, ins.coll, update)
}

func (ins *studentProfileFs) SetFaceVideoURL(
	ctx context.Context, docId, videoURL, thumbnailURL string) error {
	var update = map[string]interface{}{
		ins.fieldFaceVideoURL:     videoURL,
		ins.fieldFaceThumbnailURL: thumbnailURL,
	}

	return updateFields(ctx, docId, ins.coll, update)
}

func (ins *studentProfileFs) SetIsBlocked(
	ctx context.Context, student_id string, is_blocked bool) error {
	doc_id, _, _, err := ins.GetModelByStudentId(ctx, student_id)
	if err == model.ErrDocNotFound {
		return err
	}
	if err != nil {
		return err
	}
	var update = map[string]interface{}{
		ins.fieldIsBlocked: is_blocked,
	}

	return updateFields(ctx, doc_id, ins.coll, update)
}

func (ins *studentProfileFs) SetEkyc(
	ctx context.Context, student_id, national_id, person_id, full_name_ekyc,
	gender, face_image_url, national_id_card_url, address_ekyc,
	place_of_origin, nationality, modified_by string,
	date_of_birth, date_of_expiry time.Time) error {
	doc_id, _, _, err := ins.GetModelByStudentId(ctx, student_id)
	if err == model.ErrDocNotFound {
		return err
	}
	if err != nil {
		return err
	}
	var update = map[string]interface{}{
		ins.fieldNationalId:        national_id,
		ins.fieldPersonId:          person_id,
		ins.fieldFullNameEkyc:      full_name_ekyc,
		ins.fieldSex:               gender,
		ins.fieldFaceImageURL:      face_image_url,
		ins.fieldNationalIdCardURL: national_id_card_url,

		ins.fieldModifiedBy:    modified_by,
		ins.fieldModifiedAt:    time.Now(),
		ins.fieldAddressEkyc:   address_ekyc,
		ins.fieldPlaceOfOrigin: place_of_origin,
		ins.fieldNationality:   nationality,
		ins.fieldDateOfBirth:   date_of_birth,
		ins.fieldDateOfExpiry:  date_of_expiry,

		ins.fieldIsBlocked: false,
	}

	return updateFields(ctx, doc_id, ins.coll, update)
}
