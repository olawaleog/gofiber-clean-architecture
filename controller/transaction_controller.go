package controller

import (
	"strconv"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/middleware"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/utils"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	service.TransactionService
	service.UserService
	service.AuthorizationService
	configuration.Config
	responseBuilder *utils.ResponseBuilder
}

func NewTransactionController(
	transactionService *service.TransactionService,
	userService *service.UserService,
	authService *service.AuthorizationService,
	config configuration.Config,
) *TransactionController {
	return &TransactionController{
		TransactionService:   *transactionService,
		UserService:          *userService,
		AuthorizationService: *authService,
		Config:               config,
		responseBuilder:      utils.NewResponseBuilder(),
	}
}

func (c TransactionController) Route(app *fiber.App) {
	// Apply the JWT claims middleware to all authenticated routes
	api := app.Group("/v1/api")

	// Public routes (no auth required)
	api.Get("/payment-status/:id", c.PaymentStatus)

	// Protected routes (require authentication)
	protected := api.Group("/", middleware.ExtractClaims(c.UserService), middleware.RequireClaims())
	protected.Post("/initialize-card-transaction", c.InitiateMobileMoneyPayment)
	protected.Post("/payment-mobile-money", c.InitiateMobileMoneyPayment)
	protected.Post("/recurring-payment", c.ProcessRecurringPayment)
	protected.Get("/refinery-dashboard-data", c.GetRefineryDashboardData)
	protected.Get("/admin-dashboard-data", c.GetAdminDashboardData)
	protected.Get("/pending-orders", c.GetRefineryOrders)
	protected.Post("/approve-or-reject-order", c.ApproveOrRejectOrder)
	protected.Get("/get-driver-pending-orders", c.GetDriverPendingOrder)
	protected.Get("/get-customer-pending-orders", c.GetCustomerPendingOrder)
	protected.Get("/get-driver-completed-orders", c.GetDriverCompletedOrder)
	protected.Get("/get-customer-orders", c.GetCustomerOrders)
	protected.Get("/list", c.GetTransactions)
	protected.Get("/transactions-by-country", c.GetTransactionsByCountryCode)
	protected.Get("/mark-order-ready-for-delivery/:id", c.MarkOrderReadyForDelivery)
	protected.Get("/close-order/:id", c.CloseOrder)
	protected.Get("/order/:id", c.FindById)
	protected.Post("/submit-rating", c.SubmitRating)
	// Get rating for a specific order
	protected.Get("/order/:id/rating", c.GetOrderRating)
}

func (c TransactionController) InitiateMobileMoneyPayment(ctx *fiber.Ctx) error {
	var mobileMoneyRequestModel model.MobileMoneyRequestModel
	err := ctx.BodyParser(&mobileMoneyRequestModel)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Invalid request body"))
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	// Set user information from claims
	mobileMoneyRequestModel.UserId = uint(claims["userId"].(float64))
	mobileMoneyRequestModel.EmailAddress = claims["emailAddress"].(string)

	response := c.TransactionService.InitiateMobileMoneyTransaction(ctx.Context(), mobileMoneyRequestModel)

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(response, "Mobile money payment initiated successfully"))
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
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	refineryId, ok := claims["refineryId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Refinery ID not found in token"))
	}

	refineryData, err := c.TransactionService.GetRefineryDashboardData(ctx.Context(), uint(refineryId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(c.responseBuilder.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(refineryData, "Dashboard data retrieved successfully"))
}

func (c TransactionController) GetRefineryOrders(ctx *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	refineryId, ok := claims["refineryId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Refinery ID not found in token"))
	}

	// Extract country code from claims
	countryCode := ""
	if claims["countryCode"] != nil {
		countryCode, _ = claims["countryCode"].(string)
	}

	// Get refinery orders with country code filter
	orders, err := c.TransactionService.GetRefineryOrders(ctx.Context(), uint(refineryId), countryCode)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(c.responseBuilder.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(orders, "Orders retrieved successfully"))
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
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "User ID not found in token"))
	}

	orders := c.TransactionService.GetDriverPendingOrder(ctx.Context(), userId, 1)

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(orders, "Driver pending orders retrieved successfully"))
}

func (c TransactionController) GetCustomerPendingOrder(ctx *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "User ID not found in token"))
	}

	orders := c.TransactionService.GetCustomerPendingOrder(ctx.Context(), userId, 1)

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(orders, "Customer pending orders retrieved successfully"))
}

func (c TransactionController) GetDriverCompletedOrder(ctx *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "User ID not found in token"))
	}

	orders := c.TransactionService.GetDriverCompletedOrder(ctx.Context(), userId, 1)

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(orders, "Driver completed orders retrieved successfully"))
}

