package models

import (
	"mime/multipart"
	"time"
)

type Patient struct {
	IdentityNumber      string    `db:"identity_number" json:"id" validate:"required,identity_number"`
	PhoneNumber         string    `db:"phone_number" json:"phoneNumber" validate:"required,min=10,max=15,phone_number"`
	Name                string    `db:"name" json:"name" validate:"required,min=5,max=50"`
	BirthDate           string    `db:"birth_date" json:"birthDate" validate:"required,birth_date"`
	Gender              string    `db:"gender" json:"gender" validate:"required,oneof='male' 'female'"`
	IdentityCardScanImg string    `db:"identity_card_scan_img" json:"identityCardScanImg" validate:"required,img_url"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
}

type PatientRegistrationPayload struct {
	IdentityNumber            int64                 `db:"identity_number" json:"identityNumber" form:"identityNumber" validate:"required,identity_number"`
	PhoneNumber               string                `db:"phone_number" json:"phoneNumber" form:"phoneNumber" validate:"required,min=10,max=15,phone_number"`
	Name                      string                `db:"name" json:"name" form:"name" validate:"required,min=5,max=50"`
	BirthDate                 string                `db:"birth_date" json:"birthDate" form:"birthDate" validate:"required,birth_date"`
	Gender                    string                `db:"gender" json:"gender" form:"gender" validate:"required,oneof='male' 'female'"`
	IdentityCardScanImg       *multipart.FileHeader `json:"identityCardScanImg" form:"identityCardScanImg" validate:"required,img_file"`
	IdentityCardScanImgString string                `db:"identity_card_scan_img"`
}

type GetPatientQueries struct {
	IdentityNumber *int64 `db:"identity_number" json:"identityNumber" query:"identityNumber" validate:"identity_number"`
	Limit          int    `json:"limit" query:"limit"`
	Offset         int    `json:"offset" query:"offset"`
	Name           string `db:"name" json:"name" query:"name"`
	PhoneNumber    string `db:"phone_number" json:"phoneNumber" query:"phoneNumber"`
	CreatedAt      string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetPatientResponse struct {
	IdentityNumber int64  `db:"identity_number" json:"identityNumber"`
	PhoneNumber    string `db:"phone_number" json:"phoneNumber"`
	Name           string `db:"name" json:"name"`
	BirthDate      string `db:"birth_date" json:"birthDate"`
	Gender         string `db:"gender" json:"gender"`
	CreatedAt      string `db:"created_at" json:"createdAt"`
}
