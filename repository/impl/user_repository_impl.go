package impl

import (
	"context"
	"errors"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
	"time"
)

func NewUserRepositoryImpl(DB *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{DB: DB}
}

type userRepositoryImpl struct {
	*gorm.DB
}

func (u *userRepositoryImpl) FindByEmailOrPhone(ctx context.Context, userModel model.UserModel) (entity.User, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("username = ? or phone_number = ? or email = ?", userModel.PhoneNumber, userModel.PhoneNumber, userModel.EmailAddress).First(&user).Error
	if err != nil {
		return entity.User{}, errors.New("user not found")
	}

	return user, nil
}

func (u *userRepositoryImpl) FindAllAddress(ctx context.Context, id uint) ([]model.AddressResponseModel, error) {
	var addresses []entity.Address
	err := u.Where("user_id = ?", id).Find(&addresses).Error
	exception.PanicLogging(err)
	var result []model.AddressResponseModel
	for _, address := range addresses {
		result = append(result, model.AddressResponseModel{
			Description: address.Description,
			PlaceId:     address.PlaceId,
			Id:          address.ID,
		})
	}
	return result, nil
}

func (u *userRepositoryImpl) SaveAddress(ctx context.Context, request model.AddressModel) (interface{}, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("id = ?", request.UserId).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}
	streetNumber := request.Terms[0].Value
	city := ""
	street := request.Terms[1].Value
	postalCode := request.Terms[0].Value
	region := request.Terms[2].Value

	address := entity.Address{
		StreetNumber: streetNumber,
		City:         city,
		Street:       street,
		IsMain:       true,
		PostalCode:   postalCode,
		Region:       region,
		UserId:       request.UserId,
		Description:  request.Description,
		PlaceId:      request.PlaceId,
	}

	err = u.DB.WithContext(ctx).Create(&address).Error
	exception.PanicLogging(err)

	return address, nil
}

