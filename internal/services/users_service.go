package services

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/config"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/database/models"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/dtos"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/repositories"
	"github.com/ayushWeb07/AirBnb-Go-Api-Gateway/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(userPayload *dtos.CreateUserPayload) *utils.AppError
	LoginUser(userPayload *dtos.LoginUserPayload) (string, string, *utils.AppError)
	GetAllUsers() ([]*models.UserModel, *utils.AppError)
	GetUserById(userParams *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError)
	UpdateUserById(userParams *dtos.UpdateUserByIdParams, userPayload *dtos.UpdateUserByIdPayload) *utils.AppError
	DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError
	SendOtpForVerification(otpPayload *dtos.CreateOtpServicePayload) *utils.AppError
	VerifyOtp(otpPayload *dtos.VerifyOtpPayload) *utils.AppError
	RefreshAccessToken(tokenPayload *dtos.RefreshAccessTokenPayload) (string, *utils.AppError)
	LogoutUser(tokenPayload *dtos.LogoutUserPayload) *utils.AppError
	LogoutUserFromAllSessions(tokenPayload *dtos.LogoutUserFromAllSessionsPayload) *utils.AppError
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (us *UserService) CreateUser(userPayload *dtos.CreateUserPayload) *utils.AppError {
	us.logger.Info("Create user service called...")

	// check if the user already exists
	_, getUserRepositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmailPayload{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if getUserRepositoryErr == nil {
		return utils.BadRequest("User with such username and email, already exists")
	}

	// hash the password
	bytes, hashErr := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		us.logger.Fatal("Something went wrong while hashing the password",
			zap.String("error", hashErr.Error()))

		return utils.InternalServerError("Something went wrong while hashing the password: " + hashErr.Error())
	}

	userPayload.Password = string(bytes)

	// call the create user repository
	createUserRepositoryErr := us.UserRepository.CreateUser(userPayload)

	if createUserRepositoryErr != nil {
		return createUserRepositoryErr
	}

	us.logger.Info("Create user service was successful")

	return nil
}

