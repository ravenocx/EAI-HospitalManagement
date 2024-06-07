package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	ID                  string     `db:"id" json:"id" validate:"required,uuid"`
	Nip                 string     `db:"nip" json:"nip" validate:"required"`
	Name                string     `db:"name" json:"name" validate:"required,min=5,max=50"`
	Role                string     `db:"role" json:"role" validate:"required"`
	Password            string     `db:"password" json:"password,omitempty" validate:"required,min=5,max=33"`
	IdentityCardScanImg string     `db:"identity_card_scan_img" json:"identityCardScanImg" validate:"required,img_url"`
	Access              bool       `db:"access" json:"access"`
	RefreshToken        string     `db:"refresh_token"`
	CreatedAt           time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt           *time.Time `db:"updated_at" json:"-"`
}

type NurseRegistrationPayload struct {
	Nip                       int64                 `json:"nip" form:"nip" validate:"required,nip_nurse"`
	Name                      string                `json:"name,omitempty" form:"name" validate:"required,min=5,max=50"`
	IdentityCardScanImg       *multipart.FileHeader `json:"identityCardScanImg" form:"identityCardScanImg" validate:"required,img_file"`
	IdentityCardScanImgString string                `db:"identity_card_scan_img"`
}

type NurseUpdatePayload struct {
	Nip  int64  `json:"nip" form:"nip" validate:"required,nip_nurse"`
	Name string `json:"name,omitempty" form:"name" validate:"required,min=5,max=50"`
}

type NurseAccessPayload struct {
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}

type GetUserQueries struct {
	UserId    string `db:"id" json:"userId" query:"userId" validate:"uuid"`
	Limit     int    `json:"limit" query:"limit"`
	Offset    int    `json:"offset" query:"offset"`
	Name      string `db:"name" json:"name" query:"name"`
	Nip       string `db:"nip" json:"nip" query:"nip"`
	Role      string `db:"role" json:"role" query:"role"`
	CreatedAt string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type GetUserResponse struct {
	UserId    string `db:"id" json:"userId" query:"userId"`
	Nip       int64  `db:"nip" json:"nip" query:"nip"`
	Name      string `db:"name" json:"name" query:"name"`
	CreatedAt string `db:"created_at" json:"createdAt" query:"createdAt"`
}

type Credential struct {
	Nip      string `json:"nip" form:"nip" validate:"required,nip_admin"`
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}
