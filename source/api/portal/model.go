package portal

import (
	"ekyc-app/source/auth"
	"ekyc-app/source/fsdb"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	qrActionLoginWebPortal = "login-web-portal"
	// qrActionLinkDevice     = "link-business-device"
)

type traceField struct {
	RequestId string `json:"reqId"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* */
type rendQRLoginRequest struct {
	traceField
}

type rendQRLoginResponse struct {
	traceField
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Payload rend_qr_code_login_resp `json:"payload"`
}
type rend_qr_code_login_resp struct {
	Path      string `json:"path"`
	JWT       string `json:"jwt"`
	ExpiresIn int64  `json:"expiresIn"`
	IssuedAt  int64  `json:"issuedAt"`
}

/* */
type loginBasicRequest struct {
	traceField
	Credentials
}

type loginBasicResponse struct {
	traceField
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Payload login_basic_data `json:"payload"`
}

type login_basic_data struct {
	AccountId   string    `json:"accountId"`
	FullName    string    `json:"fullName"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Birthday    time.Time `json:"birthday"`
	Avt         string    `json:"avt"`
	Banner      string    `json:"banner"`
	Token       string    `json:"token"`
}

func (ins *Credentials) validate() error {
	if len(ins.Email) < 1 {
		return errors.New("field email invalid")
	}
	if len(ins.Password) < 1 {
		return errors.New("field password invalid")
	}

	return nil
}

/* */
type signupBasicRequest struct {
	traceField
	Credentials
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	// some fields
}

type signupBasicResponse struct {
	traceField
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Payload signup_basic_data `json:"payload"`
}

type signup_basic_data struct {
	AccountId string `json:"accountId"`
	Email     string `json:"email"`
}

func (ins *signupBasicRequest) validate() error {
	if err := ins.Credentials.validate(); err != nil {
		return err
	}
	if len(ins.PhoneNumber) < 1 {
		return errors.New("field phone number invalid")
	}
	if len(ins.FullName) < 1 {
		return errors.New("field full name invalid")
	}
	return nil

}

/* */
type filterListUserRequest struct {
	traceField
	Permit *auth.DataJWT
	//	Payload list_user_req `json:"payload"`
}

type filterListUserResponse struct {
	traceField
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Payload list_user_resp `json:"payload"`
}

type list_user_resp struct {
	TotalUser int         `json:"totalUser"`
	ListUser  []user_data `json:"listUser"`
}

type user_data struct {
	FullName    string `json:"fullName"`
	Image       string `json:"image"`
	PhoneNumber string `json:"phoneNumber"`
	IsBlocked   bool   `json:"isBlocked"`
}

func withUserModel(um *fsdb.PersonProfileModel) user_data {
	return user_data{
		FullName:    um.FullName,
		PhoneNumber: um.PhoneNumber,
		// Image     : ,
		IsBlocked: um.IsBlocked,
	}
}

/* */
type filterListStudentRequest struct {
	traceField
	Permit *auth.DataJWT
	//	Payload list_user_req `json:"payload"`
}

type filterListStudentResponse struct {
	traceField
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Payload list_student_resp `json:"payload"`
}

type list_student_resp struct {
	Stash        int            `json:"stash"`
	TotalStudent int            `json:"totalStudent"`
	ListStudent  []student_data `json:"listStudent"`
}

type student_data struct {
	//student_detail
	FullName    string    `json:"fullName"`
	Image       string    `json:"image"`
	StudentId   string    `json:"studentId"`
	PhoneNumber string    `json:"phoneNumber"`
	UnitId      string    `json:"unitId" `
	Birthday    time.Time `json:"birthday"`
	Gender      string    `json:"gender"`

	//student_ekyc

	//student_censorship
	IsBlocked bool `json:"isBlocked"`
}

