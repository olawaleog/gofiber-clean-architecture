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
	protected.Get("/transaction-list", c.GetTransactions)
	protected.Get("/transactions-by-country", c.GetTransactionsByCountryCode)
	protected.Get("/mark-order-ready-for-delivery/:id", c.MarkOrderReadyForDelivery)
	protected.Get("/close-order/:id", c.CloseOrder)
	protected.Get("/order/:id", c.FindById)
	protected.Post("/submit-rating", c.SubmitRating)
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
	mobileMoneyRequestModel.UserId = claims["userId"].(uint)
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

func (c TransactionController) GetCustomerPendingOrder(ctx *fiber.Ctx) error {
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	userId := claims["userId"].(float64)
	orders := c.TransactionService.GetCustomerPendingOrder(ctx.Context(), userId, 1)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    orders,
		Success: true,
	})
}

func (c TransactionController) GetDriverCompletedOrder(ctx *fiber.Ctx) error {
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	userId := claims["userId"].(float64)
	orders := c.TransactionService.GetDriverCompletedOrder(ctx.Context(), userId, 1)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    orders,
		Success: true,
	})
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
	var claims map[string]interface{}
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	userId := claims["userId"].(float64)
	orders, err := c.TransactionService.GetCustomerOrders(ctx.Context(), uint(userId))
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    orders,
		Success: true,
	})
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

func (c TransactionController) ProcessRecurringPayment(ctx *fiber.Ctx) error {
	var mobileMoneyRequestModel model.MobileMoneyRequestModel
	err := ctx.BodyParser(&mobileMoneyRequestModel)
	exception.PanicLogging(err)
	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	//mobileMoneyRequestModel.PhoneNumber = claims["phoneNumber"].(string)
	mobileMoneyRequestModel.UserId = claims["userId"].(uint)
	mobileMoneyRequestModel.EmailAddress = claims["emailAddress"].(string)
	response := c.TransactionService.ProcessRecurringPayment(ctx.Context(), mobileMoneyRequestModel)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    response,
		Success: true,
	})

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
	exception.PanicLogging(err)

	token := ctx.Get("Authorization")
	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	ratingModel.UserId = uint(claims["userId"].(float64))

	err = c.TransactionService.SubmitRating(ctx.Context(), ratingModel)
	exception.PanicLogging(err)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Rating submitted successfully",
		Data:    nil,
		Success: true,
	})
}

// GetTransactionsByCountryCode returns transactions filtered by country code or all transactions
func (c TransactionController) GetTransactionsByCountryCode(ctx *fiber.Ctx) error {
	// Get pagination parameters from the query string
	pageStr := ctx.Query("page", "1")    // default to page 1 if not provided
	limitStr := ctx.Query("limit", "10") // default to 10 items per page if not provided

	// Convert string parameters to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get country code from token
	countryCode := ""
	token := ctx.Get("Authorization")
	if token != "" {
		claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
		if err == nil && claims["countryCode"] != nil {
			// Use the country code from JWT claims
			countryCode, _ = claims["countryCode"].(string)
		}
	}

	// Allow override via query parameter for admin users (optional)
	queryCountryCode := ctx.Query("country_code", "")
	if queryCountryCode != "" && c.isAdminUser(ctx) {
		countryCode = queryCountryCode
	}

	// Get transactions with country code filter and pagination
	transactions, totalCount := c.TransactionService.GetTransactionsByCountryCode(ctx.Context(), countryCode, page, limit)

	// Create pagination response
	paginatedResponse := model.NewPaginationResponse(transactions, page, limit, totalCount)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    paginatedResponse,
		Success: true,
	})
}

// isAdminUser checks if the current user has admin privileges
func (c TransactionController) isAdminUser(ctx *fiber.Ctx) bool {
	token := ctx.Get("Authorization")
	if token == "" {
		return false
	}

	claims, err := c.UserService.GetClaimsFromToken(ctx.Context(), token)
	if err != nil {
		return false
	}

	// Check the role in claims
	if claims["roles"] != nil {
		roles, ok := claims["roles"].([]interface{})
		if ok && len(roles) > 0 {
			for _, role := range roles {
				if roleMap, ok := role.(map[string]interface{}); ok {
					if roleVal, exists := roleMap["role"].(string); exists && roleVal == "admin" {
						return true
					}
				}
			}
		}
	}

	// Check user role field directly
	if userRole, exists := claims["role"].(string); exists && userRole == "admin" {
		return true
	}

	return false
}
