package dto

type ValidateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (v *ValidateUser) GetLogin() string {
	return v.Email
}

func (v *ValidateUser) GetPassword() string {
	return v.Password
}
