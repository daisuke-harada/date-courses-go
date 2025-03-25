package handler

import "github.com/daisuke-harada/date-courses-go/internal/api"

type handler struct{}

func NewHandler() api.ServerInterface {
	return &handler{}
}
