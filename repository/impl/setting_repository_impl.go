package impl

import (
	"context"
	"errors"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

func NewSettingRepository(DB *gorm.DB) repository.SettingRepository {
	return &settingRepositoryImpl{DB: DB}
}

type settingRepositoryImpl struct {
	*gorm.DB
}

func (s *settingRepositoryImpl) Create(ctx context.Context, settingModel model.SettingModel) (entity.Setting, error) {
	var existingSetting entity.Setting
	err := s.DB.WithContext(ctx).Where("key = ?", settingModel.Key).First(&existingSetting).Error
	if err == nil {
		return entity.Setting{}, errors.New("setting with this key already exists")
	}

	setting := entity.Setting{
		Key:   settingModel.Key,
		Value: settingModel.Value,
	}

	err = s.DB.WithContext(ctx).Create(&setting).Error
	if err != nil {
		return entity.Setting{}, err
	}

	return setting, nil
}

func (s *settingRepositoryImpl) Update(ctx context.Context, key string, settingModel model.SettingModel) (entity.Setting, error) {
	var setting entity.Setting
	err := s.DB.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Setting{}, errors.New("setting not found")
		}
		return entity.Setting{}, err
	}

	// Update values
	setting.Value = settingModel.Value
	if settingModel.Key != key {
		// Check if new key already exists (only if key is being changed)
		var existingSetting entity.Setting
		err = s.DB.WithContext(ctx).Where("key = ?", settingModel.Key).First(&existingSetting).Error
		if err == nil {
			return entity.Setting{}, errors.New("setting with new key already exists")
		}
		setting.Key = settingModel.Key
	}

	err = s.DB.WithContext(ctx).Save(&setting).Error
	if err != nil {
		return entity.Setting{}, err
	}

	return setting, nil
}

func (s *settingRepositoryImpl) List(ctx context.Context) ([]entity.Setting, error) {
	var settings []entity.Setting
	err := s.DB.WithContext(ctx).Find(&settings).Error
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *settingRepositoryImpl) FindByKey(ctx context.Context, key string) (entity.Setting, error) {
	var setting entity.Setting
	err := s.DB.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Setting{}, errors.New("setting not found")
		}
		return entity.Setting{}, err
	}

	return setting, nil
}
