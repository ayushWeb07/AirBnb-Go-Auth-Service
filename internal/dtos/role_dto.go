package dtos

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateRole struct {
	Name        string `json:"name" validate:"required,min=6,max=100"`
	Description string `json:"description" validate:"required,min=6,max=100"`
}

type GetUserByIdParams struct {
	ID int `validate:"required,number"`
}

type DeleteUserById struct {
	ID string `validate:"required,number"`
}

func (u DeleteUserById) SetUrlParams(req *http.Request) UrlParamSetterInterface {
	u.ID = chi.URLParam(req, "id")
	return u
}

func (u DeleteUserById) Describe() string {
	return "Delete User By Id DTO ideally will be used in the request params"
}
