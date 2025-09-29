package impl

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServiceImpl(userRepository *repository.UserRepository, messageService *service.MessageService, localGovernmentService *service.LocalGovernmentService, config configuration.Config) service.UserService {
	return &userServiceImpl{
		UserRepository:         *userRepository,
		MessageService:         *messageService,
		LocalGovernmentService: *localGovernmentService,
		Config:                 config,
	}
}

type userServiceImpl struct {
	repository.UserRepository
	service.MessageService
	service.LocalGovernmentService
	Config configuration.Config
}

func (u *userServiceImpl) FindByEmailOrPhone(ctx context.Context, userModel model.UserModel) entity.User {
	userResult, err := u.UserRepository.FindByEmailOrPhone(ctx, userModel)
	if err != nil {
		return entity.User{}
	}
	return userResult
}

func (u *userServiceImpl) GetAddresses(ctx context.Context, id uint) (interface{}, error) {
	addresses, err := u.UserRepository.FindAllAddress(ctx, id)
	exception.PanicLogging(err)
	return addresses, nil
}

func (u *userServiceImpl) SaveAddress(ctx context.Context, request model.AddressModel) (interface{}, error) {
	if request.Longitude == 0 || request.Latitude == 0 {
		locationGeometryResult := u.LocalGovernmentService.GetPlaceDetail(ctx, request.PlaceId)
		request.Longitude = locationGeometryResult.Result.Geometry.Location.Lng
		request.Latitude = locationGeometryResult.Result.Geometry.Location.Lat
	}

	address, err := u.UserRepository.SaveAddress(ctx, request)
	if err != nil {
		return nil, err
	}
	addressModel := model.AddressResponseModel{
		Id:          address.ID,
		Longitude:   common.ToFloat64(address.Longitude),
		Latitude:    common.ToFloat64(address.Latitude),
		PlaceId:     address.PlaceId,
		Description: address.Description,
	}
	return addressModel, nil
}

func (u *userServiceImpl) UpdateProfile(ctx context.Context, request model.UserModel, username string) (model.UserModel, error) {
	request.Username = username
	userResult, err := u.UserRepository.UpdateProfile(ctx, request)
	if err != nil {
		return model.UserModel{}, err
	}

	return userResult, nil
}

func (u *userServiceImpl) UpdateUserPassword(ctx context.Context, request model.ResetPasswordViewModel) model.UserModel {
	userResult, err := u.UserRepository.FindById(ctx, request.UserId)
	exception.PanicLogging(err)
	if userResult.ID == 0 {
		panic(exception.BadRequestError{
			Message: "User not found",
		})
	}
	if userResult.IsActive == false {
		panic(exception.BadRequestError{
			Message: "User is not active",
		})
	}

	otpModel := model.OtpModel{
		Code:      request.OTP,
		UserId:    request.UserId,
		Operation: "reset-password",
	}

	_, err = u.UserRepository.ValidateOtp(ctx, otpModel)
	exception.PanicLogging(err)

	passwordResetResult, err := u.UserRepository.SetPassword(request.UserId, request.Password)
	exception.PanicLogging(err)
	if passwordResetResult.Id == 0 {
		panic(exception.BadRequestError{
			Message: "User not found",
		})
	}

	return passwordResetResult
}

func (u *userServiceImpl) ResetPassword(ctx context.Context, request model.UserModel) model.UserModel {
	userResult, err := u.UserRepository.ResetPassword(ctx, request)
	if err != nil {
		panic(exception.BadRequestError{
			Message: "User not found",
		})
	}
	if userResult.Id == 0 {
		panic(exception.BadRequestError{
			Message: "User not found",
		})
	}
	otp, err := u.MessageService.GenerateOneTimePassword(ctx, userResult.Id)
	exception.PanicLogging(err)
	message := "Hello " + userResult.FirstName + "  your password reset otp is " + otp.Code + ", Do not share your this otp with a third-party."
	// Send SMS
	emailMessageModel := model.SMSMessageModel{
		PhoneNumber: request.PhoneNumber,
		CountryCode: userResult.AreaCode,
		Message:     message,
	}
	u.MessageService.SendSMS(ctx, emailMessageModel)
	return userResult
}

func (u *userServiceImpl) ValidateOtp(ctx context.Context, request model.OtpModel) entity.OneTimePassword {
	otp, err := u.UserRepository.ValidateOtp(ctx, request)
	if err != nil {
		panic(exception.BadRequestError{
			Message: "Invalid OTP",
		})
	}
	return otp
}

