package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"io/ioutil"
	"net/http"
)

type localGovernmentServiceImpl struct {
	repository.LocalGovernmentRepository
}

func NewLocalGovernmentServiceImpl(r *repository.LocalGovernmentRepository) service.LocalGovernmentService {
	return &localGovernmentServiceImpl{LocalGovernmentRepository: *r}
}

func (l localGovernmentServiceImpl) FindAll(ctx context.Context) ([]model.LocalGovernmentModel, error) {
	var localGovernments []model.LocalGovernmentModel
	localGovernments, err := l.LocalGovernmentRepository.FindAll(ctx)
	exception.PanicLogging(err)
	return localGovernments, nil
}

func (l localGovernmentServiceImpl) ToggleLocalGovernmentActive(ctx context.Context, id string) error {
	err := l.LocalGovernmentRepository.ToggleLocalGovernmentActive(ctx, id)
	exception.PanicLogging(err)
	return nil
}

func (l localGovernmentServiceImpl) GetPlaceSuggestion(ctx context.Context, placeString string) interface{} {
	googleUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/autocomplete/json?input=%v&key=AIzaSyAFYfTvR_8IzpQb7DHMl9HA6h1kskcz2ok&language=en", placeString)
	response, err := http.Get(googleUrl)
	exception.PanicLogging(err)
	responseData, err := ioutil.ReadAll(response.Body)
	data := string(responseData)
	//desrialize json
	var result map[string]interface{}
	err = json.Unmarshal([]byte(data), &result)
	exception.PanicLogging(err)
	return result
}
