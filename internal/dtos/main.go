package dtos

import "net/http"

type DtoInterface interface {
	Describe() string
}

type UrlParamSetterInterface interface {
	SetUrlParams(req *http.Request)
}
