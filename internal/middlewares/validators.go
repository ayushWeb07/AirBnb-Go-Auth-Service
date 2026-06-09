package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/go-playground/validator/v10"
)

// HTTP middleware to decode and validate JSON request body
func DecodeAndValidateRequestBody[T dtos.DtoInterface](next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		userPayload := new(T)

		// read the data from the request body
		decodeErr := json.NewDecoder(req.Body).Decode(userPayload)

		if decodeErr != nil {
			utils.WriteJsonResponse(http.StatusBadRequest, resWriter, map[string]any{
				"success": false,
				"message": "Failed to decode the json body",
				"error":   decodeErr.Error(),
			})

			return
		}

		// validate the request body
		validate := validator.New(validator.WithRequiredStructEnabled())
		validateErr := validate.Struct(userPayload)

		if validateErr != nil {
			utils.WriteJsonResponse(http.StatusBadRequest, resWriter, map[string]any{
				"success": false,
				"message": "Invalid json body has been provided",
				"error":   validateErr.Error(),
			})

			return
		}

		ctx := context.WithValue(req.Context(), "payload", userPayload)
		next.ServeHTTP(resWriter, req.WithContext(ctx))
	})
}

// HTTP middleware to decode and validate request params
func DecodeAndValidateParams[T dtos.UrlParamSetterInterface](next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		userPayload := new(T)

		// assign req url params
		*userPayload = (*userPayload).SetUrlParams(req).(T)

		// validate the request params
		validate := validator.New(validator.WithRequiredStructEnabled())
		validateErr := validate.Struct(userPayload)

		if validateErr != nil {
			utils.WriteJsonResponse(http.StatusBadRequest, resWriter, map[string]any{
				"success": false,
				"message": "Invalid json body has been provided",
				"error":   validateErr.Error(),
			})

			return
		}

		ctx := context.WithValue(req.Context(), "payload", userPayload)
		next.ServeHTTP(resWriter, req.WithContext(ctx))
	})
}
