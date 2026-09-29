package models

type Res struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type UserAuthRes struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Error   string `json:"error"`
}

type UserRes struct {
	Success bool   `json:"success"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Error   string `json:"error"`
}

type PhonesRes struct {
	Success bool          `json:"success"`
	Data    []UserPhoneDb `json:"data"`
	Error   string        `json:"error"`
}
