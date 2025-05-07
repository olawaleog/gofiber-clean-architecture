package repository

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type LocalGovernmentRepository interface {
	FindAll(ctx context.Context) ([]model.LocalGovernmentModel, error)
	ToggleLocalGovernmentActive(ctx context.Context, id string) error
}
