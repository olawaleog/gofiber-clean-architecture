package impl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/go-playground/validator/v10"
)

type SettingServiceImpl struct {
	SettingRepository repository.SettingRepository
	RedisService      service.RedisService
	Validate          *validator.Validate
}

func NewSettingService(settingRepository repository.SettingRepository, redisService service.RedisService, validate *validator.Validate) service.SettingService {
	return &SettingServiceImpl{
		SettingRepository: settingRepository,
		RedisService:      redisService,
		Validate:          validate,
	}
}

func (service *SettingServiceImpl) Create(ctx context.Context, request model.SettingModel) (entity.Setting, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return entity.Setting{}, err
	}

	return service.SettingRepository.Create(ctx, request)
}

func (service *SettingServiceImpl) Update(ctx context.Context, key string, request model.SettingModel) (entity.Setting, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return entity.Setting{}, err
	}

	return service.SettingRepository.Update(ctx, key, request)
}

func (service *SettingServiceImpl) List(ctx context.Context) ([]model.SettingResponseModel, error) {
	settings, err := service.SettingRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	var settingResponses []model.SettingResponseModel
	for _, setting := range settings {
		settingResponses = append(settingResponses, model.SettingResponseModel{
			ID:    setting.ID,
			Key:   setting.Key,
			Value: setting.Value,
		})
	}

	return settingResponses, nil
}

func (service *SettingServiceImpl) FindByKey(ctx context.Context, key string) (model.SettingResponseModel, error) {
	// Try to get setting from Redis first
	redisKey := "setting_" + key
	cachedData, err := service.RedisService.Get(redisKey)

	// If found in Redis, unmarshal and return
	if err == nil {
		var settingResponse model.SettingResponseModel
		if err := json.Unmarshal([]byte(cachedData), &settingResponse); err == nil {
			return settingResponse, nil
		}
		// If we can't unmarshal, continue to get from DB
	}

	// If not found in Redis or unmarshalling failed, get from DB
	setting, err := service.SettingRepository.FindByKey(ctx, key)
	if err != nil {
		// Not found in DB either, return error
		return model.SettingResponseModel{}, err
	}

	// Found in DB, prepare response
	settingResponse := model.SettingResponseModel{
		ID:    setting.ID,
		Key:   setting.Key,
		Value: setting.Value,
	}

	// Cache in Redis before returning
	// Convert to JSON
	jsonData, err := json.Marshal(settingResponse)
	if err == nil {
		// Store in Redis with a TTL of 1 hour
		err := service.RedisService.Set(redisKey, jsonData, time.Hour)
		if err != nil {
			return model.SettingResponseModel{}, err
		}
	}

	return settingResponse, nil
}
