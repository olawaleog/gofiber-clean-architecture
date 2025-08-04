package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type UserService interface {
	Authentication(ctx context.Context, model model.LoginModel) entity.User
	ChangePassword(ctx context.Context, token string, model model.ChangePasswordModel) entity.User
	GetClaimsFromToken(ctx context.Context, tokenString string) (map[string]interface{}, error)
	Register(ctx context.Context, model model.UserModel) entity.User
	SeedUser(ctx context.Context)
	List(ctx context.Context) ([]entity.User, error)
	FindByID(ctx context.Context, id int) (model.UserModel, error)
	RegisterCustomer(ctx context.Context, request model.UserModel) interface{}
	ValidateOtp(ctx context.Context, request model.OtpModel) entity.OneTimePassword
	ResetPassword(ctx context.Context, request model.UserModel) model.UserModel
	UpdateUserPassword(ctx context.Context, request model.ResetPasswordViewModel) model.UserModel
	UpdateProfile(ctx context.Context, request model.UserModel, token string) (model.UserModel, error)
	SaveAddress(ctx context.Context, request model.AddressModel) (interface{}, error)
	GetAddresses(ctx context.Context, id uint) (interface{}, error)
	FindByEmailOrPhone(ctx context.Context, userModel model.UserModel) entity.User
	UpdateFcmToken(ctx context.Context, request model.UpdateFcmToken) error
}
