package dto

type UserDto struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type AuthUserDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewAuthUserDto(login string, password string) AuthUserDto {
	return AuthUserDto{
		Login:    login,
		Password: password,
	}
}
