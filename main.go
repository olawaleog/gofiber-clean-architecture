package main

import (
	"context"
	"embed"
	"encoding/json"
	"os"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/client/restclient"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/controller"
	_ "github.com/RizkiMufrizal/gofiber-clean-architecture/docs"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/middleware"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	repository "github.com/RizkiMufrizal/gofiber-clean-architecture/repository/impl"
	service "github.com/RizkiMufrizal/gofiber-clean-architecture/service/impl"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

var f embed.FS

// @title Go Fiber Clean Architecture
// @version 1.0.0
// @description Baseline project using Go Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9999
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Authorization For JWT
func main() {
	configPath := ""
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	//setup configuration
	//config := configuration.New(".env")
	logger.NewLogger()
	logger.Logger.Info("Starting the application...")

	config := configuration.New(configPath + ".env")
	//config := configuration.New(".env")

	database := configuration.NewDatabase(config)
	redis := configuration.NewRedis(config)
	rabbitMQ := configuration.NewRabbitMQ(config)

	// Initialize Redis and RabbitMQ services
	rabbitMQService := service.NewRabbitMQService(rabbitMQ, "aqua_wizz_exchange", "topic")
	// Redis service initialization if needed in the future
	redisService := service.NewRedisService(redis)

	//repository
	messageTemplateRepository := repository.NewMessageTemplateRepositoryImpl(database)
	transactionRepository := repository.NewTransactionRepositoryImpl(database)
	transactionDetailRepository := repository.NewTransactionDetailRepositoryImpl(database)
	userRepository := repository.NewUserRepositoryImpl(database)
	truckRepository := repository.NewTruckRepositoryImpl(database)
	refineryRepository := repository.NewRefineryRepositoryImpl(database)
	paymentRepository := repository.NewPaymentRepository(database)
	localGovernmentRepository := repository.NewLocalGovernmentRepository(database)
	orderRepository := repository.NewOrderRepository(database)
	paymentMethodRepository := repository.NewPaymentMethodRepository(database)
	paymentConfigRepository := repository.NewPaymentConfigRepository(database, redis)
	notificationRepository := repository.NewNotificationRepository(database)
	settingRepository := repository.NewSettingRepository(database)

	//rest client
	httpRestClient := restclient.NewHttpRestClient(config)

	//service
	notificationService := service.NewNotificationService(config.Get("FCM_CREDENTIALS_PATH"))
	httpService := service.NewHttpBinServiceImpl(&httpRestClient)
	messageService := service.NewMessageServiceImpl(config, messageTemplateRepository, &httpService, rabbitMQService, notificationRepository)
	// Initialize the email consumer service
	emailConsumerService := service.NewEmailConsumerService(config, rabbitMQService)
	transactionService := service.NewTransactionServiceImpl(&transactionRepository, &orderRepository, &paymentMethodRepository, &httpService, config, &notificationService, &paymentConfigRepository)
	transactionDetailService := service.NewTransactionDetailServiceImpl(&transactionDetailRepository)
	localGovernmentService := service.NewLocalGovernmentServiceImpl(&localGovernmentRepository, config)
	userService := service.NewUserServiceImpl(&userRepository, &messageService, &localGovernmentService, config)
	truckService := service.NewTruckServiceImpl(&truckRepository, &userService, &messageService)
	settingService := service.NewSettingService(settingRepository, redisService, validator.New())
	refineryService := service.NewRefineryServiceImpl(&refineryRepository, &userService, &messageService, &settingService, config)
	paymentService := service.NewPaymentService(&paymentRepository)
	paymentConfigService := service.NewPaymentConfigService(paymentConfigRepository)

	// Google Maps Service
	mapsService := service.NewGoogleMapsService(config)

	// Authorization Service
	authorizationService := service.NewAuthorizationService()

	//controller
	transactionController := controller.NewTransactionController(&transactionService, &userService, &authorizationService, config)
	transactionDetailController := controller.NewTransactionDetailController(&transactionDetailService, config)
	userController := controller.NewUserController(&userService, config)
	messageController := controller.NewMessageController(&messageService)
	truckController := controller.NewTruckController(&truckService, config)
	refineryController := controller.NewRefineryController(&refineryService, config)
	paymentController := controller.NewPaymentController(&paymentService, &userService, config)
	paymentConfigController := controller.NewPaymentConfigController(paymentConfigService)
	localGovernmentController := controller.NewLocalGovernmentAreaController(&localGovernmentService)
	settingController := controller.NewSettingController(settingService)

	// Google Maps Controller
	mapsController := controller.NewGoogleMapsController(mapsService)

	// Notification Controller
	notificationController := controller.NewNotificationController(notificationService, userRepository)

	//setup fiber
	app := fiber.New(configuration.NewFiberConfiguration())

	app.Use(recover.New())
	app.Use(cors.New())
	app.Static("/uploads", "./uploads")
	app.Use(middleware.RequestLogger)

	//routing
	transactionController.Route(app)
	transactionDetailController.Route(app)
	userController.Route(app)
	messageController.Route(app)
	truckController.Route(app)
	refineryController.Route(app)
	paymentController.Route(app)
	paymentConfigController.Route(app)
	localGovernmentController.Route(app)
	// Register Google Maps endpoints
	mapsController.RegisterRoutes(app)
	// Register setting routes
	settingController.Route(app)
	// Register notification routes
	notificationController.Route(app)

	// Payment configuration routes are registered through the controller's Route method

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Set up cron jobs
	//cronManager := jobs.SetupCronJobs(&transactionService, &truckService)
	//cronManager.Start()
	//defer cronManager.Stop()
	//logger.Logger.Info("Cron jobs scheduled and started")

	// Start the email consumer service
	err := emailConsumerService.StartConsumer()
	if err != nil {
		logger.Logger.Error("Failed to start email consumer service: " + err.Error())
	}

	// Set up example RabbitMQ subscription
	err = rabbitMQService.SubscribeToTopic("notifications", func(message []byte) error {
		logger.Logger.Info("Received notification message: " + string(message))
		return nil
	})
	if err != nil {
		logger.Logger.Error("Failed to subscribe to notifications topic: " + err.Error())
	}

	// Start the SMS consumer service
	err = rabbitMQService.SubscribeToTopic("sms.send", func(message []byte) error {
		var sms model.SMSMessageModel
		if err := json.Unmarshal(message, &sms); err != nil {
			logger.Logger.Error("Failed to unmarshal SMS message: " + err.Error())
			return err
		}
		if err := messageService.SendSMSDirect(sms); err != nil {
			logger.Logger.Error("Failed to send SMS: " + err.Error())
			return nil
		}
		logger.Logger.Info("SMS sent to: " + sms.PhoneNumber)
		return nil
	})
	if err != nil {
		logger.Logger.Error("Failed to subscribe to sms.send topic: " + err.Error())
	}

	// Properly close connections when app terminates
	defer func() {
		rabbitMQService.Close()
	}()
	//err = messageService.SendSMSDirect(model.SMSMessageModel{
	//	CountryCode: "+233",
	//	PhoneNumber: "0243101864",
	//	Message:     "Hello world",
	//})
	//start app
	logger.Logger.Info("Application Started")

	userService.SeedUser(context.TODO())

	err = app.Listen(config.Get("SERVER.PORT"))
	exception.PanicLogging(err)
}
