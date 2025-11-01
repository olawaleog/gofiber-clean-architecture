package impl

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserRepositoryImpl(DB *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{DB: DB}
}

type userRepositoryImpl struct {
	*gorm.DB
}

func (u userRepositoryImpl) Update(ctx context.Context, userData entity.User) (entity.User, error) {
	err := u.DB.WithContext(ctx).Updates(userData).Error
	if err != nil {
		return userData, err
	}
	return userData, nil
}
func (u *userRepositoryImpl) FineAddressById(ctx context.Context, id uint) (entity.Address, error) {
	var address entity.Address
	err := u.DB.WithContext(ctx).Where("id = ?", id).First(&address).Error
	if err != nil {
		return entity.Address{}, errors.New("address not found")
	}
	return address, nil
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
			Longitude:   toFloat64(address.Longitude),
			Latitude:    toFloat64(address.Latitude),
			CountryCode: address.CountryCode,
		})
	}
	return result, nil
}

func (u *userRepositoryImpl) SaveAddress(ctx context.Context, request map[string]interface{}) (entity.Address, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("id = ?", request["userId"]).First(&user).Error
	if err != nil {
		return entity.Address{}, errors.New("user not found")
	}

	// Check for existing address by PlaceId or Description
	var existingAddress entity.Address
	err = u.DB.WithContext(ctx).Where("user_id = ? AND (place_id = ? OR description = ?)", request["userId"], request["placeId"], request["description"]).First(&existingAddress).Error
	if err == nil {
		return existingAddress, nil
	} else if err != gorm.ErrRecordNotFound {
		return entity.Address{}, err
	}

	streetNumber := "" //request.Terms[0].Value
	city := ""
	street := "  "   //request.Terms[1].Value
	postalCode := "" //r equest.Terms[0].Value
	region := ""     // request.Terms[2].Value

	address := entity.Address{
		StreetNumber: streetNumber,
		City:         city,
		Street:       street,
		IsMain:       true,
		PostalCode:   postalCode,
		Region:       region,
		UserId:       request["userId"].(uint),
		Description:  request["description"].(string),
		PlaceId:      request["place_id"].(string),
		Longitude:    strconv.FormatFloat(request["longitude"].(float64), 'f', -1, 64),
		Latitude:     strconv.FormatFloat(request["latitude"].(float64), 'f', -1, 64),
		CountryCode:  request["country_code"].(string),
	}

	err = u.DB.WithContext(ctx).Create(&address).Error
	exception.PanicLogging(err)

	return address, nil
}

