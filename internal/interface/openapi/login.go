package openapi

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/daisuke-harada/date-courses-go/internal/domain/model"
)

func NewLoginResponse(user *model.User, token string) (LoginResponseData, error) {
	gender, err := NewGender(user.Gender)
	if err != nil {
		return LoginResponseData{}, err
	}
	return LoginResponseData{
		User: UserData{
			Id:     int(user.ID),
			Name:   user.Name,
			Email:  openapi_types.Email(user.Email),
			Gender: gender,
			Admin:  user.Admin,
			Image:  ImageData{Url: user.Image},
		},
		LoginStatus: true,
		Token:       token,
	}, nil
}
