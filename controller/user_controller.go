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
	app.Post("/v1/api/authentication", controller.Authentication)
	app.Post("/v1/api/register", controller.Register)
	app.Put("/v1/api/users/:id", controller.UpdateUser)
	app.Post("/v1/api/register-customer", controller.RegisterCustomer)
	app.Post("/v1/api/change-password", controller.ChangePassword)
	app.Post("/v1/api/reset-password", controller.ResetPassword)
	app.Post("/v1/api/save-address", controller.SaveAddress)
	app.Get("/v1/api/get-addresses", controller.GetAddresses)
	app.Get("/v1/api/users", controller.ListUsers)
	app.Post("/v1/api/verify-otp", controller.ValidateOtp)
	app.Post("/v1/api/post-new-password", controller.UpdateUserPassword)
	app.Post("/v1/api/update-profile", controller.UpdateProfile)
	app.Post("/v1/api/update-fcm-token", controller.UpdateFcmToken)

	app.Get("/v1/api/users/:id", controller.FindUserById)
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
	exception.PanicLogging(err)

	token := c.Get("Authorization")

	_ = controller.UserService.ChangePassword(c.Context(), token, request)
	//for _, userRole := range result.UserRole {
	//	userRoles = append(userRoles, map[string]interface{}{
	//		"role": userRole.Role,
	//	})
	//}

	return c.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
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
	exception.PanicLogging(err)
	token := ctx.Get("Authorization")
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	request.Id = uint(id)
	user, err := controller.UserService.UpdateProfile(ctx.Context(), request, token)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User updated successfully",
		Success: true,
		Data:    user,
	})
}

func (controller UserController) SaveAddress(ctx *fiber.Ctx) error {
	var request model.AddressModel
	err := ctx.BodyParser(&request)
	exception.PanicLogging(err)
	token := ctx.Get("Authorization")
	claims, err := controller.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	request.UserId = uint(claims["userId"].(float64))
	fmt.Println(request)
	res, err := controller.UserService.SaveAddress(ctx.Context(), request)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "User updated successfully",
		Success: true,
		Data:    res,
	})
}

func (controller UserController) GetAddresses(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	claims, err := controller.UserService.GetClaimsFromToken(ctx.Context(), token)
	exception.PanicLogging(err)
	userId := uint(claims["userId"].(float64))
	addresses, err := controller.UserService.GetAddresses(ctx.Context(), userId)
	exception.PanicLogging(err)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Successful",
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
