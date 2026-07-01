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
	LogoutUserFromAllSessions(resWriter http.ResponseWriter, req *http.Request)
}

type UserController struct {
	UserService  services.UserServiceInterface
	logger       *zap.Logger
	serverConfig *config.ServerConfig
}

func (uc *UserController) CreateUser(resWriter http.ResponseWriter, req *http.Request) {
	userPayload := req.Context().Value("payload").(*dtos.CreateUserPayload)

	// call the create user service
	createUserServiceErr := uc.UserService.CreateUser(userPayload)

	if createUserServiceErr != nil {
		utils.WriteJsonResponse(createUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": createUserServiceErr.Success,
			"message": "Something went wrong while creating the user",
			"error":   createUserServiceErr.Error(),
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
	accessToken, refreshToken, loginUserServiceErr := uc.UserService.LoginUser(userPayload)

	if loginUserServiceErr != nil {
		utils.WriteJsonResponse(loginUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": loginUserServiceErr.Success,
			"message": "Something went wrong while logging in",
			"error":   loginUserServiceErr.Error(),
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
	userModels, getUsersServiceErr := uc.UserService.GetAllUsers()

	if getUsersServiceErr != nil {
		utils.WriteJsonResponse(getUsersServiceErr.StatusCode, resWriter, map[string]any{
			"success": getUsersServiceErr.Success,
			"message": "Something went wrong while getting all users",
			"error":   getUsersServiceErr.Error(),
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
	userModel, getUserServiceErr := uc.UserService.GetUserById(userParams)

	if getUserServiceErr != nil {
		utils.WriteJsonResponse(getUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": getUserServiceErr.Success,
			"message": "Something went wrong while getting the user by id",
			"error":   getUserServiceErr.Error(),
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
	deleteUserServiceErr := uc.UserService.DeleteUserById(userParams)

	if deleteUserServiceErr != nil {
		utils.WriteJsonResponse(deleteUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": deleteUserServiceErr.Success,
			"message": "Something went wrong while deleting the user",
			"error":   deleteUserServiceErr.Error(),
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
	updateUserServiceErr := uc.UserService.UpdateUserById(userParams, userPayload)

	if updateUserServiceErr != nil {
		utils.WriteJsonResponse(updateUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": updateUserServiceErr.Success,
			"message": "Something went wrong while updating the user",
			"error":   updateUserServiceErr.Error(),
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
	userModel, getUserServiceErr := uc.UserService.GetUserById(userParams)

	if getUserServiceErr != nil {
		utils.WriteJsonResponse(getUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": getUserServiceErr.Success,
			"message": "Something went wrong while getting the profile",
			"error":   getUserServiceErr.Error(),
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
	sendOtpServiceErr := uc.UserService.SendOtpForVerification(userPayload)

	if sendOtpServiceErr != nil {
		utils.WriteJsonResponse(sendOtpServiceErr.StatusCode, resWriter, map[string]any{
			"success": sendOtpServiceErr.Success,
			"message": "Something went wrong while sending otp for verification",
			"error":   sendOtpServiceErr.Error(),
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
	verifyOtpServiceErr := uc.UserService.VerifyOtp(userPayload)

	if verifyOtpServiceErr != nil {
		utils.WriteJsonResponse(verifyOtpServiceErr.StatusCode, resWriter, map[string]any{
			"success": verifyOtpServiceErr.Success,
			"message": "Something went wrong while user verification",
			"error":   verifyOtpServiceErr.Error(),
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
	cookie, cookieErr := req.Cookie("refresh_token")
	if cookieErr != nil {
		if errors.Is(cookieErr, http.ErrNoCookie) {
			utils.WriteJsonResponse(http.StatusNotFound, resWriter, map[string]any{
				"success": false,
				"message": "Missing tokens, please login",
				"error":   cookieErr.Error(),
			})

			return
		}

		utils.WriteJsonResponse(http.StatusInternalServerError, resWriter, map[string]any{
			"success": false,
			"message": "Something went wrong while retrieving tokens",
			"error":   cookieErr.Error(),
		})

		return
	}

	// call the refresh access token service
	tokenPayload := &dtos.RefreshAccessTokenPayload{
		RefreshToken: cookie.Value,
	}

	accessToken, refreshTokenServiceErr := uc.UserService.RefreshAccessToken(tokenPayload)

	if refreshTokenServiceErr != nil {
		utils.WriteJsonResponse(refreshTokenServiceErr.StatusCode, resWriter, map[string]any{
			"success": refreshTokenServiceErr.Success,
			"message": "Something went wrong while refreshing access token",
			"error":   refreshTokenServiceErr.Error(),
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
	cookie, cookieErr := req.Cookie("refresh_token")
	if cookieErr != nil {
		if errors.Is(cookieErr, http.ErrNoCookie) {
			utils.WriteJsonResponse(http.StatusNotFound, resWriter, map[string]any{
				"success": false,
				"message": "Missing tokens, please login",
				"error":   cookieErr.Error(),
			})

			return
		}

		utils.WriteJsonResponse(http.StatusInternalServerError, resWriter, map[string]any{
			"success": false,
			"message": "Something went wrong while retrieving tokens",
			"error":   cookieErr.Error(),
		})

		return
	}

	// call the logout service
	tokenPayload := &dtos.LogoutUserPayload{
		RefreshToken: cookie.Value,
	}

	logoutUserServiceErr := uc.UserService.LogoutUser(tokenPayload)

	if logoutUserServiceErr != nil {
		utils.WriteJsonResponse(logoutUserServiceErr.StatusCode, resWriter, map[string]any{
			"success": logoutUserServiceErr.Success,
			"message": "Something went wrong while logging out",
			"error":   logoutUserServiceErr.Error(),
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

func (uc *UserController) LogoutUserFromAllSessions(resWriter http.ResponseWriter, req *http.Request) {

	// Retrieve the cookie
	cookie, cookieErr := req.Cookie("refresh_token")
	if cookieErr != nil {
		if errors.Is(cookieErr, http.ErrNoCookie) {
			utils.WriteJsonResponse(http.StatusNotFound, resWriter, map[string]any{
				"success": false,
				"message": "Missing tokens, please login",
				"error":   cookieErr.Error(),
			})

			return
		}

		utils.WriteJsonResponse(http.StatusInternalServerError, resWriter, map[string]any{
			"success": false,
			"message": "Something went wrong while retrieving tokens",
			"error":   cookieErr.Error(),
		})

		return
	}

	// call the logout from all sessions service
	tokenPayload := &dtos.LogoutUserFromAllSessionsPayload{
		RefreshToken: cookie.Value,
	}

	logoutUserFromSessionsServiceErr := uc.UserService.LogoutUserFromAllSessions(tokenPayload)

	if logoutUserFromSessionsServiceErr != nil {
		utils.WriteJsonResponse(logoutUserFromSessionsServiceErr.StatusCode, resWriter, map[string]any{
			"success": logoutUserFromSessionsServiceErr.Success,
			"message": "Something went wrong while logging out from all sessions",
			"error":   logoutUserFromSessionsServiceErr.Error(),
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
		"message": "You have been successfully logged out from all sessions",
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
