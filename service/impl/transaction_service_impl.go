package impl

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type transactionServiceImpl struct {
	repository.TransactionRepository
	repository.OrderRepository
	repository.PaymentMethodRepository
	repository.PaymentConfigRepository
	service.HttpService
	configuration.Config
	service.NotificationService // Add this line
}

// Update your constructor to accept NotificationService
func NewTransactionServiceImpl(
	transactionRepository *repository.TransactionRepository,
	orderRepo *repository.OrderRepository,
	paymentRepo *repository.PaymentMethodRepository,
	client *service.HttpService,
	config configuration.Config,
	notificationService *service.NotificationService, // Add this param
	paymentConfigRepository *repository.PaymentConfigRepository,
) service.TransactionService {
	return &transactionServiceImpl{
		TransactionRepository:   *transactionRepository,
		HttpService:             *client,
		Config:                  config,
		OrderRepository:         *orderRepo,
		PaymentMethodRepository: *paymentRepo,
		NotificationService:     *notificationService,
		PaymentConfigRepository: *paymentConfigRepository,
	}
}

func (t *transactionServiceImpl) ProcessRecurringPayment(ctx context.Context, request model.MobileMoneyRequestModel) interface{} {
	transaction := t.insertTransaction(ctx, request)

	paymentMethod, err := t.PaymentMethodRepository.GetByID(ctx, strconv.FormatUint(uint64(request.PaymentMethodId), 10))
	exception.PanicLogging(err)

	paymentConfig, err := t.PaymentConfigRepository.GetPaymentConfig(ctx, request.CountryCode)
	exception.PanicLogging(err)
	if paymentConfig == nil {
		exception.PanicLogging(errors.New("payment configuration not found for the given country code"))
	}

	var email string
	if request.EmailAddress == "" {
		email = "Intelblue28@gmail.com"
	} else {
		email = request.EmailAddress
	}

	data := make(map[string]interface{})
	data["authorization_code"] = paymentMethod.AuthCode
	data["amount"] = int(request.Amount * 100)
	data["email"] = email
	data["currency"] = request.Currency
	paystackUrl := t.Config.Get("PAYSTACK_BASE_URL") + "/transaction/charge_authorization"
	response, err := t.SendToPayStack(ctx, paystackUrl, data, paymentConfig.SecretKey)
	exception.PanicLogging(err)
	transaction.Reference = response["data"].(map[string]interface{})["reference"].(string)
	err = t.TransactionRepository.Update(ctx, transaction)
	var trxData = response["data"].(map[string]interface{})
	trxData["transaction_id"] = transaction.ID
	response["data"] = trxData
	return response
}

