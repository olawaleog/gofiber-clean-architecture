package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type MessageController struct {
	service.MessageService
}

func (c MessageController) Route(app *fiber.App) {
	app.Get("/", c.WelcomeToAquaWizz)
	app.Get("/message-templates", c.FindAllMessageTemplates)
	app.Get("/send-sms", c.SendSms)
}

func NewMessageController(messageService *service.MessageService) *MessageController {
	return &MessageController{MessageService: *messageService}
}

func (c MessageController) FindAllMessageTemplates(ctx *fiber.Ctx) error {
	messageTemplates := c.MessageService.FindAllMessageTemplate(ctx.Context())
	return ctx.Status(fiber.StatusOK).JSON(messageTemplates)
}

func (c MessageController) FindMessageTemplateById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	exception.PanicLogging(err)
	messageTemplate := c.MessageService.FindMessageTemplateById(ctx.Context(), id)
	return ctx.Status(fiber.StatusOK).JSON(messageTemplate)
}

func (c MessageController) CreateMessageTemplate(ctx *fiber.Ctx) error {
	var messageTemlate model.MessageTemplateModel
	err := ctx.BodyParser(&messageTemlate)
	exception.PanicLogging(err)
	c.MessageService.CreateMessageTemplate(ctx.Context(), messageTemlate)
	return ctx.Status(fiber.StatusCreated).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    nil,
		Success: true,
	})
}

func (c MessageController) UpdateMessageTemplate(ctx *fiber.Ctx) error {
	var messageTemlate model.MessageTemplateModel
	err := ctx.BodyParser(&messageTemlate)
	exception.PanicLogging(err)
	c.MessageService.UpdateMessageTemplate(ctx.Context(), messageTemlate)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    nil,
		Success: true,
	})
}

func (c MessageController) WelcomeToAquaWizz(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Welcome to Aqua Wizz",
		Data:    "Welcome to Aqua Wizz",
		Success: true,
	})
}

func (c MessageController) SendSms(ctx *fiber.Ctx) error {
	data := model.SMSMessageModel{
		CountryCode: "+233",
		PhoneNumber: "0543798411",
		Message:     "Hello from Aqua Wizz",
	}
	c.MessageService.SendSMS(ctx.Context(), data)
	return ctx.Status(fiber.StatusOK).JSON(model.GeneralResponse{
		Code:    200,
		Message: "Success",
		Data:    nil,
		Success: true,
	})
}
