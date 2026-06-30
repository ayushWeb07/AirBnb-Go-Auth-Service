package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/services"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	renderPkg "github.com/unrolled/render"
	"go.uber.org/zap"
)

var render *renderPkg.Render

func init() {
	render = renderPkg.New()
}

type UserControllerInterface interface {
	CreateUser(resWriter http.ResponseWriter, req *http.Request)
	LoginUser(resWriter http.ResponseWriter, req *http.Request)
	GetAllUsers(resWriter http.ResponseWriter, req *http.Request)
	GetUserById(resWriter http.ResponseWriter, req *http.Request)
	UpdateUserById(resWriter http.ResponseWriter, req *http.Request)
	DeleteUserById(resWriter http.ResponseWriter, req *http.Request)
	GetProfile(resWriter http.ResponseWriter, req *http.Request)
	SendOtpForVerification(resWriter http.ResponseWriter, req *http.Request)
	VerifyOtp(resWriter http.ResponseWriter, req *http.Request)
	RefreshAccessToken(resWriter http.ResponseWriter, req *http.Request)
	LogoutUser(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService  services.UserServiceInterface
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(*dtos.CreateUserPayload)

	// call the create user service
	serviceErr := uc.UserService.CreateUser(userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while creating the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusCreated, resWriter, map[string]any{
		"success":  true,
		"message":  "Successfully created the user",
		"email":    userPayload.Email,
		"username": userPayload.Username,
	})
}

func (uc *UserController) LoginUser(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(*dtos.LoginUserPayload)

	// call the login user service
	accessToken, refreshToken, serviceErr := uc.UserService.LoginUser(userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while logging in",
			"error":   serviceErr.Error(),
		})

		return
	}

	// store the refresh token in cookies
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(resWriter, &cookie)

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Login was successful",
		"token":   accessToken,
	})
}

func (uc *UserController) GetAllUsers(resWriter http.ResponseWriter, req *http.Request) {
	// call the fetch all users service
	userModels, serviceErr := uc.UserService.GetAllUsers()

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting all users",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched all the users",
		"users":   userModels,
	})
}

func (uc *UserController) GetUserById(resWriter http.ResponseWriter, req *http.Request) {
	userParams := req.Context().Value("params").(*dtos.GetUserByIdParams)

	// call the fetch user by id service
	userModel, serviceErr := uc.UserService.GetUserById(userParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the user by id",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the user by id",
		"user":    userModel,
	})
}

func (uc *UserController) DeleteUserById(resWriter http.ResponseWriter, req *http.Request) {
	userParams := req.Context().Value("params").(*dtos.DeleteUserByIdParams)

	// call the delete user service
	serviceErr := uc.UserService.DeleteUserById(userParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while deleting the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully deleted the user",
	})
}

func (uc *UserController) UpdateUserById(resWriter http.ResponseWriter, req *http.Request) {
	userParams := req.Context().Value("params").(*dtos.UpdateUserByIdParams)
	userPayload := req.Context().Value("payload").(*dtos.UpdateUserByIdPayload)

	// call the update user service
	serviceErr := uc.UserService.UpdateUserById(userParams, userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while updating the user",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully updated the user",
	})
}

func (uc *UserController) GetProfile(resWriter http.ResponseWriter, req *http.Request) {
	userParams := req.Context().Value("params").(*dtos.GetUserByIdParams)

	// call the fetch user by id service
	userModel, serviceErr := uc.UserService.GetUserById(userParams)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while getting the profile",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully fetched the profile",
		"user":    userModel,
	})
}

func (uc *UserController) SendOtpForVerification(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(*dtos.CreateOtpServicePayload)

	// call the create user service
	serviceErr := uc.UserService.SendOtpForVerification(userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while sending otp for verification",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully sent the otp for verification",
		"email":   userPayload.UserEmail,
	})
}

func (uc *UserController) VerifyOtp(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(*dtos.VerifyOtpPayload)

	// call the verify otp user service
	serviceErr := uc.UserService.VerifyOtp(userPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while user verification",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Successfully verified the user by otp",
		"email":   userPayload.UserEmail,
	})
}

func (uc *UserController) RefreshAccessToken(resWriter http.ResponseWriter, req *http.Request) {

	// Retrieve the cookie
	cookie, err := req.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			utils.WriteJsonResponse(http.StatusNotFound, resWriter, map[string]any{
				"success": false,
				"message": "Cookie not found",
				"error":   err.Error(),
			})

			return
		}

		utils.WriteJsonResponse(http.StatusInternalServerError, resWriter, map[string]any{
			"success": false,
			"message": "Something went wrong while retrieving cookie",
			"error":   err.Error(),
		})

		return
	}

	// call the refresh access token service
	tokenPayload := &dtos.RefreshAccessTokenPayload{
		RefreshToken: cookie.Value,
	}

	accessToken, serviceErr := uc.UserService.RefreshAccessToken(tokenPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while refreshing access token",
			"error":   serviceErr.Error(),
		})

		return
	}

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "Refreshing access token was successful",
		"token":   accessToken,
	})

}

func (uc *UserController) LogoutUser(resWriter http.ResponseWriter, req *http.Request) {

	// Retrieve the cookie
	cookie, err := req.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			utils.WriteJsonResponse(http.StatusNotFound, resWriter, map[string]any{
				"success": false,
				"message": "Cookie not found",
				"error":   err.Error(),
			})

			return
		}

		utils.WriteJsonResponse(http.StatusInternalServerError, resWriter, map[string]any{
			"success": false,
			"message": "Something went wrong while retrieving cookie",
			"error":   err.Error(),
		})

		return
	}

	// call the logout service
	tokenPayload := &dtos.LogoutUserPayload{
		RefreshToken: cookie.Value,
	}

	serviceErr := uc.UserService.LogoutUser(tokenPayload)

	if serviceErr != nil {
		utils.WriteJsonResponse(serviceErr.StatusCode, resWriter, map[string]any{
			"success": serviceErr.Success,
			"message": "Something went wrong while logging out",
			"error":   serviceErr.Error(),
		})

		return
	}

	// clear cookie to remove refresh token
	cookie = &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(resWriter, cookie)

	utils.WriteJsonResponse(http.StatusOK, resWriter, map[string]any{
		"success": true,
		"message": "You have been successfully logged out",
	})
}

func NewUserController(service services.UserServiceInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserControllerInterface {
	newUserController := &UserController{
		UserService:  service,
		logger:       logger,
		serverConfig: serverConfig,
	}

	return newUserController
}
