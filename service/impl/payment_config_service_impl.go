package impl

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type PaymentConfigServiceImpl struct {
	repository repository.PaymentConfigRepository
}

func NewPaymentConfigService(repo repository.PaymentConfigRepository) service.PaymentConfigService {
	return &PaymentConfigServiceImpl{
		repository: repo,
	}
}

func (s *PaymentConfigServiceImpl) GetPaymentConfig(ctx context.Context, countryCode string) (*model.PaymentConfigModel, error) {
	config, err := s.repository.GetPaymentConfig(ctx, countryCode)
	if err != nil {
		return nil, err
	}
	return s.entityToModel(config), nil
}

func (s *PaymentConfigServiceImpl) ListPaymentConfigs(ctx context.Context) ([]model.PaymentConfigModel, error) {
	configs, err := s.repository.ListPaymentConfigs(ctx)
	if err != nil {
		return nil, err
	}

	var result []model.PaymentConfigModel
	for _, config := range configs {
		result = append(result, *s.entityToModel(&config))
	}
	return result, nil
}

func (s *PaymentConfigServiceImpl) CreatePaymentConfig(ctx context.Context, config *model.PaymentConfigModel) error {
	entity := s.modelToEntity(config)
	return s.repository.SavePaymentConfig(ctx, entity)
}

func (s *PaymentConfigServiceImpl) UpdatePaymentConfig(ctx context.Context, config *model.PaymentConfigModel) error {
	entity := s.modelToEntity(config)
	return s.repository.UpdatePaymentConfig(ctx, entity)
}

func (s *PaymentConfigServiceImpl) DeletePaymentConfig(ctx context.Context, countryCode string) error {
	return s.repository.DeletePaymentConfig(ctx, countryCode)
}

func (s *PaymentConfigServiceImpl) entityToModel(entity *entity.PaymentConfiguration) *model.PaymentConfigModel {
	return &model.PaymentConfigModel{
		ID:          entity.ID,
		CountryCode: entity.CountryCode,
		SecretKey:   entity.SecretKey,
		PublicKey:   entity.PublicKey,
	}
}

func (s *PaymentConfigServiceImpl) RefreshCache(ctx context.Context) error {
	configs, err := s.repository.ListPaymentConfigs(ctx)
	if err != nil {
		return err
	}

	// Force refresh by getting each config again
	for _, config := range configs {
		_, err := s.repository.GetPaymentConfig(ctx, config.CountryCode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PaymentConfigServiceImpl) modelToEntity(model *model.PaymentConfigModel) *entity.PaymentConfiguration {
	return &entity.PaymentConfiguration{
		CountryCode: model.CountryCode,
		SecretKey:   model.SecretKey,
		PublicKey:   model.PublicKey,
	}
}
