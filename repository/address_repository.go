package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
)

type AddressRepository interface {
	FindById(ctx context.Context, id int, userId int) (entity.Address, error)
	FindAll(ctx context.Context, userId int) ([]entity.Address, error)
}