func (us *UserService) LoginUser(userPayload *dtos.LoginUserPayload) (string, string, *utils.AppError) {
	us.logger.Info("Login user service called...")

	// fetch the user by username and email repository
	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmailPayload{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if getUserRepositoryErr != nil {
		return "", "", getUserRepositoryErr
	}

	// check if user is verified
	if !existingUserModel.Verified {
		us.logger.Error("User is not verified",
			zap.Int("user_id", existingUserModel.ID))

		return "", "", utils.Forbidden("You must first verify your email address to login")
	}

	// check if passwords match
	compareErr := bcrypt.CompareHashAndPassword([]byte(existingUserModel.Password), []byte(userPayload.Password))

	if compareErr != nil {
		us.logger.Error("Invalid password has been provided",
			zap.String("error", compareErr.Error()))

		return "", "", utils.BadRequest("Invalid password has been provided")
	}

	// generate the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUserModel.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	accessTokenString, accessTokenErr := accessToken.SignedString([]byte(us.serverConfig.AccessTokenSecretKey))

	if accessTokenErr != nil {
		us.logger.Fatal("Something went wrong while generating the access token",
			zap.String("error", accessTokenErr.Error()))

		return "", "", utils.InternalServerError("Something went wrong while generating the access token: " + accessTokenErr.Error())
	}

	// generate the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUserModel.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenString, refreshTokenErr := refreshToken.SignedString([]byte(us.serverConfig.RefreshTokenSecretKey))

	if refreshTokenErr != nil {
		us.logger.Fatal("Something went wrong while generating the refresh token",
			zap.String("error", refreshTokenErr.Error()))

		return "", "", utils.InternalServerError("Something went wrong while generating the refresh token: " + refreshTokenErr.Error())
	}

	// generate the session
	refreshTokenBytes := sha256.Sum256([]byte(refreshTokenString))
	refreshTokenHash := hex.EncodeToString(refreshTokenBytes[:])

	sessionPayload := &dtos.CreateSessionPayload{
		UserID:           existingUserModel.ID,
		RefreshTokenHash: refreshTokenHash,
	}

	createSessionRepositoryErr := us.UserRepository.CreateSession(sessionPayload)

	if createSessionRepositoryErr != nil {
		us.logger.Fatal("Something went wrong while generating the session")

		return "", "", createSessionRepositoryErr
	}

	us.logger.Info("Login user service was successful",
		zap.Int("user_id", existingUserModel.ID))

	return accessTokenString, refreshTokenString, nil
}

func (us *UserService) GetAllUsers() ([]*models.UserModel, *utils.AppError) {
	us.logger.Info("Get all users service called...")

	// call the fetch all users repository
	userModels, getUsersRepositoryErr := us.UserRepository.GetAllUsers()
	return userModels, getUsersRepositoryErr
}

func (us *UserService) GetUserById(userParams *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError) {
	us.logger.Info("Get by id user service called...")

	// call the fetch user by id repository
	userModel, getUserRepositoryErr := us.UserRepository.GetUserById(userParams)
	return userModel, getUserRepositoryErr
}

func (us *UserService) UpdateUserById(userParams *dtos.UpdateUserByIdParams, userPayload *dtos.UpdateUserByIdPayload) *utils.AppError {
	us.logger.Info("Update user service called...")

	// call the update user by id repository
	updateUserRepositoryErr := us.UserRepository.UpdateUserById(userParams, userPayload)
	return updateUserRepositoryErr
}

func (us *UserService) DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError {
	us.logger.Info("Delete user service called...")

	// call the delete user by id repository
	deleteUserRepositoryErr := us.UserRepository.DeleteUserById(userParams)
	return deleteUserRepositoryErr
}

func (us *UserService) SendOtpForVerification(otpPayload *dtos.CreateOtpServicePayload) *utils.AppError {
	us.logger.Info("Send otp for verification service called...")

	// fetch the user by email repository
	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserByEmail(&dtos.GetUserByEmailPayload{
		Email: otpPayload.UserEmail,
	})

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if user is verified
	if existingUserModel.Verified {
		us.logger.Error("User is already verified",
			zap.Int("user_id", existingUserModel.ID))

		return utils.Forbidden("Your email address is already verified")
	}

	// generate otp and hash
	otpString, otpGenerationErr := utils.GenerateRandomOtp(10)

	if otpGenerationErr != nil {
		us.logger.Error("Something went wrong while generating otp")

		return utils.InternalServerError("Something went wrong while generating otp: " + otpGenerationErr.Error())
	}

	otpBytes := sha256.Sum256([]byte(otpString))
	otpHash := hex.EncodeToString(otpBytes[:])

	// insert otp into db
	otpRepoPayload := &dtos.CreateOtpRepoPayload{
		UserID:    existingUserModel.ID,
		UserEmail: existingUserModel.Email,
		OtpHash:   otpHash,
	}

	createOtpRepositoryErr := us.UserRepository.CreateOtp(otpRepoPayload)

	if createOtpRepositoryErr != nil {
		us.logger.Fatal("Something went wrong while generating the otp")

		return createOtpRepositoryErr
	}

	// render the email template
	result, emailTemplateErr := utils.Render("otp_verification.txt", &utils.EmailData{
		UserName: existingUserModel.Username,
		AppName:  "Hajjme No Ippo",
		Otp:      otpString,
	})

	if emailTemplateErr != nil {
		us.logger.Fatal("Something went wrong while rendering the email template")

		return utils.InternalServerError("Something went wrong while rendering the email template: " + emailTemplateErr.Error())
	}

	// send the email
	emailErr := sendEmail(otpPayload.UserEmail, "Complete Your Account Verification", result)

	if emailErr != nil {
		us.logger.Fatal("Something went wrong while sending the otp verification email")

		return utils.InternalServerError("Something went wrong while sending the otp verification email: " + emailErr.Error())
	}

	us.logger.Info("Sending otp for verification was successful")

	return nil
}

func (us *UserService) VerifyOtp(otpPayload *dtos.VerifyOtpPayload) *utils.AppError {
	us.logger.Info("Verify otp service called...")

	// fetch the user by email repository
	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserByEmail(&dtos.GetUserByEmailPayload{
		Email: otpPayload.UserEmail,
	})

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if user is verified
	if existingUserModel.Verified {
		us.logger.Error("User is already verified",
			zap.Int("user_id", existingUserModel.ID))

		return utils.Forbidden("Your email address is already verified")
	}

	// hash the user sent otp
	otpBytes := sha256.Sum256([]byte(otpPayload.Otp))
	otpHash := hex.EncodeToString(otpBytes[:])

	// find the otp in the db
	otpRepoPayload := &dtos.FetchOtpRepoPayload{
		UserID:    existingUserModel.ID,
		UserEmail: existingUserModel.Email,
		OtpHash:   otpHash,
	}

	_, fetchOtpRepositoryErr := us.UserRepository.FetchOtp(otpRepoPayload)

	if fetchOtpRepositoryErr != nil {
		return fetchOtpRepositoryErr
	}

	// update user status to verified
	userParams := &dtos.UpdateUserByIdParams{
		ID: existingUserModel.ID,
	}

	userPayload := &dtos.UpdateUserByIdPayload{
		Verified: true,
	}

	updateUserRepositoryErr := us.UserRepository.UpdateUserById(userParams, userPayload)

	if updateUserRepositoryErr != nil {
		return updateUserRepositoryErr
	}

	// delete all the otps for the user
	deleteOtpsPayload := &dtos.DeleteOtpsRepoPayload{
		UserID: existingUserModel.ID,
	}

	deleteOtpsRepositoryErr := us.UserRepository.DeleteOtps(deleteOtpsPayload)

	if deleteOtpsRepositoryErr != nil {
		return deleteOtpsRepositoryErr
	}

	us.logger.Info("User verification was successful")

	return nil
}

func (us *UserService) RefreshAccessToken(tokenPayload *dtos.RefreshAccessTokenPayload) (string, *utils.AppError) {
	// verify the refresh token
	token, tokenErr := jwt.Parse(tokenPayload.RefreshToken, func(token *jwt.Token) (any, error) {
		// invalid signing method had been used for token generating
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			us.logger.Error("Invalid signing method had been used while refreshing access token")

			return utils.BadRequest("Invalid signing method had been used"), nil
		}

		// else return the jwt secret key
		return []byte(us.serverConfig.RefreshTokenSecretKey), nil
	})

	// check if there's an error or the token is invalid
	if tokenErr != nil {
		us.logger.Fatal("Invalid token has been provided while refreshing access token",
			zap.String("error", tokenErr.Error()))

		return "", utils.Unauthorized("Invalid token has been provided while refreshing access token:  " + tokenErr.Error())
	}

	if !token.Valid {
		us.logger.Fatal("Invalid or expired token has been provided while refreshing access token")

		return "", utils.Unauthorized("Invalid or expired token has been provided while refreshing access token")
	}

	// parse token to decode the payload
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		us.logger.Fatal("Failed to decode payload from an invalid token")

		return "", utils.Unauthorized("Failed to decode payload from an invalid token")
	}

	// access the payload
	userId := claims["user_id"].(float64)
	expiryTime := claims["exp"].(float64)

	// check if token has expired
	if time.Now().Unix() > int64(expiryTime) {
		us.logger.Fatal("Token has expired. Please login again")

		return "", utils.Unauthorized("Token has expired. Please login again")
	}

	// fetch the user
	getUserPayload := &dtos.GetUserByIdParams{
		ID: int(userId),
	}

	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserById(getUserPayload)

	if getUserRepositoryErr != nil {
		return "", getUserRepositoryErr
	}

	// check if user is verified
	if !existingUserModel.Verified {
		us.logger.Error("User is not verified",
			zap.Int("user_id", existingUserModel.ID))

		return "", utils.Forbidden("You must first verify your email address to login")
	}

	// fetch the session
	refreshTokenBytes := sha256.Sum256([]byte(tokenPayload.RefreshToken))
	refreshTokenHash := hex.EncodeToString(refreshTokenBytes[:])

	sessionPayload := &dtos.FetchSessionPayload{
		UserID:           existingUserModel.ID,
		RefreshTokenHash: refreshTokenHash,
		Revoked:          false,
	}

	_, fetchSessionRepositoryErr := us.UserRepository.FetchSession(sessionPayload)

	if fetchSessionRepositoryErr != nil {
		return "", fetchSessionRepositoryErr
	}

	// generate the access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUserModel.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	accessTokenString, accessTokenErr := accessToken.SignedString([]byte(us.serverConfig.AccessTokenSecretKey))

	if accessTokenErr != nil {
		us.logger.Fatal("Something went wrong while generating the access token",
			zap.String("error", accessTokenErr.Error()))

		return "", utils.InternalServerError("Something went wrong while generating the access token: " + accessTokenErr.Error())
	}

	return accessTokenString, nil
}