func (u *userRepositoryImpl) UpdateProfile(ctx context.Context, request model.UserModel) (model.UserModel, error) {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("id = ?", request.Id).First(&user).Error
	if err != nil {
		return model.UserModel{}, errors.New("user not found")
	}

	if user.IsActive == false {
		return model.UserModel{}, errors.New("user is not active")
	}

	// Check if email is being updated and already exists for another user
	//if user.Email != request.EmailAddress {
	//	var existingUser entity.User
	//	err := u.DB.WithContext(ctx).Where("email = ? AND id != ?", request.EmailAddress, user.ID).First(&existingUser).Error
	//	if err == nil {
	//		return model.UserModel{}, errors.New("email address already exists for another user")
	//	} else if err != gorm.ErrRecordNotFound {
	//		return model.UserModel{}, fmt.Errorf("error checking email uniqueness: %w", err)
	//	}
	user.Email = request.EmailAddress
	//}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Region = request.Region
	user.CountryCode = request.CountryCode

	err = u.DB.WithContext(ctx).Updates(user).Error
	if err != nil {
		return model.UserModel{}, fmt.Errorf("failed to update user profile: %w", err)
	}

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
	user.IsActive = true

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
		return otp, errors.New("invalid OTP code")
	}
	otp.IsUsed = true
	err = u.DB.Save(&otp).Error
	if err != nil {
		return otp, errors.New("failed to update OTP status")
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

	// First, check if user exists with the same phone number or username
	err := u.DB.Where("username = ? or phone_number = ?", model.PhoneNumber, model.PhoneNumber).Find(&user).Error
	if user.Username != "" {
		return entity.User{}, errors.New("User with phone number already exists")
	}

	// Check if user exists with the same email address
	if model.EmailAddress != "" {
		var userWithEmail entity.User
		err := u.DB.Where("email = ?", model.EmailAddress).Find(&userWithEmail).Error
		if err == nil && userWithEmail.ID != 0 {
			return entity.User{}, errors.New("User with email address already exists")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	exception.PanicLogging(err)
	user = entity.User{
		Username:       model.PhoneNumber,
		Password:       string(hashedPassword),
		IsActive:       model.IsActive,
		FirstName:      model.FirstName,
		LastName:       model.LastName,
		Email:          model.EmailAddress,
		UserRole:       model.Role,
		PhoneNumber:    model.PhoneNumber,
		FileName:       model.FileName,
		RefineryId:     model.RefineryId,
		AreaCode:       model.AreaCode,
		CountryCode:    model.CountryCode,
		EmailValidated: false,
	}
	var addresses []entity.Address

	if model.Street != "" {
		addresses = append(addresses, entity.Address{
			City:       model.City,
			Street:     model.Street,
			IsMain:     false,
			PostalCode: model.PostalCode,
			Region:     model.Region,
		})
		user.Addresses = addresses
	}

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

	//parse username to int
	var _, ok = strconv.ParseInt(username, 10, 64)
	if ok == nil && username[0] != '0' {
		username = "0" + username
	}
	result := u.DB.WithContext(ctx).
		Where("(tb_users.username = ?  or tb_users.phone_number = ? or tb_users.email = ? )", username, username, username).
		Find(&userResult)
	if result.RowsAffected == 0 {
		return entity.User{}, errors.New("user not found")
	}
	return userResult, nil
}

func (u *userRepositoryImpl) SeedUser(ctx context.Context) {
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	//if err != nil {
	//	return
	//}
	adminModel := model.UserModel{
		Username:     "08011111111",
		Password:     "admin",
		FirstName:    "Admin",
		LastName:     "User",
		EmailAddress: "walenewera@gmail.com",
		PhoneNumber:  "08011111111",
		Role:         "admin",
		IsActive:     true,
	}
	var user entity.User
	_ = u.DB.WithContext(ctx).
		Where("(tb_users.username = ?  or tb_users.phone_number = ?) ", adminModel.Username, adminModel.PhoneNumber).
		First(&user).Error
	//exception.PanicLogging(err)

	if user.FirstName == "" {
		_, err := u.Create(adminModel)
		if err == nil {
			logger.Logger.Info("Admin user created successfully")
		}
	}

	// Seed refinery_admin user attached to the only active refinery
	var activeRefinery entity.Refinery
	err := u.DB.WithContext(ctx).Where("is_active = ?", true).First(&activeRefinery).Error
	if err == nil && activeRefinery.ID != 0 {
		refineryAdminModel := model.UserModel{
			Username:     "08099999999",
			Password:     "Vnp-1234",
			FirstName:    "Refinery",
			LastName:     "Admin",
			EmailAddress: "refineryadmin@example.com",
			PhoneNumber:  "08099999999",
			Role:         "refinery_admin",
			IsActive:     true,
			RefineryId:   activeRefinery.ID,
		}
		var refineryAdmin entity.User
		_ = u.DB.WithContext(ctx).
			Where("(tb_users.username = ?  or tb_users.phone_number = ?)", refineryAdminModel.Username, refineryAdminModel.PhoneNumber).
			First(&refineryAdmin).Error
		if refineryAdmin.FirstName == "" {
			_, err := u.Create(refineryAdminModel)
			if err == nil {
				logger.Logger.Info("Refinery admin user created successfully")
			}
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

func (u *userRepositoryImpl) UpdateFcmToken(ctx context.Context, request model.UpdateFcmToken) error {
	var user entity.User
	err := u.DB.WithContext(ctx).Where("id = ?", request.Id).First(&user).Error
	if err != nil {
		return err
	}
	user.FcmToken = request.FcmToken
	err = u.DB.Save(&user).Error

	return err
}
func (u *userRepositoryImpl) FindAllWithFcmToken(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := u.DB.WithContext(ctx).Where("fcm_token IS NOT NULL AND fcm_token != ''").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func toFloat64(value string) float64 {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0 // or handle the error as needed
	}
	return floatValue
}
