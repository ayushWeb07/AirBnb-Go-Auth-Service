package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/go-playground/validator/v10"
)

// HTTP middleware to decode and validate JSON request body
func DecodeAndValidateRequestBody[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
		userPayload := new(T)

		// read the data from the request body
		decodeErr := json.NewDecoder(req.Body).Decode(userPayload)

		if decodeErr != nil {
			utils.WriteJsonResponse(http.StatusBadRequest, resWriter, map[string]any{
				"success": false,
				"message": "Invalid json body has been provided",
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
func DecodeAndValidateParams[T any](extractor func(req *http.Request) (*T, *utils.AppError)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
			payload, err := extractor(req)

			if err != nil {
				utils.WriteJsonResponse(err.StatusCode, resWriter, map[string]any{
					"success": err.Success,
					"message": "Invalid req params has been provided",
					"error":   err.Error(),
				})

				return
			}

			// validate the request params
			validate := validator.New(validator.WithRequiredStructEnabled())
			validateErr := validate.Struct(payload)

			if validateErr != nil {
				utils.WriteJsonResponse(http.StatusBadRequest, resWriter, map[string]any{
					"success": false,
					"message": "Invalid req params has been provided",
					"error":   validateErr.Error(),
				})

				return
			}

			ctx := context.WithValue(req.Context(), "params", payload)
			next.ServeHTTP(resWriter, req.WithContext(ctx))
		})
	}
}
