package impl

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/repository"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"time"
)

func NewTransactionServiceImpl(transactionRepository *repository.TransactionRepository, orderRepo *repository.OrderRepository, paymentRepo *repository.PaymentMethodRepository, client *service.HttpService, config configuration.Config) service.TransactionService {
	return &transactionServiceImpl{TransactionRepository: *transactionRepository, HttpService: *client, Config: config, OrderRepository: *orderRepo, PaymentMethodRepository: *paymentRepo}
}

type transactionServiceImpl struct {
	repository.TransactionRepository
	repository.OrderRepository
	repository.PaymentMethodRepository
	service.HttpService
	configuration.Config
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
				RawRequest:  order.Transaction.RawRequest,
				RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			Status:        order.Status,
			TransactionId: order.TransactionId,
			TruckId:       order.TruckId,
			CreatedAt:     order.CreatedAt,
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
			RawRequest:  transaction.RawRequest,
			RawResponse: transaction.RawResponse,
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
				RawRequest:  order.Transaction.RawRequest,
				RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			DeliveryAddress: order.DeliveryAddress,
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
		order.Status += 1
		if orderModel.TruckId != 0 {
			order.TruckId = orderModel.TruckId
		}
	} else {
		order.Status -= 1
	}

	err = t.OrderRepository.Update(ctx, order)
	exception.PanicLogging(err)

	return order, nil
}

func (t *transactionServiceImpl) GetRefineryOrders(ctx context.Context, u uint) ([]model.OrderModel, error) {
	orders, err := t.OrderRepository.GetRefineryOrders(ctx, u)
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
				RawRequest:  order.Transaction.RawRequest,
				RawResponse: order.Transaction.RawResponse,
				Reference:   order.Transaction.Reference,
				DeliveryFee: order.Transaction.DeliveryFee,
				WaterCost:   order.Transaction.WaterCost,
			},
			DeliveryAddress: order.DeliveryAddress,
			RefineryAddress: order.RefineryAddress,
			Status:          order.Status,
			TransactionId:   order.TransactionId,
			TruckId:         order.TruckId,
			CreatedAt:       order.CreatedAt,
		})
	}

	return orderModels, nil
}

func (t *transactionServiceImpl) PaymentStatus(ctx context.Context, id string) model.TransactionStatusModel {
	transaction, err := t.TransactionRepository.FindByReference(ctx, id)
	exception.PanicLogging(err)

	if transaction.ID == 0 {

		exception.PanicLogging(errors.New("transaction does not exist"))
	}
	header := make(map[string]interface{})
	header["Authorization"] = "Bearer " + t.Config.Get("PAYSTACK_SECRET_KEY")
	paystackUrl := t.Config.Get("PAYSTACK_BASE_URL")
	response := t.HttpService.PostMethod(ctx, paystackUrl+"/transaction/verify/"+transaction.Reference, "GET", &map[string]interface{}{}, &header, false)
	jsn, err := json.Marshal(response)
	exception.PanicLogging(err)

	transactionStatus := GetTransactionStatus(string(jsn))
	transaction.Status = transactionStatus.Data.Status
	transaction.RawResponse = string(jsn)
	transaction.UpdatedAt = time.Now()
	err = t.TransactionRepository.Update(ctx, transaction)
	exception.PanicLogging(err)
	var order entity.Order
	if transactionStatus.Data.Status == "success" {
		order, _ = t.OrderRepository.FindByTransactionId(ctx, transaction.ID)
		var request model.MobileMoneyRequestModel
		err := json.Unmarshal([]byte(transaction.RawRequest), &request)
		exception.PanicLogging(err)

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
			}
			order = t.OrderRepository.Insert(ctx, order)
		}

		//save payment method
		paymentMethod := entity.PaymentMethod{
			UserID:   request.UserId,
			Provider: request.Provider,
			UniqueId: transaction.PhoneNumber,
			Scheme:   transaction.Scheme,
			RawData:  transaction.RawResponse,
			AuthCode: transactionStatus.Data.Authorization.AuthorizationCode,
		}

		_, err = t.PaymentMethodRepository.Create(ctx, paymentMethod)
	}
	transactionStatus.Data.TransactionId = order.ID
	return transactionStatus
}

func (t *transactionServiceImpl) InitiateMobileMoneyTransaction(ctx context.Context, request model.MobileMoneyRequestModel) interface{} {
	//amount, err := strconv.ParseFloat(request.Amount, 64)
	//exception.PanicLogging(err)
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
	}
	transaction = t.TransactionRepository.Insert(ctx, transaction)

	mobileMoney := make(map[string]interface{})
	mobileMoney["phone"] = request.PhoneNumber
	mobileMoney["provider"] = request.Provider
	var email string
	if request.EmailAddress == "" {
		email = "Intelblue28@gmail.com"
	} else {
		email = request.EmailAddress
	}

	data := make(map[string]interface{})
	data["amount"] = int(request.Amount * 100)
	data["email"] = email
	data["currency"] = request.Currency
	data["mobile_money"] = mobileMoney
	paystackUrl := t.Config.Get("PAYSTACK_BASE_URL")

	header := make(map[string]interface{})
	header["Authorization"] = "Bearer " + t.Config.Get("PAYSTACK_SECRET_KEY")

	response := t.HttpService.PostMethod(ctx, paystackUrl+"/charge", "POST", &data, &header, false)
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
	orderModel := model.OrderModel{
		Id:              order.ID,
		TransactionId:   order.TransactionId,
		Amount:          order.Amount,
		Currency:        order.Currency,
		WaterCost:       order.WaterCost,
		DeliveryFee:     order.DeliveryFee,
		DeliveryAddress: order.DeliveryAddress,
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
	}
	return orderModel, nil
}
