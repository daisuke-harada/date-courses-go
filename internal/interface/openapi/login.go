package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

// LoginResponseBody は POST /api/v1/login のレスポンス型です。
// Rails の SessionsSerializer / UserSerializer に対応します。
type LoginResponseBody struct {
	User        LoginUserResponseData `json:"user"`
	LoginStatus bool                  `json:"login_status"`
}

// LoginUserResponseData はログインレスポンスに含まれるユーザーデータです。
type LoginUserResponseData struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Gender string  `json:"gender"`
	Image  *string `json:"image"`
	Admin  bool    `json:"admin"`
}

// NewLoginResponse は LoginOutput からレスポンス型を生成します。
func NewLoginResponse(user *model.User) LoginResponseBody {
	return LoginResponseBody{
		User: LoginUserResponseData{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Gender: user.Gender,
			Image:  user.Image,
			Admin:  user.Admin,
		},
		LoginStatus: true,
	}
}
