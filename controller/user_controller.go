package controller

import (
	"fmt"
	"strconv"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
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
	app.Post("/v1/api/authentication", controller.Authentication)
	app.Post("/v1/api/register", controller.Register)
	app.Post("/v1/api/register-customer", controller.RegisterCustomer)
	app.Post("/v1/api/reset-password", controller.ResetPassword)
	app.Post("/v1/api/verify-otp", controller.ValidateOtp)
	app.Post("/v1/api/post-new-password", controller.UpdateUserPassword)

	// Protected routes (require authentication)
	protected := app.Group("/v1/api", middleware.ExtractClaims(controller.UserService), middleware.RequireClaims())
	protected.Put("/users/:id", controller.UpdateUser)
	protected.Post("/change-password", controller.ChangePassword)
	protected.Post("/save-address", controller.SaveAddress)
	protected.Get("/get-addresses", controller.GetAddresses)
	protected.Get("/users", controller.ListUsers)
	protected.Post("/update-profile", controller.UpdateProfile)
	protected.Post("/update-fcm-token", controller.UpdateFcmToken)
	protected.Get("/users/:id", controller.FindUserById)
}

// Register Registeration func Register user.
// @Description register new user.
// @Summary register user
// @Tags Authenticate user
// @Accept json
// @Produce json
// @Param request body model.UserModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Router /v1/api/register [post]
func (controller UserController) Register(c *fiber.Ctx) error {
	var request model.UserModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	if request.CountryCode == "NG" {
		request.AreaCode = "+234"
	} else {
		request.CountryCode = "+233"
	}

	exception.PanicLogging(err)

	request.IsActive = false
	_ = controller.UserService.Register(c.Context(), request)
	exception.PanicLogging(err)

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User created successfully",
		Data:    nil,
	})
}

func (controller UserController) RegisterCustomer(c *fiber.Ctx) error {
	var request model.UserModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	exception.PanicLogging(err)
	request.Role = common.CUSTOMER_ROLE
	user := controller.UserService.RegisterCustomer(c.Context(), request)

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User created successfully",
		Data:    user,
	})
}

// Authentication func Authenticate user.
// @Description authenticate user.
// @Summary authenticate user
// @Tags Authenticate user
// @Accept json
// @Produce json
// @Param request body model.LoginModel true "Request Body"
// @Success 200 {object} model.GeneralResponse
// @Router /v1/api/authentication [post]
func (controller UserController) Authentication(c *fiber.Ctx) error {
	var request model.LoginModel
	err := c.BodyParser(&request)
	exception.PanicLogging(err)

	result := controller.UserService.Authentication(c.Context(), request)
	var userRoles []map[string]interface{}
	//for _, userRole := range result.UserRole {
	//	userRoles = append(userRoles, map[string]interface{}{
	//		"role": userRole.Role,
	//	})
	//}
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
func (controller UserController) ChangePassword(c *fiber.Ctx) error {
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

	// Extract the token from header - we still need this for the service
	token := c.Get("Authorization")

	err = controller.UserService.ChangePassword(c.Context(), token, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GeneralResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Password changed successfully",
		Success: true,
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
	exception.PanicLogging(err)
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
	return c.Status(fiber.StatusOK).
		JSON(users)
}

func (u UserController) FindUserById(c *fiber.Ctx) error {
	var user model.UserModel
	userId, err := strconv.Atoi(c.Params("id"))
	exception.PanicLogging(err)
	user, err = u.UserService.FindByID(c.Context(), userId)
	exception.PanicLogging(err)
	return c.Status(fiber.StatusOK).JSON(user)
}

func (controller UserController) ValidateOtp(ctx *fiber.Ctx) error {
	var request model.OtpModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	otp := controller.UserService.ValidateOtp(ctx.Context(), request)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    otp,
	})
}

func (controller UserController) ResetPassword(ctx *fiber.Ctx) error {
	var request model.UserModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)

	user := controller.UserService.ResetPassword(ctx.Context(), request)
	user.Password = ""
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
		Success: true,
	})
}

func (controller UserController) UpdateUserPassword(ctx *fiber.Ctx) error {
	var request model.ResetPasswordViewModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)
	user := controller.UserService.UpdateUserPassword(ctx.Context(), request)

	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    user,
		Success: true,
	})
}

func (controller UserController) UpdateProfile(ctx *fiber.Ctx) error {
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

	// Extract the token from header - we still need this for the service
	token := ctx.Get("Authorization")

	// Get id from params if provided
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

	user, err := controller.UserService.UpdateProfile(ctx.Context(), request, token)
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

func (controller UserController) SaveAddress(ctx *fiber.Ctx) error {
	var request model.AddressModel
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

	request.UserId = uint(userId)

	res, err := controller.UserService.SaveAddress(ctx.Context(), request)
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

func (controller UserController) GetAddresses(ctx *fiber.Ctx) error {
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

	addresses, err := controller.UserService.GetAddresses(ctx.Context(), uint(userId))
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

func (controller UserController) UpdateFcmToken(ctx *fiber.Ctx) error {
	var request model.UpdateFcmToken
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)
	err = controller.UserService.UpdateFcmToken(ctx.Context(), request)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Successful",
		Data:    nil,
		Success: true,
	})
}

func (controller *UserController) UpdateUser(ctx *fiber.Ctx) error {
	var request model.UserModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)
	userId, err := strconv.Atoi(ctx.Params("id"))
	exception.PanicLogging(err)
	request.Id = uint(userId)
	token := ctx.Get("Authorization")
	user, err := controller.UserService.UpdateProfile(ctx.Context(), request, token)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.GeneralResponse{
			Code:    400,
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User updated successfully",
		Success: true,
		Data:    user,
	})
}
