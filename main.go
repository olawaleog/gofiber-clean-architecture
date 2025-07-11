package main

import (
	"context"
	"embed"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/client/restclient"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/controller"
	_ "github.com/RizkiMufrizal/gofiber-clean-architecture/docs"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/jobs"
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
	common.NewLogger()
	common.Logger.Info("Starting the application...")

	//config := configuration.New("/var/www/api/.env")
	config := configuration.New(".env")

	database := configuration.NewDatabase(config)
	//redis := configuration.NewRedis(config)
	//e := event.New()

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
	httpService := service.NewHttpBinServiceImpl(&httpRestClient)
	messageService := service.NewMessageServiceImpl(config, messageTemplateRepository, &httpService)
	transactionService := service.NewTransactionServiceImpl(&transactionRepository, &orderRepository, &paymentMethodRepository, &httpService, config)
	transactionDetailService := service.NewTransactionDetailServiceImpl(&transactionDetailRepository)
	localGovernmentService := service.NewLocalGovernmentServiceImpl(&localGovernmentRepository, config)
	userService := service.NewUserServiceImpl(&userRepository, &messageService, &localGovernmentService)
	truckService := service.NewTruckServiceImpl(&truckRepository, &userService, &messageService)
	refineryService := service.NewRefineryServiceImpl(&refineryRepository, &userService, &messageService)
	paymentService := service.NewPaymentService(&paymentRepository)

	// Initialize FCM Notification Service
	//notificationService := service.NewNotificationService(config.Get("FCM_CREDENTIALS_PATH"))

	//controller
	transactionController := controller.NewTransactionController(&transactionService, &userService, config)
	transactionDetailController := controller.NewTransactionDetailController(&transactionDetailService, config)
	userController := controller.NewUserController(&userService, config)
	messageController := controller.NewMessageController(&messageService)
	truckController := controller.NewTruckController(&truckService, config)
	refineryController := controller.NewRefineryController(&refineryService, config)
	paymentController := controller.NewPaymentController(&paymentService, &userService, config)
	localGovernmentController := controller.NewLocalGovernmentAreaController(&localGovernmentService)

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

	//swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Set up cron jobs
	cronManager := jobs.SetupCronJobs(&transactionService, &truckService)
	cronManager.Start()
	defer cronManager.Stop()
	common.Logger.Info("Cron jobs scheduled and started")
	//start app
	common.Logger.Info("Application Started")

	userService.SeedUser(context.TODO())

	err := app.Listen(config.Get("SERVER.PORT"))
	exception.PanicLogging(err)
}
