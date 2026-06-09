package dtos

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateUser struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (u CreateUser) Describe() string {
	return "Create User DTO ideally will be used in the request json body"
}

type GetUserByUsernameAndEmail struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
}

func (u GetUserByUsernameAndEmail) Describe() string {
	return "Get User By Username And Email DTO ideally will be used in the request json body"
}

type LoginUser struct {
	Username string `json:"username" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

func (u LoginUser) Describe() string {
	return "Login User DTO ideally will be used in the request json body"
}

type GetUserById struct {
	ID string `validate:"required,number"`
}

func (u GetUserById) SetUrlParams(req *http.Request) {
	u.ID = chi.URLParam(req, "id")
}

func (u GetUserById) Describe() string {
	return "Get User By Id DTO ideally will be used in the request params"
}

type DeleteUserById struct {
	ID string `validate:"required,number"`
}

func (u DeleteUserById) SetUrlParams(req *http.Request) {
	u.ID = chi.URLParam(req, "id")
}

func (u DeleteUserById) Describe() string {
	return "Delete User By Id DTO ideally will be used in the request params"
}