func (t *transactionServiceImpl) GetCustomerOrders(ctx context.Context, u uint) ([]model.OrderModel, error) {
	orders, err := t.OrderRepository.GetUserOrders(ctx, u)
	exception.PanicLogging(err)
	if orders == nil {
		return nil, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	var orderModels []model.OrderModel
	for _, order := range orders {
		longitude, _ := strconv.ParseFloat(order.Transaction.Address.Longitude, 64)
		latitude, _ := strconv.ParseFloat(order.Transaction.Address.Latitude, 64)
		orderModels = append(orderModels, model.OrderModel{
			Id:          order.ID,
			Amount:      order.Amount,
			Currency:    order.Currency,
			WaterCost:   order.WaterCost,
			DeliveryFee: order.DeliveryFee,
			RefineryId:  order.RefineryId,
			Capacity:    order.Capacity,
			Type:        order.Type,
			Transaction: model.TransactionModel{
				ID:          order.Transaction.ID,
				Email:       order.Transaction.Email,
				PhoneNumber: order.Transaction.PhoneNumber,
				Amount:      order.Transaction.Amount,
				Currency:    order.Transaction.Currency,
				PaymentID:   order.Transaction.PaymentID,
				Provider:    order.Transaction.Provider,
				PaymentType: order.Transaction.PaymentType,
				Status:      order.Transaction.Status,
				//RawRequest:  order.Transaction.RawRequest,
				//RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
				DeliveryAddress: model.AddressModel{
					Longitude:   longitude,
					Latitude:    latitude,
					Description: order.Transaction.Address.Description,
					PlaceId:     order.Transaction.Address.PlaceId,
				},
			},
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
			DeliveryAddress: order.Transaction.Address.Description,
			//UserId:       u,
			//UserName:     "",
			//UserPhoneNumber:"",
			//UserEmailAddress:"",
			//UserFirstName:"",
			//UserLastName:"",
		})
	}

	return orderModels, nil
}

func (t *transactionServiceImpl) GetAdminDashboardData(ctx context.Context) (map[string]interface{}, error) {

	response := t.TransactionRepository.GetAdminDashboardData(ctx)

	return response, nil

}

func (t *transactionServiceImpl) GetTransactions(ctx context.Context) []model.TransactionModel {
	var list []model.TransactionModel
	var transactions []entity.Transaction
	transactions = t.TransactionRepository.FindAll(ctx)
	for _, transaction := range transactions {
		list = append(list, model.TransactionModel{
			ID:          transaction.ID,
			Email:       transaction.Email,
			PhoneNumber: transaction.PhoneNumber,
			Amount:      transaction.Amount,
			Currency:    transaction.Currency,
			PaymentID:   transaction.PaymentID,
			Provider:    transaction.Provider,
			PaymentType: transaction.PaymentType,
			Status:      transaction.Status,
			//RawRequest:  transaction.RawRequest,
			//RawResponse: transaction.RawResponse,
			Reference:   transaction.Reference,
			DeliveryFee: transaction.DeliveryFee,
			WaterCost:   transaction.WaterCost,
			CreatedAt:   transaction.CreatedAt,
		})
	}
	return list
}

func (t *transactionServiceImpl) GetDriverPendingOrder(ctx context.Context, userId float64, stage uint) []model.OrderModel {
	orders, err := t.OrderRepository.FindDriverOrdersByUserId(ctx, userId, stage)
	exception.PanicLogging(err)
	if orders == nil {
		return nil
	}
	var orderModels []model.OrderModel
	for _, order := range orders {
		orderModels = append(orderModels, model.OrderModel{
			Id:          order.ID,
			Amount:      order.Amount,
			Currency:    order.Currency,
			WaterCost:   order.WaterCost,
			DeliveryFee: order.DeliveryFee,
			RefineryId:  order.RefineryId,
			Capacity:    order.Capacity,
			Transaction: model.TransactionModel{
				ID:          order.Transaction.ID,
				Email:       order.Transaction.Email,
				PhoneNumber: order.Transaction.PhoneNumber,
				Amount:      order.Transaction.Amount,
				Currency:    order.Transaction.Currency,
				PaymentID:   order.Transaction.PaymentID,
				Provider:    order.Transaction.Provider,
				PaymentType: order.Transaction.PaymentType,
				Status:      order.Transaction.Status,
				//RawRequest:  order.Transaction.RawRequest,
				//RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			DeliveryAddress: order.Transaction.Address.Description,
			DeliveryPlaceId: order.DeliveryPlaceId,
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
			RefineryAddress: order.RefineryAddress,
			RefineryPlaceId: order.RefineryPlaceId,
			Refinery: model.RefineryModel{
				Id:      order.RefineryId,
				Name:    order.Refinery.Name,
				Email:   order.Refinery.Email,
				Phone:   order.Refinery.Phone,
				Address: order.Refinery.Address,
				Region:  order.Refinery.Region,
				PlaceId: order.Refinery.PlaceId,
			},
			User: model.UserModel{
				Username:     order.Transaction.User.Username,
				EmailAddress: order.Transaction.User.Email,
				PhoneNumber:  order.Transaction.User.PhoneNumber,
				FirstName:    order.Transaction.User.FirstName,
				LastName:     order.Transaction.User.LastName,
			},
		})

	}
	return orderModels
}

func (t *transactionServiceImpl) GetCustomerPendingOrder(ctx context.Context, userId float64, stage uint) []model.OrderModel {
	orders, err := t.OrderRepository.FindCustomerOrdersByUserId(ctx, userId, stage)
	exception.PanicLogging(err)
	if orders == nil {
		return nil
	}
	var orderModels []model.OrderModel
	for _, order := range orders {
		orderModels = append(orderModels, model.OrderModel{
			Id:          order.ID,
			Amount:      order.Amount,
			Currency:    order.Currency,
			WaterCost:   order.WaterCost,
			DeliveryFee: order.DeliveryFee,
			RefineryId:  order.RefineryId,
			Capacity:    order.Capacity,
			Transaction: model.TransactionModel{
				ID:          order.Transaction.ID,
				Email:       order.Transaction.Email,
				PhoneNumber: order.Transaction.PhoneNumber,
				Amount:      order.Transaction.Amount,
				Currency:    order.Transaction.Currency,
				PaymentID:   order.Transaction.PaymentID,
				Provider:    order.Transaction.Provider,
				PaymentType: order.Transaction.PaymentType,
				Status:      order.Transaction.Status,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			DeliveryAddress: order.Transaction.Address.Description,
			DeliveryPlaceId: order.DeliveryPlaceId,
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
			RefineryAddress: order.RefineryAddress,
			RefineryPlaceId: order.RefineryPlaceId,
			Refinery: model.RefineryModel{
				Id:      order.RefineryId,
				Name:    order.Refinery.Name,
				Email:   order.Refinery.Email,
				Phone:   order.Refinery.Phone,
				Address: order.Refinery.Address,
				Region:  order.Refinery.Region,
				PlaceId: order.Refinery.PlaceId,
			},
			User: model.UserModel{
				Username:     order.Transaction.User.Username,
				EmailAddress: order.Transaction.User.Email,
				PhoneNumber:  order.Transaction.User.PhoneNumber,
				FirstName:    order.Transaction.User.FirstName,
				LastName:     order.Transaction.User.LastName,
			},
		})

	}
	return orderModels
}

func (t *transactionServiceImpl) GetDriverCompletedOrder(ctx context.Context, userId float64, stage uint) []model.OrderModel {
	orders, err := t.OrderRepository.FindCompletedDriverOrdersByUserId(ctx, userId, stage)
	exception.PanicLogging(err)
	if orders == nil {
		return nil
	}
	var orderModels []model.OrderModel
	for _, order := range orders {
		orderModels = append(orderModels, model.OrderModel{
			Id:          order.ID,
			Amount:      order.Amount,
			Currency:    order.Currency,
			WaterCost:   order.WaterCost,
			DeliveryFee: order.DeliveryFee,
			RefineryId:  order.RefineryId,
			Capacity:    order.Capacity,
			Transaction: model.TransactionModel{
				ID:          order.Transaction.ID,
				Email:       order.Transaction.Email,
				PhoneNumber: order.Transaction.PhoneNumber,
				Amount:      order.Transaction.Amount,
				Currency:    order.Transaction.Currency,
				PaymentID:   order.Transaction.PaymentID,
				Provider:    order.Transaction.Provider,
				PaymentType: order.Transaction.PaymentType,
				Status:      order.Transaction.Status,
				//RawRequest:  order.Transaction.RawRequest,
				//RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			DeliveryAddress: order.Transaction.Address.Description,
			DeliveryPlaceId: order.DeliveryPlaceId,
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
			RefineryAddress: order.RefineryAddress,
			RefineryPlaceId: order.RefineryPlaceId,
			Refinery: model.RefineryModel{
				Id:      order.RefineryId,
				Name:    order.Refinery.Name,
				Email:   order.Refinery.Email,
				Phone:   order.Refinery.Phone,
				Address: order.Refinery.Address,
				Region:  order.Refinery.Region,
				PlaceId: order.Refinery.PlaceId,
			},
			User: model.UserModel{
				Username:     order.Transaction.User.Username,
				EmailAddress: order.Transaction.User.Email,
				PhoneNumber:  order.Transaction.User.PhoneNumber,
				FirstName:    order.Transaction.User.FirstName,
				LastName:     order.Transaction.User.LastName,
			},
		})

	}
	return orderModels
}

func (t *transactionServiceImpl) ApproveOrRejectOrder(ctx context.Context, orderModel model.ApproveOrRejectOrderModel) (interface{}, error) {
	order, err := t.OrderRepository.FindById(ctx, orderModel.OrderId)
	exception.PanicLogging(err)

	if order.ID == 0 {
		return nil, exception.BadRequestError{
			Message: "Order not found",
		}
	}

	if orderModel.Action == "approve" {
		order.Status = 2
		if orderModel.TruckId != 0 {
			order.TruckId = orderModel.TruckId

			//ToDo: Push notification to rider suing rabbitmq
			//driver := truck.User
			//if driver.Token != "" {
			//	err = t.NotificationService.SendToDevice(context.Background(), model.NotificationModel{
			//		Token:       driver.Token,
			//		Title:       "Order Assigned",
			//		Body:        "An order has been assigned to you!",
			//		Data:        map[string]string{"orderId": strconv.Itoa(int(order.ID)), "status": "ready_for_delivery"},
			//		ClickAction: "OPEN_ORDER_DETAILS",
			//	})
			//}
		}
	} else {
		order.Status -= 1
	}

	err = t.OrderRepository.Update(ctx, order)
	exception.PanicLogging(err)

	return order, nil
}

func (t *transactionServiceImpl) GetRefineryOrders(ctx context.Context, u uint, countryCode string) ([]model.OrderModel, error) {
	orders, err := t.OrderRepository.GetRefineryOrders(ctx, u, countryCode)
	exception.PanicLogging(err)
	if orders == nil {
		return nil, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}
	var orderModels []model.OrderModel
	for _, order := range orders {
		orderModels = append(orderModels, model.OrderModel{
			Id:          order.ID,
			Amount:      order.Amount,
			Currency:    order.Currency,
			WaterCost:   order.WaterCost,
			DeliveryFee: order.DeliveryFee,
			RefineryId:  order.RefineryId,
			Capacity:    order.Capacity,
			Transaction: model.TransactionModel{
				ID:          order.Transaction.ID,
				Email:       order.Transaction.Email,
				PhoneNumber: order.Transaction.PhoneNumber,
				Amount:      order.Transaction.Amount,
				Currency:    order.Transaction.Currency,
				PaymentID:   order.Transaction.PaymentID,
				Provider:    order.Transaction.Provider,
				PaymentType: order.Transaction.PaymentType,
				Status:      order.Transaction.Status,
				//RawRequest:  order.Transaction.RawRequest,
				//RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			Refinery: model.RefineryModel{
				Id:      order.RefineryId,
				Name:    order.Refinery.Name,
				Email:   order.Refinery.Email,
				Phone:   order.Refinery.Phone,
				Address: order.Refinery.Address,
				Region:  order.Refinery.Region,
				PlaceId: order.Refinery.PlaceId,
			},

			User: model.UserModel{
				Username:     order.Transaction.User.Username,
				EmailAddress: order.Transaction.User.Email,
				PhoneNumber:  order.Transaction.User.PhoneNumber,
				FirstName:    order.Transaction.User.FirstName,
				LastName:     order.Transaction.User.LastName,
			},
			Truck: model.TruckModel{
				Id:          order.Truck.ID,
				PlateNumber: order.Truck.PlateNumber,
				Capacity:    strconv.Itoa(order.Truck.Capacity),
				User: model.UserModel{
					Id:           order.Truck.User.ID,
					Username:     order.Truck.User.Username,
					EmailAddress: order.Truck.User.Email,
					FirstName:    order.Truck.User.FirstName,
					LastName:     order.Truck.User.LastName,
					PhoneNumber:  order.Truck.User.PhoneNumber,
				},
			},
			DeliveryAddress: order.Transaction.Address.Description,
			RefineryAddress: order.RefineryAddress,
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
			Rating:          order.Rating,
			Review:          order.Review,
		})
	}

	return orderModels, nil
}

func (t *transactionServiceImpl) PaymentStatus(ctx context.Context, id string) map[string]interface{} {
	transaction, err := t.TransactionRepository.FindByReference(ctx, id)
	exception.PanicLogging(err)

	if transaction.ID == 0 {

		exception.PanicLogging(errors.New("transaction does not exist"))
	}

	paymentConfig, err := t.PaymentConfigRepository.GetPaymentConfig(ctx, transaction.CountryCode)
	exception.PanicLogging(err)
	if paymentConfig == nil {
		exception.PanicLogging(errors.New("payment configuration not found for the given country code"))
	}
	header := make(map[string]interface{})
	header["Authorization"] = "Bearer " + paymentConfig.SecretKey
	paystackUrl := t.Config.Get("PAYSTACK_BASE_URL")
	response, err := t.HttpService.PostMethod(ctx, paystackUrl+"/transaction/verify/"+transaction.Reference, "GET", &map[string]interface{}{}, &header, false)
	jsn, err := json.Marshal(response)
	exception.PanicLogging(err)

	transactionStatus := GetTransactionStatus(string(jsn))
	transaction.Status = transactionStatus.Data.Status
	transaction.RawResponse = string(jsn)
	transaction.UpdatedAt = time.Now()
	transaction.Reference = transactionStatus.Data.Reference
	pan := transactionStatus.Data.Authorization.Bin + "****" + transactionStatus.Data.Authorization.Last4
	transaction.PaymentID = pan
	transaction.PaymentType = transactionStatus.Data.Channel
	transaction.Scheme = transactionStatus.Data.Authorization.Brand
	err = t.TransactionRepository.Update(ctx, transaction)
	exception.PanicLogging(err)
	var order entity.Order = entity.Order{}
	var request model.MobileMoneyRequestModel
	err = json.Unmarshal([]byte(transaction.RawRequest), &request)
	exception.PanicLogging(err)
	if transactionStatus.Data.Status == "success" {
		order, _ = t.OrderRepository.FindByTransactionId(ctx, transaction.ID)

		if order.ID == 0 {
			order = entity.Order{
				TransactionId: transaction.ID,
				//UserId:          transaction.UserId,
				Amount:          transaction.Amount,
				Currency:        transaction.Currency,
				WaterCost:       transaction.WaterCost,
				DeliveryFee:     transaction.DeliveryFee,
				DeliveryAddress: request.CustomerAddress,
				DeliveryPlaceId: request.CustomerPlaceId,
				RefineryAddress: request.RefineryAddress,
				RefineryPlaceId: request.RefineryPlaceId,
				RefineryId:      request.RefineryId,
				Capacity:        transaction.Capacity,
				Type:            transaction.Type,
				WaterType:       request.Type,
				AddressId:       transaction.AddressId,
			}
			order = t.OrderRepository.Insert(ctx, order)
		}
		var uniqueId string
		if transactionStatus.Data.Channel == "mobile_money" {
			uniqueId = transaction.PhoneNumber
		} else {
			uniqueId = pan
		}

		//save payment method
		paymentMethod := entity.PaymentMethod{
			UserID:   request.UserId,
			Provider: request.Provider,
			UniqueId: uniqueId,
			Scheme:   transaction.Scheme,
			RawData:  transaction.RawResponse,
			AuthCode: transactionStatus.Data.Authorization.AuthorizationCode,
			Name: transactionStatus.Data.Customer.FirstName + " " +
				transactionStatus.Data.Customer.LastName,
		}

		_, err = t.PaymentMethodRepository.Create(ctx, paymentMethod)
	}
	transactionStatus.Data.TransactionId = transaction.ID
	status := map[string]interface{}{
		"status":          transactionStatus.Data.Status,
		"transactionId":   transaction.ID,
		"reference":       transactionStatus.Data.Reference,
		"amount":          transaction.Amount,
		"currency":        transactionStatus.Data.Currency,
		"createdAt":       transactionStatus.Data.CreatedAt,
		"scheme":          transactionStatus.Data.Authorization.Brand,
		"paymentType":     transactionStatus.Data.Channel,
		"paymentId":       transaction.PaymentID,
		"capacity":        transaction.Capacity,
		"deliveryAddress": request.CustomerAddress,
		"pan":             pan,
		"waterType":       order.WaterType,
		"orderId":         order.ID,
		"orderStatus":     order.Status,
		"user": model.UserModel{
			Id:           transaction.User.ID,
			Username:     transaction.User.Username,
			EmailAddress: transaction.User.Email,
			PhoneNumber:  transaction.User.PhoneNumber,
			FirstName:    transaction.User.FirstName,
			LastName:     transaction.User.LastName,
		},
	}
	return status
}

func (t *transactionServiceImpl) InitiateMobileMoneyTransaction(ctx context.Context, request model.MobileMoneyRequestModel) interface{} {
	// Load payment configuration based on country code
	paymentConfig, err := t.PaymentConfigRepository.GetPaymentConfig(ctx, request.CountryCode)
	exception.PanicLogging(err)
	if paymentConfig == nil {
		exception.PanicLogging(errors.New("payment configuration not found for the given country code"))
	}

	transaction := t.insertTransaction(ctx, request)

	var email string
	if request.EmailAddress == "" || !strings.Contains(request.EmailAddress, "@") {
		email = "Intelblue28@gmail.com"
	} else {
		email = request.EmailAddress
	}

	data := make(map[string]interface{})
	data["amount"] = int(request.Amount * 100)
	data["email"] = email
	data["currency"] = transaction.Currency
	var paystackUrl string
	if request.Provider != "card" {
		paystackUrl = t.Config.Get("PAYSTACK_BASE_URL") + "/charge"
		mobileMoney := make(map[string]interface{})
		mobileMoney["phone"] = request.PhoneNumber
		mobileMoney["provider"] = request.Provider
		data["mobile_money"] = mobileMoney
	} else {
		paystackUrl = t.Config.Get("PAYSTACK_BASE_URL") + "/transaction/initialize"
		var chn []string
		chn = append(chn, "card")
		data["callback_url"] = "https://www.aquawizz.com/redirect-url"
	}

	response, err := t.SendToPayStack(ctx, paystackUrl, data, paymentConfig.SecretKey)
	if request.Provider == "card" {
		transaction.Reference = response["data"].(map[string]interface{})["reference"].(string)
		err = t.TransactionRepository.Update(ctx, transaction)
		var trxData = response["data"].(map[string]interface{})
		trxData["transaction_id"] = transaction.ID
		response["data"] = trxData
		return response
	}
	jsn, err := json.Marshal(response)

	transactionStatus := GetTransactionStatus(string(jsn))
	transactionStatus.Data.TransactionId = transaction.ID
	transaction.Status = transactionStatus.Data.Status
	transaction.RawResponse = string(jsn)
	transaction.UpdatedAt = time.Now()
	transaction.Reference = transactionStatus.Data.Reference
	err = t.TransactionRepository.Update(ctx, transaction)
	exception.PanicLogging(err)

	return transactionStatus
}

func (t *transactionServiceImpl) insertTransaction(ctx context.Context, request model.MobileMoneyRequestModel) entity.Transaction {
	rawRequest, err := json.Marshal(request)
	exception.PanicLogging(err)
	transaction := entity.Transaction{
		Email:       request.EmailAddress,
		PhoneNumber: request.PhoneNumber,
		Amount:      request.Amount,
		Currency:    request.Currency,
		UserId:      request.UserId,
		PaymentID:   request.PhoneNumber,
		Provider:    request.Provider,
		PaymentType: "mobile-money",
		Status:      "Initiated",
		RawRequest:  string(rawRequest),
		WaterCost:   request.WaterCost,
		DeliveryFee: request.DeliveryFee,
		Capacity:    request.Capacity,
		Type:        request.Type,
		AddressId:   request.AddressId,
		RefineryId:  request.RefineryId,
		CountryCode: request.CountryCode,
	}
	transaction = t.TransactionRepository.Insert(ctx, transaction)
	return transaction
}

func GetTransactionStatus(data string) model.TransactionStatusModel {
	var transactionStatus model.TransactionStatusModel
	err := json.Unmarshal([]byte(data), &transactionStatus)
	exception.PanicLogging(err)
	return transactionStatus
}

func (t *transactionServiceImpl) GetRefineryDashboardData(ctx context.Context, u uint) (map[string]interface{}, error) {
	refineryData, err := t.TransactionRepository.GetRefineryDashboardData(ctx, u)
	exception.PanicLogging(err)
	if refineryData == nil {
		return nil, exception.BadRequestError{
			Message: "Refinery not found",
		}
	}

	return refineryData, nil
}
func (t *transactionServiceImpl) FindById(ctx context.Context, id uint) (model.OrderModel, error) {
	order, err := t.OrderRepository.FindById(ctx, id)
	exception.PanicLogging(err)
	longitude, _ := strconv.ParseFloat(order.Transaction.Address.Longitude, 64)
	latitude, _ := strconv.ParseFloat(order.Transaction.Address.Latitude, 64)
	orderModel := model.OrderModel{
		Id:            order.ID,
		TransactionId: order.TransactionId,
		Transaction: model.TransactionModel{
			ID:          order.Transaction.ID,
			Email:       order.Transaction.Email,
			PhoneNumber: order.Transaction.PhoneNumber,
			Amount:      order.Transaction.Amount,
			Currency:    order.Transaction.Currency,
			PaymentID:   order.Transaction.PaymentID,
			PaymentType: order.Transaction.PaymentType,
			Status:      order.Transaction.Status,
			//RawRequest:  order.Transaction.RawRequest,
			//RawResponse: order.Transaction.RawResponse,
			WaterCost:   order.Transaction.WaterCost,
			DeliveryFee: order.Transaction.DeliveryFee,

			CreatedAt: order.Transaction.CreatedAt,
			DeliveryAddress: model.AddressModel{
				Description: order.Transaction.Address.Description,
				Latitude:    latitude,
				Longitude:   longitude,
				PlaceId:     order.Transaction.Address.PlaceId,
			},
		},
		Amount:          order.Amount,
		Currency:        order.Currency,
		WaterCost:       order.WaterCost,
		DeliveryFee:     order.DeliveryFee,
		DeliveryAddress: order.Transaction.Address.Description,
		DeliveryPlaceId: order.DeliveryPlaceId,
		RefineryAddress: order.RefineryAddress,
		RefineryPlaceId: order.RefineryPlaceId,
		RefineryId:      order.RefineryId,
		Capacity:        order.Capacity,
		Type:            order.Type,
		Refinery: model.RefineryModel{
			Id:        order.RefineryId,
			Name:      order.Refinery.Name,
			Address:   order.Refinery.Address,
			Phone:     order.Refinery.Phone,
			Email:     order.Refinery.Email,
			Latitude:  order.Refinery.Latitude,
			Longitude: order.Refinery.Longitude,
			PlaceId:   order.Refinery.PlaceId,
		},
		User: model.UserModel{
			Id:           order.Transaction.User.ID,
			FirstName:    order.Transaction.User.FirstName,
			LastName:     order.Transaction.User.LastName,
			EmailAddress: order.Transaction.User.Email,
			PhoneNumber:  order.Transaction.User.PhoneNumber,
		},
		Status:      order.Status,
		OrderStatus: order.Status,
	}
	return orderModel, nil
}
func (t *transactionServiceImpl) SendToPayStack(ctx context.Context, url string, data map[string]interface{}, secretKey string) (map[string]interface{}, error) {

	header := make(map[string]interface{})
	header["Authorization"] = "Bearer " + secretKey

	response, err := t.HttpService.PostMethod(ctx, url, "POST", &data, &header, false)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Add this method to your TransactionServiceImpl
func (t *transactionServiceImpl) ProcessPendingTransactions(ctx context.Context, truck model.TruckModel) error {
	// Get transactions that have been pending for more than 24 hours
	initiatedOrders, err := t.OrderRepository.FindInitiatedOrders(ctx, time.Hour*24)
	if err != nil {
		return err
	}

	for _, order := range initiatedOrders {

		order.Status = 1
		order.UpdatedAt = time.Now()
		//order.TruckId = &truck.Id

		err := t.OrderRepository.Update(ctx, order)

		user := order.Transaction.User
		if user.FcmToken != "" {
			//return errors.New("user FCM token not found")

			err = t.NotificationService.SendToDevice(context.Background(), model.NotificationModel{
				Token:       user.FcmToken,
				Title:       "Order Ready",
				Body:        "Your order has been confirm and is processing!",
				Data:        map[string]string{"orderId": strconv.Itoa(int(order.ID)), "status": "ready_for_delivery"},
				ClickAction: "OPEN_ORDER_DETAILS",
			})
		}

		if err != nil {
			logger.Logger.Error("Failed to update transaction " + ": " + err.Error())
			// Continue processing other transactions even if one fails
			continue
		}

		// Additional processing like sending notifications could be done here
	}

	return nil
}

func (t *transactionServiceImpl) MarkOrderReadyForDelivery(id string) error {
	order, err := t.OrderRepository.MarkOrderReadyForDelivery(id)
	if err != nil {
		return err
	}
	// Assume you can get the user and their FCM token from the order
	user := order.Transaction.User
	if user.FcmToken != "" {
		//return errors.New("user FCM token not found")

		err = t.NotificationService.SendToDevice(context.Background(), model.NotificationModel{
			Token:       user.FcmToken,
			Title:       "Order Ready",
			Body:        "Your order is ready for delivery!",
			Data:        map[string]string{"orderId": strconv.Itoa(int(order.ID)), "status": "ready_for_delivery"},
			ClickAction: "OPEN_ORDER_DETAILS",
		})
	}
	//driver := order.Truck.User
	//if driver.FcmToken != "" {
	//	err = t.NotificationService.SendToDevice(context.Background(), model.NotificationModel{
	//		Token:       user.FcmToken,
	//		Title:       "Order Assigned",
	//		Body:        "An order has been assigned to you!",
	//		Data:        map[string]string{"orderId": strconv.Itoa(int(order.ID)), "status": "ready_for_delivery"},
	//		ClickAction: "OPEN_ORDER_DETAILS",
	//	})
	//}

	//if err != nil {
	//	return err
	//}

	return nil

}

func (t *transactionServiceImpl) CloseOrder(id string) error {
	order, err := t.OrderRepository.CloseOrder(id)
	if err != nil {
		return err
	}
	// Assume you can get the user and their FCM token from the order
	user := order.Transaction.User
	if user.FcmToken == "" {
		//return errors.New("user FCM token not found")
		return errors.New("user FCM token not found")
	}

	notification := model.NotificationModel{
		Token:       user.FcmToken,
		Title:       "Order Completed",
		Body:        "Your order has been completed successfully!",
		Data:        map[string]string{"orderId": strconv.Itoa(int(order.ID)), "status": "order_completed"},
		ClickAction: "FLUTTER_NOTIFICATION_CLICK",
	}

	ctx := context.Background()
	err = t.NotificationService.SendToDevice(ctx, notification)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (t *transactionServiceImpl) SubmitRating(ctx context.Context, ratingModel model.RatingModel) error {
	if ratingModel.OrderId == 0 {
		return errors.New("order ID is required")
	}
	if ratingModel.Rating < 1 || ratingModel.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	order, err := t.OrderRepository.FindById(ctx, ratingModel.OrderId)
	if err != nil {
		return err
	}
	if order.ID == 0 {
		return errors.New("order not found")
	}

	order.Rating = ratingModel.Rating
	order.Review = ratingModel.Review
	err = t.OrderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderRating returns the rating and review for a given order id
func (t *transactionServiceImpl) GetOrderRating(ctx context.Context, orderId uint) (model.RatingModel, error) {
	var rating model.RatingModel
	if orderId == 0 {
		return rating, errors.New("order ID is required")
	}

	order, err := t.OrderRepository.FindById(ctx, orderId)
	if err != nil {
		return rating, err
	}

	if order.ID == 0 {
		return rating, errors.New("order not found")
	}

	rating = model.RatingModel{
		Review:  order.Review,
		Rating:  order.Rating,
		UserId:  order.Transaction.UserId,
		OrderId: order.ID,
	}

	return rating, nil
}

// GetTransactionsByCountryCode returns transactions filtered by country code or all if code is empty
func (t *transactionServiceImpl) GetTransactionsByCountryCode(ctx context.Context, countryCode string, page, limit int) ([]model.TransactionModel, int64) {
	transactions, totalCount := t.TransactionRepository.FindByCountryCode(ctx, countryCode, page, limit)

	var transactionModels []model.TransactionModel
	for _, transaction := range transactions {
		transactionModels = append(transactionModels, model.TransactionModel{
			ID:          transaction.ID,
			UserId:      transaction.UserId,
			Amount:      transaction.Amount,
			PhoneNumber: transaction.PhoneNumber,
			Email:       transaction.Email,
			Provider:    transaction.Provider,
			Status:      transaction.Status,
			PaymentID:   transaction.PaymentID,
			PaymentType: transaction.PaymentType,
			CompletedAt: transaction.CompletedAt,
			Currency:    transaction.Currency,
			Reference:   transaction.Reference,
			Scheme:      transaction.Scheme,
			CountryCode: transaction.CountryCode,
			CreatedAt:   transaction.CreatedAt,
		})
	}

	return transactionModels, totalCount
}