func withStudentModel(sm *fsdb.StudentProfileModel) student_data {
	return student_data{
		//student_detail
		FullName:    sm.FullName,
		StudentId:   sm.StudentId,
		Image:       "https://tuk-cdn.s3.amazonaws.com/assets/components/advance_tables/at_1.png",
		PhoneNumber: sm.PhoneNumber,
		Birthday:    sm.Birthday,
		UnitId:      sm.UnitId,
		Gender:      sm.Sex,
		//student_ekyc

		//student_censorship
		IsBlocked: sm.IsBlocked,
	}
}

/* */
type createStudentProfileRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload student_profile_data `json:"payload"`
}

type createStudentProfileResponse struct {
	traceField
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Payload list_student_resp `json:"payload"`
}

type student_profile_data struct {
	StudentId     string    `json:"studentId"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	PhoneNumber   string    `json:"phoneNumber"`
	UnitId        string    `json:"unitId"`
	NationalId    string    `json:"nationalId"`
	Birthday      time.Time `json:"birthday"`
	Sex           string    `json:"sex"`
	Address       string    `json:"address"`
	AddressOrigin string    `json:"addressOrigin"`
	Image         string    `json:"image" `
	ImageEkyc     string    `json:"imageEkyc"`
	IsBlocked     bool      `json:"isBlocked"`
}

func (ins *student_profile_data) validate() error {
	if len(ins.StudentId) < 1 {
		log.Println("ins.StudentId", ins.StudentId)
		return errors.New("field student_id is invalid")
	}
	if len(ins.Email) < 1 {
		return errors.New("field Email is invalid")
	}
	if len(ins.FullName) < 1 {
		return errors.New("field full_name is invalid")
	}
	if len(ins.UnitId) < 1 {
		return errors.New("field unit_id is invalid")
	}
	if len(ins.NationalId) < 1 {
		return errors.New("field national_id is invalid")
	}
	if len(ins.Address) < 1 {
		return errors.New("field address is invalid")
	}
	if len(ins.Image) < 1 {
		return errors.New("field image is invalid")
	}
	if len(ins.ImageEkyc) < 1 {
		return errors.New("field image_ekyc is invalid")
	}

	return nil
}

/* */

type studentDetailsRequest struct {
	traceField
	Permit    *auth.DataJWT
	StudentId string
	// Payload   student_profile_data `json:"payload"`
}

type studentDetailsResponse struct {
	traceField
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Payload student_details_data `json:"payload"`
}

type student_details_data struct {
	StudentId     string    `json:"studentId"`
	Email         string    `json:"email"`
	FullName      string    `json:"fullName"`
	PhoneNumber   string    `json:"phoneNumber"`
	UnitId        string    `json:"unitId"`
	NationalId    string    `json:"nationalId"`
	Birthday      time.Time `json:"birthday"`
	Sex           string    `json:"sex"`
	Address       string    `json:"address"`
	AddressOrigin string    `json:"addressOrigin"`
	Image         string    `json:"image" `
	ImageEkyc     string    `json:"imageEkyc"`
	IsBlocked     bool      `json:"isBlocked"`
	ModifiedBy    string    `json:"modifiedBy"`
	ModifiedAt    time.Time `json:"modifiedAt"`
	CreatedBy     string    `json:"createdBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

/* */
type uploadFaceImageRequest struct {
	traceField
	Permit  *auth.DataJWT
	Payload face_image_req `json:"payload"`
}

type uploadFaceImageResponse struct {
	traceField
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Payload face_image_resp `json:"payload"`
}

type face_image_req struct {
	StudentId string `json:"studentId"`
	FileName  string `json:"fileName"`
	File      []byte `json:"file"`
}

type face_image_resp struct {
	URL  string `json:"url"`
	Path string `json:"photoPath"`
}

func (ins *face_image_req) validate() error {
	if len(ins.StudentId) < 1 {
		return errors.New("student_id invalids")
	}
	if len(ins.FileName) < 3 {
		ins.FileName = fmt.Sprintf("%s_%s.bin", ins.StudentId, primitive.NewObjectID().Hex())
	}
	// if len(ins.File) < 128 {
	// 	return errors.New("file is too small")
	// }
	return nil
}