func (u *userRepositoryImpl) UpdateProfile(ctx context.Context, request model.UserModel) (model.UserModel, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("username = ?", request.Username).First(&user).Error
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	if user.IsActive == false {
		return model.UserModel{}, errors.New("user is not active")
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.EmailAddress
	//user.PhoneNumber = request.PhoneNumber
	//user.FileName = request.FileName

	err = u.DB.Save(&user).Error
	exception.PanicLogging(err)

	return model.UserModel{
		Id:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}, nil
}

func (u *userRepositoryImpl) SetPassword(id int, password string) (model.UserModel, error) {
	var user entity.User
	err := u.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	exception.PanicLogging(err)
	user.Password = string(hashedPassword)

	err = u.DB.Save(&user).Error
	exception.PanicLogging(err)

	return model.UserModel{
		Id:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (u *userRepositoryImpl) ResetPassword(ctx context.Context, request model.UserModel) (model.UserModel, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("username = ? or phone_number = ?", request.Username, request.PhoneNumber).First(&user).Error
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	if user.IsActive == false {
		return model.UserModel{}, errors.New("user is not active")
	}

	if user.Username == "" {
		return model.UserModel{}, errors.New("user not found")
	}
	return model.UserModel{
		Id:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (u *userRepositoryImpl) ValidateOtp(ctx context.Context, request model.OtpModel) (entity.OneTimePassword, error) {
	var otp entity.OneTimePassword
	err := u.DB.WithContext(ctx).Where("code = ? and user_id = ?", request.Code, request.UserId).First(&otp).Error
	if err != nil {
		return otp, errors.New("invalid OTP")
	}
	if otp.IsUsed {
		return otp, errors.New("OTP already used")
	}
	if otp.ExpiredAt.Before(time.Now()) {
		return otp, errors.New("OTP expired")
	}
	if otp.Code != request.Code {
		return otp, errors.New("Invalid OTP code")
	}
	otp.IsUsed = true
	err = u.DB.Save(&otp).Error
	if err != nil {
		return otp, errors.New("Failed to update OTP status")
	}

	return otp, nil
}

func (u *userRepositoryImpl) FindById(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	exception.PanicLogging(err)
	return user, nil
}

func (u *userRepositoryImpl) List(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := u.DB.WithContext(ctx).Find(&users).Error
	exception.PanicLogging(err)
	return users, nil
}

func (u *userRepositoryImpl) ChangePassword(ctx context.Context, claims map[string]interface{}, passwordModel model.ChangePasswordModel) (entity.User, error) {
	var userResult entity.User
	username := claims["username"].(string)
	result := u.DB.WithContext(ctx).
		Where("(tb_users.username = ?  or tb_users.phone_number = ?)and tb_users.is_active = ?", username, username, true).
		Find(&userResult)
	if result.RowsAffected == 0 {
		return entity.User{}, errors.New("User not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(userResult.Password), []byte(passwordModel.OldPassword))
	if err != nil {
		return entity.User{}, errors.New("Incorrect old password")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordModel.NewPassword), bcrypt.DefaultCost)
	exception.PanicLogging(err)
	userResult.Password = string(hashedPassword)
	err = u.DB.Save(&userResult).Error
	exception.PanicLogging(err)
	return userResult, nil
}

func (u *userRepositoryImpl) Create(model model.UserModel) (entity.User, error) {
	var user entity.User
	err := u.DB.Where("username = ? or phone_number = ? or email = ?", model.PhoneNumber, model.PhoneNumber, model.EmailAddress).Find(&user).Error
	if user.Username != "" {
		return entity.User{}, errors.New("User already exist")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	exception.PanicLogging(err)
	user = entity.User{
		Username:    model.PhoneNumber,
		Password:    string(hashedPassword),
		IsActive:    model.IsActive,
		FirstName:   model.FirstName,
		LastName:    model.LastName,
		Email:       model.EmailAddress,
		UserRole:    model.Role,
		PhoneNumber: model.PhoneNumber,
		FileName:    model.FileName,
		RefineryId:  model.RefineryId,
	}
	var addresses []entity.Address

	addresses = append(addresses, entity.Address{
		City:       model.City,
		Street:     model.Street,
		IsMain:     true,
		PostalCode: model.PostalCode,
		Region:     model.Region,
	})

	user.Addresses = addresses
	err = u.DB.Create(&user).Error
	exception.PanicLogging(err)

	//Todo: Send email/sms notification to user

	return user, nil
}

func (u *userRepositoryImpl) DeleteAll() {
	err := u.DB.Where("1=1").Delete(&entity.User{}).Error
	exception.PanicLogging(err)
}

func (u *userRepositoryImpl) Authentication(ctx context.Context, username string) (entity.User, error) {
	var userResult entity.User
	result := u.DB.WithContext(ctx).
		Where("(tb_users.username = ?  or tb_users.phone_number = ?)and tb_users.is_active = ?", username, username, true).
		Find(&userResult)
	if result.RowsAffected == 0 {
		return entity.User{}, errors.New("user not found")
	}
	return userResult, nil
}

func (u *userRepositoryImpl) SeedUser(ctx context.Context) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	model := model.UserModel{
		Username:     "08011111111",
		Password:     string(hashedPassword),
		FirstName:    "Admin",
		LastName:     "User",
		EmailAddress: "walenewera@gmail.com",
		PhoneNumber:  "08011111111",
		Role:         "admin",
		IsActive:     true,
	}
	var user entity.User
	_ = u.DB.WithContext(ctx).
		Where("(tb_users.username = ?  or tb_users.phone_number = ?) ", model.Username, model.PhoneNumber).
		First(&user).Error
	//exception.PanicLogging(err)

	if user.FirstName == "" {
		_, err := u.Create(model)
		if err == nil {
			common.Logger.Info("Admin user created successfully")
		}
	}

}
func (u *userRepositoryImpl) GetClaimsFromToken(tokenString string) (map[string]interface{}, error) {
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
