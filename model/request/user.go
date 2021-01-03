package request

type User struct {
	ID         int `gorm:"primarykey"`
	Name       string
	Password   string
	Token      string
	ExpireTime int64
	FailNum    int64
	FailTime   int64
}

// User LoginStruct
type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordStruct
type ChangePasswordStruct struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
