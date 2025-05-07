package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	service.TransactionService
	service.UserService
	configuration.Config
}

func NewTransactionController(transactionService *service.TransactionService, userService *service.UserService, config configuration.Config) *TransactionController {
	return &TransactionController{TransactionService: *transactionService, Config: config, UserService: *userService}
}

func (c TransactionController) Route(app *fiber.App) {
	app.Post("/v1/api/payment-mobile-money", c.InitiateMobileMoneyPayment)
	app.Get("/v1/api/payment-status/:id", c.PaymentStatus)
	app.Get("/v1/api/refinery-dashboard-data", c.GetRefineryDashboardData)
	app.Get("/v1/api/pending-orders", c.GetRefineryOrders)
	app.Post("/v1/api/approve-or-reject-order", c.ApproveOrRejectOrder)
	app.Get("/v1/api/get-driver-pending-orders", c.GetDriverPendingOrder)

	//app.Post("/v1/api/paystack/webook", controller.PayStackWebhook)
}

func (c TransactionController) InitiateMobileMoneyPayment(ctx *fiber.Ctx) error {
	var mobileMoneyRequestModel model.MobileMoneyRequestModel

	err := ctx.BodyParser(&mobileMoneyRequestModel)

	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	//mobileMoneyRequestModel.PhoneNumber = claims["phoneNumber"].(string)
	mobileMoneyRequestModel.UserId = claims["userId"].(float64)
	mobileMoneyRequestModel.EmailAddress = claims["emailAddress"].(string)
	response := c.TransactionService.InitiateMobileMoneyTransaction(ctx.Context(), mobileMoneyRequestModel)
	var transactionStatus = response.(model.TransactionStatusModel)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Success: true,
		Code:    fiber.StatusOK,
		Message: transactionStatus.Message,
		Data:    response,
	})

}

func (c TransactionController) PaymentStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentStatus := c.TransactionService.PaymentStatus(ctx.Context(), id)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    paymentStatus,
		Success: true,
	})
}

//	func (controller TransactionController) PayStackWebhook(ctx *fiber.Ctx) error {
//		var status model.TransactionStatusModel
//		err := ctx.BodyParser(&status)
//	}
func (c TransactionController) GetRefineryDashboardData(ctx *fiber.Ctx) error {
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	refineryId := claims["refineryId"].(float64)
	refineryData, err := c.TransactionService.GetRefineryDashboardData(ctx.Context(), uint(refineryId))
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refineryData,
		Success: true,
	})
}

func (c TransactionController) GetRefineryOrders(ctx *fiber.Ctx) error {
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	refineryId := claims["refineryId"].(float64)
	orders, err := c.TransactionService.GetRefineryOrders(ctx.Context(), uint(refineryId))
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    orders,
		Success: true,
	})
}

func (c TransactionController) ApproveOrRejectOrder(ctx *fiber.Ctx) error {
	var approveOrRejectOrderModel model.ApproveOrRejectOrderModel
	err := ctx.BodyParser(&approveOrRejectOrderModel)
	exception.PanicLogging(err)

	approveOrderResponse, err := c.TransactionService.ApproveOrRejectOrder(ctx.Context(), approveOrRejectOrderModel)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    approveOrderResponse,
		Success: true,
	})
}

func (c TransactionController) GetDriverPendingOrder(ctx *fiber.Ctx) error {
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	userId := claims["userId"].(float64)
	orders := c.TransactionService.GetDriverPendingOrder(ctx.Context(), userId, 1)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    orders,
		Success: true,
	})
}
