package request

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
