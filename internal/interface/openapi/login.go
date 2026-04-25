package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

type LoginResponseBody struct {
	User        LoginUserResponseData `json:"user"`
	LoginStatus bool                  `json:"login_status"`
	Token       string                `json:"token"`
}

type LoginUserResponseData struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Gender Gender  `json:"gender"`
	Image  *string `json:"image"`
	Admin  bool    `json:"admin"`
}

func NewLoginResponse(user *model.User, token string) (LoginResponseBody, error) {
	gender, err := NewGender(user.Gender)
	if err != nil {
		return LoginResponseBody{}, err
	}
	return LoginResponseBody{
		User: LoginUserResponseData{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Gender: gender,
			Image:  user.Image,
			Admin:  user.Admin,
		},
		LoginStatus: true,
		Token:       token,
	}, nil
}
