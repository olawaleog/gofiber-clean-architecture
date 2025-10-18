package controller

import (
	"context"
	"strconv"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/middleware"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

func NewUserController(userService *service.UserService, config configuration.Config) *UserController {
	return &UserController{UserService: *userService, Config: config}
}

type UserController struct {
	service.UserService
	configuration.Config
}

func (controller UserController) Route(app *fiber.App) {
	// Public routes (no auth required)
	app.Post("/v1/api/authentication", controller.HandleAuthentication)
	app.Post("/v1/api/register", controller.HandleRegister)
	app.Post("/v1/api/register-customer", controller.HandleRegisterCustomer)
	app.Post("/v1/api/reset-password", controller.HandleResetPassword)
	app.Post("/v1/api/verify-otp", controller.HandleValidateOtp)
	app.Post("/v1/api/post-new-password", controller.HandleUpdateUserPassword)
	app.Post("/v1/api/update-fcm-token", controller.HandleUpdateFcmToken)

	// Protected routes (require authentication)
	protected := app.Group("/v1/api", middleware.ExtractClaims(controller), middleware.RequireClaims())
	protected.Put("/users/:id", controller.HandleUpdateUser)
	protected.Post("/change-password", controller.HandleChangePassword)
	protected.Post("/save-address", controller.HandleSaveAddress)
	protected.Get("/get-addresses", controller.HandleGetAddresses)
	protected.Get("/users", controller.ListUsers)
	protected.Post("/update-profile", controller.HandleUpdateProfile)
	protected.Get("/users/:id", controller.FindUserById)
}

// Register implements the UserService interface
func (controller UserController) Register(ctx context.Context, model model.UserModel) entity.User {
	return controller.UserService.Register(ctx, model)
}

// HandleRegister handles the HTTP request for user registration
func (controller UserController) HandleRegister(c *fiber.Ctx) error {
	var request model.UserModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	if request.CountryCode == "NG" {
		request.AreaCode = "+234"
	} else {
		request.CountryCode = "+233"
	}
	request.Role = common.CUSTOMER_ROLE

	request.IsActive = false
	user := controller.Register(c.Context(), request)
	exception.PanicLogging(err)

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User created successfully",
		Data:    user,
	})
}

// RegisterCustomer implements the UserService interface
func (controller UserController) RegisterCustomer(ctx context.Context, request model.UserModel) interface{} {
	return controller.UserService.RegisterCustomer(ctx, request)
}

// HandleRegisterCustomer handles the HTTP request for customer registration
func (controller UserController) HandleRegisterCustomer(c *fiber.Ctx) error {
	var request model.UserModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	request.Role = common.CUSTOMER_ROLE
	user := controller.RegisterCustomer(c.Context(), request)

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User created successfully",
		Data:    user,
	})
}

// Authentication implements the UserService interface
func (controller UserController) Authentication(ctx context.Context, model model.LoginModel) entity.User {
	return controller.UserService.Authentication(ctx, model)
}

// HandleAuthentication handles the HTTP request for user authentication
func (controller UserController) HandleAuthentication(c *fiber.Ctx) error {
	var request model.LoginModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	result := controller.Authentication(c.Context(), request)
	var userRoles []map[string]interface{}

	tokenJwtResult := common.GenerateToken(result.Username, userRoles, result, controller.Config)
	resultWithToken := map[string]interface{}{
		"token":        tokenJwtResult,
		"username":     result.Username,
		"role":         result.UserRole,
		"emailAddress": result.Email,
		"fullName":     result.FirstName + " " + result.LastName,
		"firstName":    result.FirstName,
		"lastName":     result.LastName,
		"phoneNumber":  result.PhoneNumber,
		"id":           result.ID,
		"refineryId":   result.RefineryId,
		"countryCode":  result.CountryCode,
	}
	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    resultWithToken,
	})
}

