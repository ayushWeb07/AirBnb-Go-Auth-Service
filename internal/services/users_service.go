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
}

type UserService struct {
	UserRepository repositories.UserRepositoryInterface
	logger         *zap.Logger
	serverConfig   *config.ServerConfig
}

func (us *UserService) CreateUser(userPayload *dtos.CreateUserPayload) *utils.AppError {
	us.logger.Info("Create user service called...")

	// check if the user already exists
	_, repositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmailPayload{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if repositoryErr == nil {
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
	repositoryErr = us.UserRepository.CreateUser(userPayload)

	if repositoryErr != nil {
		return repositoryErr
	}

	us.logger.Info("Create user service was successful")

	return nil
}

func (us *UserService) LoginUser(userPayload *dtos.LoginUserPayload) (string, string, *utils.AppError) {
	us.logger.Info("Login user service called...")

	// fetch the user by username and email repository
	existingUserModel, repositoryErr := us.UserRepository.GetUserByUsernameAndEmail(&dtos.GetUserByUsernameAndEmailPayload{
		Username: userPayload.Username,
		Email:    userPayload.Email,
	})

	if repositoryErr != nil {
		return "", "", repositoryErr
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

	repositoryErr = us.UserRepository.CreateSession(sessionPayload)

	if repositoryErr != nil {
		us.logger.Fatal("Something went wrong while generating the session")

		return "", "", repositoryErr
	}

	us.logger.Info("Login user service was successful",
		zap.Int("user_id", existingUserModel.ID))

	return accessTokenString, refreshTokenString, nil
}

func (us *UserService) GetAllUsers() ([]*models.UserModel, *utils.AppError) {
	us.logger.Info("Get all users service called...")

	// call the fetch all users repository
	userModels, repositoryErr := us.UserRepository.GetAllUsers()
	return userModels, repositoryErr
}

func (us *UserService) GetUserById(userParams *dtos.GetUserByIdParams) (*models.UserModel, *utils.AppError) {
	us.logger.Info("Get by id user service called...")

	// call the fetch user by id repository
	userModel, repositoryErr := us.UserRepository.GetUserById(userParams)
	return userModel, repositoryErr
}

func (us *UserService) UpdateUserById(userParams *dtos.UpdateUserByIdParams, userPayload *dtos.UpdateUserByIdPayload) *utils.AppError {
	us.logger.Info("Update user service called...")

	// call the update user by id repository
	repositoryErr := us.UserRepository.UpdateUserById(userParams, userPayload)
	return repositoryErr
}

func (us *UserService) DeleteUserById(userParams *dtos.DeleteUserByIdParams) *utils.AppError {
	us.logger.Info("Delete user service called...")

	// call the delete user by id repository
	repositoryErr := us.UserRepository.DeleteUserById(userParams)
	return repositoryErr
}

func (us *UserService) SendOtpForVerification(otpPayload *dtos.CreateOtpServicePayload) *utils.AppError {
	us.logger.Info("Send otp for verification service called...")

	// fetch the user by email repository
	existingUserModel, repositoryErr := us.UserRepository.GetUserByEmail(&dtos.GetUserByEmailPayload{
		Email: otpPayload.UserEmail,
	})

	if repositoryErr != nil {
		return repositoryErr
	}

	// check if user is verified
	if existingUserModel.Verified {
		us.logger.Error("User is already verified",
			zap.Int("user_id", existingUserModel.ID))

		return utils.Forbidden("Your email address is already verified")
	}

	// generate otp and hash
	otpString, err := utils.GenerateRandomOtp(10)

	if err != nil {
		us.logger.Error("Something went wrong while generating otp")

		return utils.InternalServerError("Something went wrong while generating otp: " + err.Error())
	}

	otpBytes := sha256.Sum256([]byte(otpString))
	otpHash := hex.EncodeToString(otpBytes[:])

	// insert otp into db
	otpRepoPayload := &dtos.CreateOtpRepoPayload{
		UserID:    existingUserModel.ID,
		UserEmail: existingUserModel.Email,
		OtpHash:   otpHash,
	}

	repositoryErr = us.UserRepository.CreateOtp(otpRepoPayload)

	if repositoryErr != nil {
		us.logger.Fatal("Something went wrong while generating the otp")

		return repositoryErr
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
	existingUserModel, repositoryErr := us.UserRepository.GetUserByEmail(&dtos.GetUserByEmailPayload{
		Email: otpPayload.UserEmail,
	})

	if repositoryErr != nil {
		return repositoryErr
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

	_, repositoryErr = us.UserRepository.FetchOtp(otpRepoPayload)

	if repositoryErr != nil {
		return repositoryErr
	}

	// update user status to verified
	userParams := &dtos.UpdateUserByIdParams{
		ID: existingUserModel.ID,
	}

	userPayload := &dtos.UpdateUserByIdPayload{
		Verified: true,
	}

	repositoryErr = us.UserRepository.UpdateUserById(userParams, userPayload)

	if repositoryErr != nil {
		return repositoryErr
	}

	// delete all the otps for the user
	deleteOtpsPayload := &dtos.DeleteOtpsRepoPayload{
		UserID: existingUserModel.ID,
	}

	repositoryErr = us.UserRepository.DeleteOtps(deleteOtpsPayload)

	if repositoryErr != nil {
		return repositoryErr
	}

	us.logger.Info("User verification was successful")

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
