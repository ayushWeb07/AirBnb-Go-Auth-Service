package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(serverConfig *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
			// read the req header
			authHeaderToken := req.Header.Get("Authorization")

			if authHeaderToken == "" {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "No token has been provided",
				})

				return
			}

			// check if it has the Bearer: prefix
			if ok := strings.HasPrefix(authHeaderToken, "Bearer "); !ok {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Invalid token has been provided",
				})

				return
			}

			// trim token and verify
			tokenString := strings.TrimPrefix(authHeaderToken, "Bearer ")

			if tokenString == "" {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Invalid token has been provided",
				})

				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				// invalid signing method had been used for token generating
				if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, fmt.Errorf("Invalid signing method had been used")
				}

				// else return the jwt secret key
				return []byte(serverConfig.AccessTokenSecretKey), nil
			})

			// check if there's an error or the token is invalid
			if err != nil {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Invalid token has been provided: " + err.Error(),
				})

				return
			}

			if !token.Valid {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Invalid or expired token has been provided",
				})

				return
			}

			// parse token to decode the payload
			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Failed to decode payload from an invalid token",
				})

				return
			}

			// access the payload
			userId := claims["id"].(float64)
			expiryTime := claims["exp"].(float64)

			// check if token has expired
			if time.Now().Unix() > int64(expiryTime) {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "Authentication failed",
					"error":   "Token has expired. Please login again",
				})

				return
			}

			// create a new context including the user details
			userParams := &dtos.GetUserByIdParams{
				ID: int(userId),
			}

			ctx := context.WithValue(req.Context(), "params", userParams)
			next.ServeHTTP(resWriter, req.WithContext(ctx))
		})
	}
}
