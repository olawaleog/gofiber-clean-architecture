package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"io/ioutil"
	"net/http"
)

type RefineryServiceImpl struct {
	repository.RefineryRepository
	service.UserService
}

func NewRefineryServiceImpl(repository *repository.RefineryRepository, userService *service.UserService) service.RefineryService {
	return &RefineryServiceImpl{RefineryRepository: *repository, UserService: *userService}
}

func (r RefineryServiceImpl) GetRefinery(context context.Context, request model.GetRefineryModel) (model.RefineryCostModel, error) {
	var basicCost float64
	refineries, err := r.RefineryRepository.ListRefinery(context)
	var selectRefinery entity.Refinery
	shortestDistance := 100000000000000.0
	exception.PanicLogging(err)
	distanceResult := make(map[string]interface{})
	// calculate route distance for each refinery
	for i := 0; i < len(refineries); i++ {
		distanceResult = CalculateDistance(refineries[i].PlaceId, request.PlaceId)
		distance := distanceResult["distance_km"].(float64)
		if distance < shortestDistance {
			selectRefinery = refineries[i]
			shortestDistance = distance
		}

	}

	if shortestDistance > 40 {
		return model.RefineryCostModel{}, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	//timeInSeconds := distanceResult["time_seconds"].(float64)
	distance := distanceResult["distance_km"].(float64)
	if request.Type == "domestic" {
		basicCost = selectRefinery.DomesticCostPerThousandLitre
	} else {
		basicCost = selectRefinery.IndustrialCostPerThousandLitre
	}

	// calculate cost
	TenThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 10) + (distance * 1) + (1 * 2),
		WaterCost:   basicCost * 10,
		DeliveryFee: distance * 1,
	}

	TwentyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 20) + (distance * 1) + (2 * 2),
		WaterCost:   basicCost * 20,
		DeliveryFee: distance * 1,
	}

	ThirtyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 30) + (distance * 1) + (3 * 2),
		WaterCost:   basicCost * 30,
		DeliveryFee: distance * 1,
	}
	FortyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 40) + (distance * 1) + (4 * 2),
		WaterCost:   basicCost * 40,
		DeliveryFee: distance * 1,
	}

	response := model.RefineryCostModel{
		RefineryId: selectRefinery.ID,
		Address: model.AddressModel{
			Description: selectRefinery.Address,
			PlaceId:     selectRefinery.PlaceId,
		},
		TenThousandLitre:     TenThousand,
		TwentyThousandLitre:  TwentyThousand,
		ThirtyThousandLitre:  ThirtyThousand,
		FortyThousandLitre:   FortyThousand,
		CostPerThousandLitre: basicCost,
		Distance:             distance,
	}

	return response, nil
}

