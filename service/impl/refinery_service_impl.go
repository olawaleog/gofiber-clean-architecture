package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type RefineryServiceImpl struct {
	repository.RefineryRepository
	service.UserService
	service.MessageService
	configuration.Config
}

func NewRefineryServiceImpl(repository *repository.RefineryRepository, userService *service.UserService, messageService *service.MessageService, configuration configuration.Config) service.RefineryService {
	return &RefineryServiceImpl{RefineryRepository: *repository, UserService: *userService, MessageService: *messageService, Config: configuration}
}

func (r RefineryServiceImpl) GetRefinery(context context.Context, request model.GetRefineryModel) (model.RefineryCostModel, error) {
	var basicCost float64
	refineries, err := r.RefineryRepository.ListRefinery(context)
	var selectRefinery entity.Refinery = refineries[0]
	shortestDistance := 60.0
	exception.PanicLogging(err)
	distanceResult := make(map[string]interface{})
	// calculate route distance for each refinery
	for i := 0; i < len(refineries); i++ {
		if !refineries[i].HasIndustrialWaterSupply && request.Type == "industrial" {
			continue
		}
		if !refineries[i].HasDomesticWaterSupply && request.Type == "domestic" {
			continue
		}
		distanceResult = r.CalculateDistance(refineries[i], request.PlaceId)
		if len(distanceResult) == 0 {
			continue
		}
		distance := distanceResult["distance_km"].(float64)
		if distance < shortestDistance {
			selectRefinery = refineries[i]
			shortestDistance = distance
		} else {
			continue
		}

	}

	if shortestDistance > 100 {
		return model.RefineryCostModel{}, nil
	}
	//timeInSeconds := distanceResult["time_seconds"].(float64)
	//distance := distanceResult["distance_km"].(float64)
	if request.Type == "domestic" {
		basicCost = selectRefinery.DomesticCostPerThousandLitre
	} else {
		basicCost = selectRefinery.IndustrialCostPerThousandLitre
	}

	// calculate cost
	TenThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 10) + (shortestDistance * 1) + (1 * 2),
		WaterCost:   basicCost * 10,
		DeliveryFee: shortestDistance * 1,
	}

	TwentyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 20) + (shortestDistance * 1) + (2 * 2),
		WaterCost:   basicCost * 20,
		DeliveryFee: shortestDistance * 1,
	}

	ThirtyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 30) + (shortestDistance * 1) + (3 * 2),
		WaterCost:   basicCost * 30,
		DeliveryFee: shortestDistance * 1,
	}
	FortyThousand := model.WaterCostModel{
		TotalCost:   (basicCost * 40) + (shortestDistance * 1) + (4 * 2),
		WaterCost:   basicCost * 40,
		DeliveryFee: shortestDistance * 1,
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
		Distance:             shortestDistance,
		Currency:             selectRefinery.Currency,
	}

	return response, nil
}

func (r *RefineryServiceImpl) CalculateDistance(origin entity.Refinery, destinationPlaceID string) map[string]interface{} {
	apiKey := r.Config.Get("GOOGLE_MAPS_API_KEY") // Replace with your actual API key
	url := ""
	if origin.PlaceId == "" {
		url = fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s,%s&destinations=place_id:%s&key=%s",
			origin.Longitude, origin.Latitude, destinationPlaceID, apiKey,
		)
	} else {
		url = fmt.Sprintf(
			"https://maps.googleapis.com/maps/api/distancematrix/json?origins=place_id:%s&destinations=place_id:%s&key=%s",
			origin.PlaceId, destinationPlaceID, apiKey,
		)
	}

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
	return map[string]interface{}{}
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
		Longitude:                      refineryModel.Longitude,
		Latitude:                       refineryModel.Latitude,
		HasDomesticWaterSupply:         refineryModel.HasDomesticWaterSupply,
		HasIndustrialWaterSupply:       refineryModel.HasIndustrialWaterSupply,
		IsActive:                       true,
		Country:                        refineryModel.Country,
		Currency:                       refineryModel.Currency,
		CurrencySymbol:                 refineryModel.CurrencySymbol,
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
		IsActive:                       refineryData.IsActive,
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
		CountryCode:  "+233",
	}

	_ = r.UserService.Register(context.TODO(), user)
	exception.PanicLogging(err)
	smsModel := model.SMSMessageModel{
		Message:     "Hello " + user.FirstName + "," + "Your Refinery has been registered successfully.Your password is " + user.Password + ":\n visit https://aqua-wizz.app\n",
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
	}
	r.MessageService.SendSMS(context.TODO(), smsModel)
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
			IsActive:                       refinery.IsActive,
			Longitude:                      refinery.Longitude,
			Latitude:                       refinery.Latitude,
			HasDomesticWaterSupply:         refinery.HasDomesticWaterSupply,
			HasIndustrialWaterSupply:       refinery.HasIndustrialWaterSupply,
		})
	}

	return result, nil

}

func (r RefineryServiceImpl) UpdateRefinery(ctx context.Context, refineryModel model.CreateRefineryModel, id string) (model.RefineryModel, error) {
	existingRefinery := r.RefineryRepository.FindById(ctx, id)
	if existingRefinery.ID == 0 {
		return model.RefineryModel{}, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	existingRefinery.Name = refineryModel.Name
	existingRefinery.PlaceId = refineryModel.PlaceId
	existingRefinery.LicenceExpiry = refineryModel.LicenceExpirationDate
	existingRefinery.DomesticCostPerThousandLitre = refineryModel.DomesticCostPerThousandLitre
	existingRefinery.IndustrialCostPerThousandLitre = refineryModel.IndustrialCostPerThousandLitre
	existingRefinery.RawLocationData = refineryModel.RawLocationData
	existingRefinery.Region = refineryModel.Region
	existingRefinery.Phone = refineryModel.Phone
	existingRefinery.Email = refineryModel.Email
	existingRefinery.Website = refineryModel.Website
	existingRefinery.Address = refineryModel.Address
	existingRefinery.HasDomesticWaterSupply = refineryModel.HasDomesticWaterSupply
	existingRefinery.HasIndustrialWaterSupply = refineryModel.HasIndustrialWaterSupply
	existingRefinery.Longitude = refineryModel.Longitude
	existingRefinery.Latitude = refineryModel.Latitude

	refineryData, err := r.RefineryRepository.Update(ctx, existingRefinery, id)
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

func (r RefineryServiceImpl) ToggleRefineryStatus(ctx context.Context, statusModel model.ToggleRefineryStatusModel) (bool, error) {
	refinery := r.RefineryRepository.FindById(ctx, statusModel.Id)
	if refinery.ID == 0 {
		return false, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	refinery.IsActive = statusModel.Status
	refinery, _ = r.RefineryRepository.Update(ctx, refinery, fmt.Sprintf("%d", statusModel.Id))
	if refinery.ID == 0 {
		return false, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	return true, nil

}
