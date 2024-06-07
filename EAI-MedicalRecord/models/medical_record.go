package models

import "time"

type MedicalRecord struct {
	IdentityNumber  string    `db:"identity_number" json:"identity_number" validate:"required,identity_number"`
	Symptoms        string    `db:"symptoms" json:"symptoms" validate:"required,min=1,max=2000"`
	Medications     string    `db:"medications" json:"medications" validate:"required,min=1,max=2000"`
	CreatedByNip    string    `db:"created_by_nip" json:"createdByNip"`
	CreatedByName   string    `db:"created_by_name" json:"createdByName"`
	CreatedByUserId string    `db:"created_by_user_id" json:"createdByUserId"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

type RecordRegistrationPayload struct {
	IdentityNumber int64  `db:"identity_number" json:"identityNumber" form:"identityNumber" validate:"required,identity_number"`
	Symptoms       string `db:"symptoms" json:"symptoms" form:"symptoms" validate:"required,min=1,max=2000"`
	Medications    string `db:"medications" json:"medications" form:"medications" validate:"required,min=1,max=2000"`
}

type GetRecordQueries struct {
	IdentityNumber  *int64 `db:"identity_number" json:"identityNumber" query:"identityNumber" validate:"required,identity_number"`
	Limit           int    `json:"limit" query:"limit"`
	Offset          int    `json:"offset" query:"offset"`
	CreatedByNip    string `db:"created_by_nip" json:"createdByNip" query:"createdBy.nip"`
	CreatedByUserId string `db:"created_by_user_id" json:"createdByUserId" query:"createdBy.userId"`
	CreatedAt       string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetRecordResponse struct {
	IdentityDetail PatientDetail   `json:"identityDetail"`
	Symptoms       string          `json:"symptoms"`
	Medications    string          `json:"medications"`
	CreatedBy      CreatedByDetail `json:"createdBy"`
	CreatedAt      string          `json:"createdAt"`
}

type CreatedByDetail struct {
	Nip    string `json:"nip"`
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type PatientDetail struct {
	IdentityNumber      int64  `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	BirthDate           string `json:"birthDate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type NurseResponse struct {
	Message string  `json:"message"`
	Data    []Nurse `json:"data"`
}

type Nurse struct {
	UserID    string `json:"userId"`
	NIP       int64  `json:"nip"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type PatientResponse struct {
	Message string  `json:"message"`
	Data    []Patient `json:"data"`
}

type Patient struct {
	IdentityNumber int64  `db:"identity_number" json:"identityNumber"`
	PhoneNumber    string `db:"phone_number" json:"phoneNumber"`
	Name           string `db:"name" json:"name"`
	BirthDate      string `db:"birth_date" json:"birthDate"`
	Gender         string `db:"gender" json:"gender"`
	CreatedAt      string `db:"created_at" json:"createdAt"`
}