func CalculateDistance(originPlaceID, destinationPlaceID string) map[string]interface{} {
	apiKey := "AIzaSyAFYfTvR_8IzpQb7DHMl9HA6h1kskcz2ok" // Replace with your actual API key
	url := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/distancematrix/json?origins=place_id:%s&destinations=place_id:%s&key=%s",
		originPlaceID, destinationPlaceID, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Failed to call Google Maps API: %v", err))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response body: %v", err))
	}

	var result struct {
		Rows []struct {
			Elements []struct {
				Distance struct {
					Value float64 `json:"value"` // Distance in meters
				} `json:"distance"`
			} `json:"elements"`
		} `json:"rows"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		panic(fmt.Sprintf("Failed to parse JSON response: %v", err))
	}

	if len(result.Rows) > 0 && len(result.Rows[0].Elements) > 0 {
		return map[string]interface{}{
			"distance":     result.Rows[0].Elements[0].Distance.Value,
			"distance_km":  result.Rows[0].Elements[0].Distance.Value / 1000, // Convert meters to kilometers
			"time":         result.Rows[0].Elements[0].Distance.Value / 60,   // Convert meters to minutes
			"time_seconds": result.Rows[0].Elements[0].Distance.Value,        // Convert meters to seconds

		} // Convert meters to kilometers
	}

	panic("No distance data found in API response")
}

func (r RefineryServiceImpl) CreateRefinery(ctx context.Context, refineryModel model.CreateRefineryModel) (model.RefineryModel, error) {
	existingUser := r.UserService.FindByEmailOrPhone(ctx, model.UserModel{
		PhoneNumber:  refineryModel.Phone,
		EmailAddress: refineryModel.Email,
	})
	if existingUser.ID != 0 {
		return model.RefineryModel{}, exception.BadRequestError{
			Message: "User already exist",
		}
	}

	refinery := entity.Refinery{
		Name:                           refineryModel.Name,
		PlaceId:                        refineryModel.PlaceId,
		LicenceExpiry:                  refineryModel.LicenceExpirationDate,
		DomesticCostPerThousandLitre:   refineryModel.DomesticCostPerThousandLitre,
		IndustrialCostPerThousandLitre: refineryModel.IndustrialCostPerThousandLitre,
		RawLocationData:                refineryModel.RawLocationData,
		Region:                         refineryModel.Region,
		Phone:                          refineryModel.Phone,
		Email:                          refineryModel.Email,
		Website:                        refineryModel.Website,
		Address:                        refineryModel.Address,
	}

	refineryData, err := r.RefineryRepository.Create(ctx, refinery)
	if err != nil {
		return model.RefineryModel{}, err
	}
	refineryModelResponse := model.RefineryModel{
		Id:                             refineryData.ID,
		Name:                           refineryData.Name,
		PlaceId:                        refineryData.PlaceId,
		LicenceExpiry:                  refineryData.LicenceExpiry,
		DomesticCostPerThousandLitre:   refineryData.DomesticCostPerThousandLitre,
		IndustrialCostPerThousandLitre: refineryData.IndustrialCostPerThousandLitre,
		RawLocationData:                refineryData.RawLocationData,
		Region:                         refineryData.Region,
		Phone:                          refineryData.Phone,
		Email:                          refineryData.Email,
		Website:                        refineryData.Website,
		Address:                        refineryData.Address,
	}
	password, err := common.GeneratePassword(8)
	exception.PanicLogging(err)
	user := model.UserModel{
		Username:     refineryModel.Phone,
		Password:     password,
		Role:         common.REFINERY_ADMIN_ROLE,
		EmailAddress: refineryModel.Email,
		FirstName:    refineryModel.FirstName,
		LastName:     refineryModel.LastName,
		PhoneNumber:  refineryModel.Phone,
		IsActive:     false,
		RefineryId:   refineryData.ID,
	}

	_ = r.UserService.Register(context.TODO(), user)
	exception.PanicLogging(err)
	return refineryModelResponse, nil
}

func (r RefineryServiceImpl) ListRefineries(ctx context.Context) ([]model.RefineryModel, error) {
	var result []model.RefineryModel
	refineries, err := r.RefineryRepository.ListRefinery(ctx)
	exception.PanicLogging(err)
	for _, refinery := range refineries {
		result = append(result, model.RefineryModel{
			Id:                             refinery.ID,
			Name:                           refinery.Name,
			PlaceId:                        refinery.PlaceId,
			LicenceExpiry:                  refinery.LicenceExpiry,
			Address:                        refinery.Address,
			Email:                          refinery.Email,
			Phone:                          refinery.Phone,
			Licence:                        refinery.Licence,
			IndustrialCostPerThousandLitre: refinery.IndustrialCostPerThousandLitre,
			DomesticCostPerThousandLitre:   refinery.DomesticCostPerThousandLitre,
			Region:                         refinery.Region,
			Website:                        refinery.Website,
		})
	}

	return result, nil

}

func (r RefineryServiceImpl) UpdateRefinery(ctx context.Context, refineryModel model.CreateRefineryModel, id string) (model.RefineryModel, error) {
	refinery := entity.Refinery{
		Name:                           refineryModel.Name,
		PlaceId:                        refineryModel.PlaceId,
		LicenceExpiry:                  refineryModel.LicenceExpirationDate,
		DomesticCostPerThousandLitre:   refineryModel.DomesticCostPerThousandLitre,
		IndustrialCostPerThousandLitre: refineryModel.IndustrialCostPerThousandLitre,
		RawLocationData:                refineryModel.RawLocationData,
		Region:                         refineryModel.Region,
		Phone:                          refineryModel.Phone,
		Email:                          refineryModel.Email,
		Website:                        refineryModel.Website,
		Address:                        refineryModel.Address,
	}

	refineryData, err := r.RefineryRepository.Update(ctx, refinery, id)
	if err != nil {
		return model.RefineryModel{}, err
	}
	refineryModelResponse := model.RefineryModel{
		Id:                             refineryData.ID,
		Name:                           refineryData.Name,
		LicenceExpiry:                  refineryData.LicenceExpiry,
		DomesticCostPerThousandLitre:   refineryData.DomesticCostPerThousandLitre,
		IndustrialCostPerThousandLitre: refineryData.IndustrialCostPerThousandLitre,
	}
	return refineryModelResponse, nil
}

func (r RefineryServiceImpl) GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error) {
	refineryData, err := r.RefineryRepository.GetRefineryDashboardData(ctx, u)
	exception.PanicLogging(err)
	if refineryData == nil {
		return nil, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}

	return refineryData, nil
}
