package repository

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type UserRepository interface {
	Authentication(ctx context.Context, username string) (entity.User, error)
	ChangePassword(ctx context.Context, claims map[string]interface{}, passwordModel model.ChangePasswordModel) (entity.User, error)
	Create(model model.UserModel) (entity.User, error)
	DeleteAll()
	SeedUser(ctx context.Context)
	List(ctx context.Context) ([]entity.User, error)
	FindById(ctx context.Context, id int) (entity.User, error)
	ValidateOtp(ctx context.Context, request model.OtpModel) (entity.OneTimePassword, error)
	ResetPassword(ctx context.Context, request model.UserModel) (model.UserModel, error)
	SetPassword(id int, password string) (model.UserModel, error)
	UpdateProfile(ctx context.Context, request model.UserModel) (model.UserModel, error)
	SaveAddress(ctx context.Context, request map[string]interface{}) (entity.Address, error)
	FindAllAddress(ctx context.Context, id uint) ([]model.AddressResponseModel, error)
	FindByEmailOrPhone(ctx context.Context, userModel model.UserModel) (entity.User, error)
	UpdateFcmToken(ctx context.Context, request model.UpdateFcmToken) error
	FineAddressById(ctx context.Context, id uint) (entity.Address, error)
	FindAllWithFcmToken(ctx context.Context) ([]entity.User, error)
}
