package service

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
)

type LocalGovernmentService interface {
	FindAll(ctx context.Context) ([]model.LocalGovernmentModel, error)
	ToggleLocalGovernmentActive(ctx context.Context, id string) error
	GetPlaceSuggestion(ctx context.Context, placeString string) interface{}
}
