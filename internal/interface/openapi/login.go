package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

type LoginUserResponseData struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Gender Gender  `json:"gender"`
	Image  *string `json:"image"`
	Admin  bool    `json:"admin"`
}

func NewLoginResponse(user *model.User, token string) (LoginResponseData, error) {
	return LoginResponseData{
		LoginStatus: true,
		Token:       token,
	}, nil
}
