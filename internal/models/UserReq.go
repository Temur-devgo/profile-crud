package models

type UserRegReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

type UserAuthReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserDb struct {
	Name     string `json:"name,omitempty" db:"name"`
	ID       int    `json:"id,omitempty" db:"id"`
	Age      int    `json:"age,omitempty" db:"age"`
	PassHash string `json:"password,omitempty" db:"password" `
}

type UserPhoneReq struct {
	Phone       string `json:"phone"`
	Description string `json:"description"`
	IsFax       bool   `json:"is_fax"`
}

type UserPhoneDb struct {
	ID          int    `json:"id" db:"id"`
	Phone       string `json:"phone" db:"phone"`
	Description string `json:"description" db:"description"`
	UserID      int    `json:"user_id" db:"user_id"`
	IsFax       bool   `json:"is_fax" db:"is_fax"`
}

type UserPhoneUpdateReq struct {
	ID          int    `json:"id"`
	Phone       string `json:"phone,omitempty"`
	Description string `json:"description,omitempty"`
	IsFax       bool   `json:"is_fax,omitempty"`
}