// ChangePassword Authentication func ChangePassword user.
// @Description change password.
// @Summary  change password
// @Tags Authenticate user
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/change-pasword [post]
func (controller UserController) ChangePassword(ctx context.Context, token string, model model.ChangePasswordModel) entity.User {
	return controller.UserService.ChangePassword(ctx, token, model)
}

// HandleChangePassword handles the HTTP request for changing password
func (controller UserController) HandleChangePassword(c *fiber.Ctx) error {
	var request model.ChangePasswordModel
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Success: false,
		})
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	// Extract the token from header
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	// Remove Bearer prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Verify user is changing their own password
	if _, ok := claims["userId"].(float64); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID in token",
			Success: false,
		})
	}

	user := controller.ChangePassword(c.Context(), token, request)
	if user.Username == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to change password",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Password changed successfully",
		Success: true,
		Data:    user,
	})
}

// ListUsers Authentication func ListUsers .
// @Description list users.
// @Summary  list users
// @Tags Authenticate user
// @Produce json
// @Success 200 {object} model.GeneralResponse
// @Security JWT
// @Router /v1/api/users [get]
func (controller UserController) ListUsers(c *fiber.Ctx) error {
	var users []model.UserModel
	result, err := controller.UserService.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	// Transform the users
	for _, user := range result {
		users = append(users, model.UserModel{
			Id:           user.ID,
			Username:     user.Username,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			EmailAddress: user.Email,
			Role:         user.UserRole,
			PhoneNumber:  user.PhoneNumber,
			FileName:     user.FileName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Users retrieved successfully",
		Success: true,
		Data:    users,
	})
}

func (controller UserController) FindUserById(c *fiber.Ctx) error {
	var user model.UserModel
	userId, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)
	user, err = controller.UserService.FindByID(c.Context(), userId)
	exception.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(user)
}

// ValidateOtp implements the UserService interface
func (controller UserController) ValidateOtp(ctx context.Context, request model.OtpModel) entity.OneTimePassword {
	return controller.UserService.ValidateOtp(ctx, request)
}

// HandleValidateOtp handles the HTTP request for OTP validation
func (controller UserController) HandleValidateOtp(ctx *fiber.Ctx) error {
	var request model.OtpModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	otp := controller.ValidateOtp(ctx.Context(), request)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    otp,
	})
}

// ResetPassword implements the UserService interface
func (controller UserController) ResetPassword(ctx context.Context, request model.UserModel) model.UserModel {
	return controller.UserService.ResetPassword(ctx, request)
}

// HandleResetPassword handles the HTTP request for password reset
func (controller UserController) HandleResetPassword(ctx *fiber.Ctx) error {
	var request model.UserModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	user := controller.ResetPassword(ctx.Context(), request)
	user.Password = ""
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
		Success: true,
	})
}

// UpdateUserPassword implements the UserService interface
func (controller UserController) UpdateUserPassword(ctx context.Context, request model.ResetPasswordViewModel) model.UserModel {
	return controller.UserService.UpdateUserPassword(ctx, request)
}

// HandleUpdateUserPassword handles the HTTP request for updating user password
func (controller UserController) HandleUpdateUserPassword(ctx *fiber.Ctx) error {
	var request model.ResetPasswordViewModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	user := controller.UpdateUserPassword(ctx.Context(), request)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
		Success: true,
	})
}

// UpdateProfile implements the UserService interface
func (controller UserController) UpdateProfile(ctx context.Context, request model.UserModel, username string) (model.UserModel, error) {
	return controller.UserService.UpdateProfile(ctx, request, username)
}