func (us *UserService) LogoutUser(tokenPayload *dtos.LogoutUserPayload) *utils.AppError {
	// verify the refresh token
	token, tokenErr := jwt.Parse(tokenPayload.RefreshToken, func(token *jwt.Token) (any, error) {
		// invalid signing method had been used for token generating
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			us.logger.Error("Invalid signing method had been used while logging out")

			return utils.BadRequest("Invalid signing method had been used"), nil
		}

		// else return the jwt secret key
		return []byte(us.serverConfig.RefreshTokenSecretKey), nil
	})

	// check if there's an error or the token is invalid
	if tokenErr != nil {
		us.logger.Fatal("Invalid token has been provided while logging out",
			zap.String("error", tokenErr.Error()))

		return utils.Unauthorized("Invalid token has been provided while logging out:  " + tokenErr.Error())
	}

	if !token.Valid {
		us.logger.Fatal("Invalid or expired token has been provided while logging out")

		return utils.Unauthorized("Invalid or expired token has been provided while generating the refresh token")
	}

	// parse token to decode the payload
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		us.logger.Fatal("Failed to decode payload from an invalid token")

		return utils.Unauthorized("Failed to decode payload from an invalid token")
	}

	// access the payload
	userId := claims["user_id"].(float64)
	expiryTime := claims["exp"].(float64)

	// check if token has expired
	if time.Now().Unix() > int64(expiryTime) {
		us.logger.Fatal("Token has expired. Please login again")

		return utils.Unauthorized("Token has expired. Please login again")
	}

	// fetch the user
	getUserPayload := &dtos.GetUserByIdParams{
		ID: int(userId),
	}

	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserById(getUserPayload)

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if user is verified
	if !existingUserModel.Verified {
		us.logger.Error("User is not verified",
			zap.Int("user_id", existingUserModel.ID))

		return utils.Forbidden("You must first verify your email address to login")
	}

	// fetch the session
	refreshTokenBytes := sha256.Sum256([]byte(tokenPayload.RefreshToken))
	refreshTokenHash := hex.EncodeToString(refreshTokenBytes[:])

	sessionPayload := &dtos.FetchSessionPayload{
		UserID:           existingUserModel.ID,
		RefreshTokenHash: refreshTokenHash,
		Revoked:          false,
	}

	existingSessionModel, fetchSessionRepositoryErr := us.UserRepository.FetchSession(sessionPayload)

	if fetchSessionRepositoryErr != nil {
		return fetchSessionRepositoryErr
	}

	// revoke the session
	updateSessionParams := &dtos.UpdateSessionByIdParams{
		ID: existingSessionModel.ID,
	}

	updateSessionPayload := &dtos.UpdateSessionByIdPayload{
		Revoked: true,
	}

	updateSessionRepositoryErr := us.UserRepository.UpdateSessionById(updateSessionParams, updateSessionPayload)

	if updateSessionRepositoryErr != nil {
		return updateSessionRepositoryErr
	}

	us.logger.Info("Logout user service was successful",
		zap.Int("user_id", existingUserModel.ID))

	return nil
}