func (c TransactionController) GetTransactions(ctx *fiber.Ctx) error {
	// Extract pagination parameters using our utility
	pagination := utils.ExtractPaginationParams(ctx)

	// Get transactions with pagination
	transactions, totalCount := c.TransactionService.GetTransactionsPaginated(ctx.Context(), pagination.Page, pagination.Limit)

	// Use response builder for consistent response format
	return ctx.Status(fiber.StatusOK).JSON(
		c.responseBuilder.Pagination(transactions, pagination.Page, pagination.Limit, totalCount),
	)
}

func (c TransactionController) GetAdminDashboardData(ctx *fiber.Ctx) error {
	refineryData, err := c.TransactionService.GetAdminDashboardData(ctx.Context())
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    refineryData,
		Success: true,
	})
}

func (c TransactionController) GetCustomerOrders(ctx *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "User ID not found in token"))
	}

	orders, err := c.TransactionService.GetCustomerOrders(ctx.Context(), uint(userId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(c.responseBuilder.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(orders, "Customer orders retrieved successfully"))
}

func (c TransactionController) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := strconv.Atoi(id)
	order, err := c.TransactionService.FindById(ctx.Context(), uint(parsedId))
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    order,
		Success: true,
	})
}

func (c TransactionController) GetOrderRating(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Invalid order id"))
	}

	rating, err := c.TransactionService.GetOrderRating(ctx.Context(), uint(parsedId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(c.responseBuilder.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(rating, "Order rating retrieved successfully"))
}

func (c TransactionController) ProcessRecurringPayment(ctx *fiber.Ctx) error {
	var mobileMoneyRequestModel model.MobileMoneyRequestModel
	err := ctx.BodyParser(&mobileMoneyRequestModel)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Invalid request body"))
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	// Set user information from claims
	mobileMoneyRequestModel.UserId = claims["userId"].(uint)
	mobileMoneyRequestModel.EmailAddress = claims["emailAddress"].(string)

	response := c.TransactionService.ProcessRecurringPayment(ctx.Context(), mobileMoneyRequestModel)

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(response, "Recurring payment processed successfully"))
}

func (c TransactionController) MarkOrderReadyForDelivery(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.TransactionService.MarkOrderReadyForDelivery(id)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    nil,
		Success: true,
	})
}

func (c TransactionController) CloseOrder(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.TransactionService.CloseOrder(id)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    nil,
		Success: true,
	})
}

func (c TransactionController) SubmitRating(ctx *fiber.Ctx) error {
	var ratingModel model.RatingModel
	err := ctx.BodyParser(&ratingModel)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "Invalid request body"))
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(c.responseBuilder.Error(fiber.StatusUnauthorized, "Authentication required"))
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(c.responseBuilder.Error(fiber.StatusBadRequest, "User ID not found in token"))
	}

	ratingModel.UserId = uint(userId)

	err = c.TransactionService.SubmitRating(ctx.Context(), ratingModel)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(c.responseBuilder.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(c.responseBuilder.Success(nil, "Rating submitted successfully"))
}

// GetTransactionsByCountryCode returns transactions filtered by country code or all transactions
func (c TransactionController) GetTransactionsByCountryCode(ctx *fiber.Ctx) error {
	// Extract pagination parameters using our utility
	pagination := utils.ExtractPaginationParams(ctx)

	// Get country code from claims
	countryCode := ""
	claims, ok := middleware.GetClaims(ctx)
	if ok && claims["countryCode"] != nil {
		countryCode, _ = claims["countryCode"].(string)
	}

	// Allow override via query parameter for admin users
	queryCountryCode := ctx.Query("country_code", "")
	if queryCountryCode != "" && c.AuthorizationService.IsAdmin(claims) {
		countryCode = queryCountryCode
	}

	// Get transactions with country code filter and pagination
	transactions, totalCount := c.TransactionService.GetTransactionsByCountryCode(
		ctx.Context(), countryCode, pagination.Page, pagination.Limit)

	// Use response builder for consistent response format
	return ctx.Status(fiber.StatusOK).JSON(
		c.responseBuilder.Pagination(transactions, pagination.Page, pagination.Limit, totalCount),
	)
}

// isAdminUser checks if the current user has admin privileges
func (c TransactionController) isAdminUser(ctx *fiber.Ctx) bool {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return false
	}

	// Use AuthorizationService to check if the user is an admin
	return c.AuthorizationService.IsAdmin(claims)
}