// HandleUpdateProfile handles the HTTP request for profile updates
func (controller UserController) HandleUpdateProfile(ctx *fiber.Ctx) error {
	var request model.UserModel
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Success: false,
		})
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	// Use claims to get user ID
	if userId, ok := claims["userId"].(float64); ok {
		request.Id = uint(userId)
	}

	// Extract the token from header
	token := ctx.Get("Authorization")

	// Get id from params if provided - this will override the id from claims if present
	param := ctx.Params("id")
	if param != "" {
		id, err := strconv.Atoi(param)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Invalid user ID",
				Success: false,
			})
		}
		request.Id = uint(id)
	}

	user, err := controller.UpdateProfile(ctx.Context(), request, token)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "User profile updated successfully",
		Success: true,
		Data:    user,
	})
}

// SaveAddress implements the UserService interface
func (controller UserController) SaveAddress(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	return controller.UserService.SaveAddress(ctx, request)
}

// HandleSaveAddress handles the HTTP request for saving addresses
func (controller UserController) HandleSaveAddress(ctx *fiber.Ctx) error {
	var request map[string]interface{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Success: false,
		})
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "User ID not found in token",
			Success: false,
		})
	}

	request["userId"] = uint(userId)

	res, err := controller.SaveAddress(ctx.Context(), request)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Address saved successfully",
		Success: true,
		Data:    res,
	})
}

// GetAddresses implements the UserService interface
func (controller UserController) GetAddresses(ctx context.Context, id uint) (interface{}, error) {
	return controller.UserService.GetAddresses(ctx, id)
}

// HandleGetAddresses handles the HTTP request for getting user addresses
func (controller UserController) HandleGetAddresses(ctx *fiber.Ctx) error {
	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	userId, ok := claims["userId"].(float64)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "User ID not found in token",
			Success: false,
		})
	}

	addresses, err := controller.GetAddresses(ctx.Context(), uint(userId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Addresses retrieved successfully",
		Data:    addresses,
		Success: true,
	})
}

// UpdateFcmToken implements the UserService interface
func (controller UserController) UpdateFcmToken(ctx context.Context, request model.UpdateFcmToken) error {
	return controller.UserService.UpdateFcmToken(ctx, request)
}

// HandleUpdateFcmToken handles the HTTP request for FCM token updates
func (controller UserController) HandleUpdateFcmToken(ctx *fiber.Ctx) error {
	var request model.UpdateFcmToken
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	err = controller.UpdateFcmToken(ctx.Context(), request)
	exception.PanicLogging(err)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Successful",
		Data:    nil,
		Success: true,
	})
}

// HandleUpdateUser handles the HTTP request for updating a user
func (controller UserController) HandleUpdateUser(ctx *fiber.Ctx) error {
	var request model.UserModel
	err := ctx.BodyParser(&request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Success: false,
		})
	}

	// Get claims from context (set by middleware)
	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	// Verify if user has permission to update the requested user
	userId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Success: false,
		})
	}

	// Check if the authenticated user has permission to update this user
	if claimUserId, ok := claims["userId"].(float64); !ok || uint(claimUserId) != uint(userId) {
		// Check if user has admin role
		roles, ok := claims["roles"].([]interface{})
		if !ok {
			return ctx.Status(fiber.StatusForbidden).JSON(model.GeneralResponse{
				Code:    fiber.StatusForbidden,
				Message: "Access denied",
				Success: false,
			})
		}

		hasAdminRole := false
		for _, role := range roles {
			if roleMap, ok := role.(map[string]interface{}); ok {
				if roleMap["role"] == "ADMIN" {
					hasAdminRole = true
					break
				}
			}
		}

		if !hasAdminRole {
			return ctx.Status(fiber.StatusForbidden).JSON(model.GeneralResponse{
				Code:    fiber.StatusForbidden,
				Message: "Access denied",
				Success: false,
			})
		}
	}

	request.Id = uint(userId)

	// Extract the token from header
	claims, ok = middleware.GetClaims(ctx)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.GeneralResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "Authentication required",
			Success: false,
		})
	}

	user, err := controller.UpdateProfile(ctx.Context(), request, claims["username"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "User updated successfully",
		Success: true,
		Data:    user,
	})
}
