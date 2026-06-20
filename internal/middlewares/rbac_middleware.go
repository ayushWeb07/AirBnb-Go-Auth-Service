package middlewares

import (
	"net/http"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
)

func RequireUserAllRoles(userRolesRepository *repositories.UserRolesRepository, roleNames []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resWriter http.ResponseWriter, req *http.Request) {
			userParams := req.Context().Value("params").(*dtos.GetUserByIdParams)

			// call the repository layer
			hasAllRoles := userRolesRepository.CheckUserHasAllRoles(&dtos.CheckUserHasAllRolesPayload{
				UserID:    userParams.ID,
				RoleNames: roleNames,
			})

			if !hasAllRoles {
				utils.WriteJsonResponse(http.StatusUnauthorized, resWriter, map[string]any{
					"success": false,
					"message": "RBAC authorization failed",
					"error":   "You do not have all the required roles",
				})

				return
			}

			next.ServeHTTP(resWriter, req)
		})
	}
}
