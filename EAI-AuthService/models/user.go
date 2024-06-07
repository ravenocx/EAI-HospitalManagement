package models

import "time"

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

type AdminRegistrationPayload struct {
	Nip      int64  `json:"nip" form:"nip" validate:"required,nip_admin"`
	Name     string `json:"name,omitempty" form:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}

type AdminCredential struct {
	Nip      int64  `json:"nip" form:"nip" validate:"required,nip_admin"`
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}

type NurseCredential struct {
	Nip      int64  `json:"nip" form:"nip" validate:"required,nip_nurse"`
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}

type Credential struct {
	Nip      string  `json:"nip" form:"nip" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=5,max=33"`
}
