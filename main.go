package main

import (
	"context"
	"embed"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/client/restclient"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/controller"
	_ "github.com/RizkiMufrizal/gofiber-clean-architecture/docs"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/jobs"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/middleware"
	repository "github.com/RizkiMufrizal/gofiber-clean-architecture/repository/impl"
	service "github.com/RizkiMufrizal/gofiber-clean-architecture/service/impl"
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
	//setup configuration
	//config := configuration.New(".env")
	logger.NewLogger()
	logger.Logger.Info("Starting the application...")

	//config := configuration.New("/var/www/api/.env")
	config := configuration.New(".env")

	database := configuration.NewDatabase(config)
	//redis := configuration.NewRedis(config)
	rabbitMQ := configuration.NewRabbitMQ(config)

	// Initialize Redis and RabbitMQ services
	//redisService := service.NewRedisService(redis)
	rabbitMQService := service.NewRabbitMQService(rabbitMQ, "aqua_wizz_exchange", "topic")

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

	//rest client
	httpRestClient := restclient.NewHttpRestClient(config)

	//service
	notificationService := service.NewNotificationService(config.Get("FCM_CREDENTIALS_PATH"))
	httpService := service.NewHttpBinServiceImpl(&httpRestClient)
	messageService := service.NewMessageServiceImpl(config, messageTemplateRepository, &httpService, rabbitMQService)
	// Initialize the email consumer service
	emailConsumerService := service.NewEmailConsumerService(config, rabbitMQService)
	transactionService := service.NewTransactionServiceImpl(&transactionRepository, &orderRepository, &paymentMethodRepository, &httpService, config, &notificationService)
	transactionDetailService := service.NewTransactionDetailServiceImpl(&transactionDetailRepository)
	localGovernmentService := service.NewLocalGovernmentServiceImpl(&localGovernmentRepository, config)
	userService := service.NewUserServiceImpl(&userRepository, &messageService, &localGovernmentService)
	truckService := service.NewTruckServiceImpl(&truckRepository, &userService, &messageService)
	refineryService := service.NewRefineryServiceImpl(&refineryRepository, &userService, &messageService, config)
	paymentService := service.NewPaymentService(&paymentRepository)

	// Google Maps Service
	mapsService := service.NewGoogleMapsService(config)

	//controller
	transactionController := controller.NewTransactionController(&transactionService, &userService, config)
	transactionDetailController := controller.NewTransactionDetailController(&transactionDetailService, config)
	userController := controller.NewUserController(&userService, config)
	messageController := controller.NewMessageController(&messageService)
	truckController := controller.NewTruckController(&truckService, config)
	refineryController := controller.NewRefineryController(&refineryService, config)
	paymentController := controller.NewPaymentController(&paymentService, &userService, config)
	localGovernmentController := controller.NewLocalGovernmentAreaController(&localGovernmentService)

	// Google Maps Controller
	mapsController := controller.NewGoogleMapsController(mapsService)

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
	localGovernmentController.Route(app)
	// Register Google Maps endpoints
	mapsController.RegisterRoutes(app)

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Set up cron jobs
	cronManager := jobs.SetupCronJobs(&transactionService, &truckService)
	cronManager.Start()
	defer cronManager.Stop()
	logger.Logger.Info("Cron jobs scheduled and started")

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

	// Properly close connections when app terminates
	defer func() {
		rabbitMQService.Close()
	}()

	//start app
	logger.Logger.Info("Application Started")

	userService.SeedUser(context.TODO())

	err = app.Listen(config.Get("SERVER.PORT"))
	exception.PanicLogging(err)
}
