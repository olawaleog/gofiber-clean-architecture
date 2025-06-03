package impl

import (
	"context"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"strings"
)

func NewUserServiceImpl(userRepository *repository.UserRepository, messageService *service.MessageService) service.UserService {
	return &userServiceImpl{UserRepository: *userRepository, MessageService: *messageService}
}

type userServiceImpl struct {
	repository.UserRepository
	service.MessageService
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
	address, err := u.UserRepository.SaveAddress(ctx, request)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (u *userServiceImpl) UpdateProfile(ctx context.Context, request model.UserModel, token string) (model.UserModel, error) {
	claims, err := u.GetClaimsFromToken(ctx, token)
	exception.PanicLogging(err)
	request.Username = claims["username"].(string)
	userResult, err := u.UserRepository.UpdateProfile(ctx, request)
	exception.PanicLogging(err)

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
		CountryCode: userResult.CountryCode,
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

func (userService *userServiceImpl) GetClaimsFromToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("missing token")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (userService *userServiceImpl) ChangePassword(ctx context.Context, token string, model model.ChangePasswordModel) entity.User {
	claims, err := userService.GetClaimsFromToken(ctx, token)
	exception.PanicLogging(err)
	userResult, err := userService.UserRepository.ChangePassword(ctx, claims, model)
	exception.PanicLogging(err)
	return userResult
}

func (userService *userServiceImpl) SeedUser(ctx context.Context) {
	userService.UserRepository.SeedUser(ctx)
}

func (userService *userServiceImpl) Authentication(ctx context.Context, model model.LoginModel) entity.User {
	userResult, err := userService.UserRepository.Authentication(ctx, model.Username)
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

func (userService *userServiceImpl) Register(ctx context.Context, userModel model.UserModel) entity.User {
	userModel.IsActive = true
	user, err := userService.UserRepository.Create(userModel)
	exception.PanicLogging(err)
	template := userService.MessageService.FindMessageTemplateByName(ctx, userModel.Role+"_template")

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
	userService.MessageService.SendEmail(ctx, emailMessageModel)
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

func (userService *userServiceImpl) Create(ctx context.Context, model model.UserModel, file *multipart.FileHeader) entity.User {
	model.IsActive = false
	user, err := userService.UserRepository.Create(model)
	exception.PanicLogging(err)
	return user
}

func (userService *userServiceImpl) List(ctx context.Context) ([]entity.User, error) {
	return userService.UserRepository.List(ctx)
}
func ReplacePlaceholders(template string, placeholders map[string]string) string {
	for key, value := range placeholders {
		placeholder := fmt.Sprintf("{{%s}}", key)
		template = strings.Replace(template, placeholder, value, -1)
	}
	return template
}