func (us *UserService) LogoutUserFromAllSessions(tokenPayload *dtos.LogoutUserFromAllSessionsPayload) *utils.AppError {
	// verify the refresh token
	token, tokenErr := jwt.Parse(tokenPayload.RefreshToken, func(token *jwt.Token) (any, error) {
		// invalid signing method had been used for token generating
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			us.logger.Error("Invalid signing method had been used while logging out from all sessions")

			return utils.BadRequest("Invalid signing method had been used"), nil
		}

		// else return the jwt secret key
		return []byte(us.serverConfig.RefreshTokenSecretKey), nil
	})

	// check if there's an error or the token is invalid
	if tokenErr != nil {
		us.logger.Fatal("Invalid token has been provided while logging out from all sessions",
			zap.String("error", tokenErr.Error()))

		return utils.Unauthorized("Invalid token has been provided while logging out from all sessions:  " + tokenErr.Error())
	}

	if !token.Valid {
		us.logger.Fatal("Invalid or expired token has been provided while logging out from all sessions")

		return utils.Unauthorized("Invalid or expired token has been provided while logging out from all sessions")
	}

	// parse token to decode the payload
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		us.logger.Fatal("Failed to decode payload from an invalid token")

		return utils.Unauthorized("Failed to decode payload from an invalid token")
	}

	// access the payload
	userId := claims["user_id"].(float64)
	expiryTime := claims["exp"].(float64)

	// check if token has expired
	if time.Now().Unix() > int64(expiryTime) {
		us.logger.Fatal("Token has expired. Please login again")

		return utils.Unauthorized("Token has expired. Please login again")
	}

	// fetch the user
	getUserPayload := &dtos.GetUserByIdParams{
		ID: int(userId),
	}

	existingUserModel, getUserRepositoryErr := us.UserRepository.GetUserById(getUserPayload)

	if getUserRepositoryErr != nil {
		return getUserRepositoryErr
	}

	// check if user is verified
	if !existingUserModel.Verified {
		us.logger.Error("User is not verified",
			zap.Int("user_id", existingUserModel.ID))

		return utils.Forbidden("You must first verify your email address to login")
	}

	// revoke all sessions
	sessionPayload := &dtos.RevokeAllSessionsPayload{
		UserID: existingUserModel.ID,
	}

	revokeSessionsRepositoryErr := us.UserRepository.RevokeAllSessions(sessionPayload)

	if revokeSessionsRepositoryErr != nil {
		return revokeSessionsRepositoryErr
	}

	us.logger.Info("Logout from all sessions user service was successful",
		zap.Int("user_id", existingUserModel.ID))

	return nil
}

func NewUserService(repo repositories.UserRepositoryInterface, logger *zap.Logger, serverConfig *config.ServerConfig) UserServiceInterface {
	newUserService := &UserService{
		UserRepository: repo,
		logger:         logger,
		serverConfig:   serverConfig,
	}

	return newUserService
}
