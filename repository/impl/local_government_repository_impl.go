package impl

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"gorm.io/gorm"
)

type localGovernemntRepositoryImpl struct {
	*gorm.DB
}

func NewLocalGovernmentRepository(db *gorm.DB) repository.LocalGovernmentRepository {
	return localGovernemntRepositoryImpl{db}
}

func (l localGovernemntRepositoryImpl) FindAll(ctx context.Context) ([]model.LocalGovernmentModel, error) {
	var localGovernments []model.LocalGovernmentModel
	err := l.DB.WithContext(ctx).Find(&localGovernments)
	exception.PanicLogging(err)
	return localGovernments, nil
}

func (l localGovernemntRepositoryImpl) ToggleLocalGovernmentActive(ctx context.Context, id string) error {
	var localGovernment model.LocalGovernmentModel
	err := l.DB.WithContext(ctx).Where("id = ?", id).First(&localGovernment)
	exception.PanicLogging(err)
	localGovernment.IsActive = !localGovernment.IsActive
	err = l.DB.WithContext(ctx).Save(&localGovernment)
	exception.PanicLogging(err)
	return nil
}