func (u *userServiceImpl) FindByID(ctx context.Context, id int) (model.UserModel, error) {
	user, err := u.UserRepository.FindById(ctx, id)
	exception.PanicLogging(err)
	userModel := model.UserModel{
		Id:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		EmailAddress: user.Email,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.UserRole,
		IsActive:     user.IsActive,
		FileName:     user.FileName,
	}
	return userModel, nil
}

func (u *userServiceImpl) GetClaimsFromToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.Config.Get("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (u *userServiceImpl) ChangePassword(ctx context.Context, token string, model model.ChangePasswordModel) entity.User {
	claims, err := u.GetClaimsFromToken(ctx, token)
	exception.PanicLogging(err)
	userResult, err := u.UserRepository.ChangePassword(ctx, claims, model)
	exception.PanicLogging(err)
	return userResult
}

func (u *userServiceImpl) SeedUser(ctx context.Context) {
	u.UserRepository.SeedUser(ctx)
}

func (u *userServiceImpl) Authentication(ctx context.Context, model model.LoginModel) entity.User {
	userResult, err := u.UserRepository.Authentication(ctx, model.Username)
	if err != nil {
		panic(exception.UnauthorizedError{
			Message: err.Error(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(model.Password))
	if err != nil {
		panic(exception.UnauthorizedError{
			Message: "incorrect username and password",
		})
	}
	if userResult.IsActive == false {
		panic(exception.UnauthorizedError{
			Message: "User is not active",
		})
	}
	return userResult
}

func (u *userServiceImpl) Register(ctx context.Context, userModel model.UserModel) entity.User {
	userModel.IsActive = true
	if userModel.Password == "" {
		// Generate a random password if not provided
		userModel.Password, _ = common.GeneratePassword(8)
	}
	user, err := u.UserRepository.Create(userModel)
	exception.PanicLogging(err)
	template := u.MessageService.FindMessageTemplateByName(ctx, userModel.Role+"_template")

	data := map[string]string{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Username":  user.Username,
		"Password":  userModel.Password,
		"Email":     user.Email,
		"Url":       "http://localhost:9999",
	}
	message := ReplacePlaceholders(template.Message, data)
	// Send email
	emailMessageModel := model.EmailMessageModel{
		To:      user.Email,
		Subject: template.Subject,
		Message: message,
	}
	u.MessageService.SendEmail(ctx, emailMessageModel)
	return user
}

func (u *userServiceImpl) RegisterCustomer(ctx context.Context, userModel model.UserModel) interface{} {
	userModel.IsActive = true
	userModel.Role = "customer"
	user, err := u.UserRepository.Create(userModel)
	exception.PanicLogging(err)

	// Send email
	otp, err := u.MessageService.GenerateOneTimePassword(ctx, user.ID)
	exception.PanicLogging(err)
	message := "Hello " + user.FirstName + "   your OTP is " + otp.Code + ", Do not share your this otp with a third-party."
	// Send SMS
	emailMessageModel := model.SMSMessageModel{
		PhoneNumber: user.PhoneNumber,
		Message:     message,
		CountryCode: user.AreaCode,
	}
	user.Password = ""
	u.MessageService.SendSMS(ctx, emailMessageModel)
	userModel = model.UserModel{
		Id:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		EmailAddress: user.Email,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.UserRole,
		IsActive:     user.IsActive,
		OtpCode:      otp.Code,
	}
	return userModel
}

func (u *userServiceImpl) Create(ctx context.Context, model model.UserModel, file *multipart.FileHeader) entity.User {
	model.IsActive = false
	user, err := u.UserRepository.Create(model)
	exception.PanicLogging(err)
	return user
}

func (u *userServiceImpl) List(ctx context.Context) ([]entity.User, error) {
	return u.UserRepository.List(ctx)
}
func ReplacePlaceholders(template string, placeholders map[string]string) string {
	for key, value := range placeholders {
		placeholder := fmt.Sprintf("{{%s}}", key)
		template = strings.Replace(template, placeholder, value, -1)
	}
	return template
}

func (u *userServiceImpl) UpdateFcmToken(ctx context.Context, request model.UpdateFcmToken) error {
	err := u.UserRepository.UpdateFcmToken(ctx, request)
	return err
}
