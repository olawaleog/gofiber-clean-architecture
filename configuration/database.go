package configuration

import (
	"fmt"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func NewDatabase(config Config) *gorm.DB {
	username := config.Get("DATASOURCE_USERNAME")
	password := config.Get("DATASOURCE_PASSWORD")
	host := config.Get("DATASOURCE_HOST")
	port := config.Get("DATASOURCE_PORT")
	dbName := config.Get("DATASOURCE_DB_NAME")
	maxPoolOpen, err := strconv.Atoi(config.Get("DATASOURCE_POOL_MAX_CONN"))
	maxPoolIdle, err := strconv.Atoi(config.Get("DATASOURCE_POOL_IDLE_CONN"))
	maxPollLifeTime, err := strconv.Atoi(config.Get("DATASOURCE_POOL_LIFE_TIME"))
	exception.PanicLogging(err)

	loggerDb := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	connectionString := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, username, password, dbName, port)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: loggerDb,
	})
	exception.PanicLogging(err)

	sqlDB, err := db.DB()
	exception.PanicLogging(err)

	sqlDB.SetMaxOpenConns(maxPoolOpen)
	sqlDB.SetMaxIdleConns(maxPoolIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(rand.Int31n(int32(maxPollLifeTime))) * time.Millisecond)

	err = db.AutoMigrate(&entity.Transaction{}, &entity.TransactionDetail{},
		&entity.User{}, &entity.UserRole{}, &entity.MessageTemplate{},
		&entity.LocalGovernmentArea{}, &entity.Address{}, &entity.Truck{},
		&entity.OneTimePassword{}, &entity.Refinery{}, &entity.PaymentMethod{},
		&entity.Order{},
	)
	//autoMigrate
	//err = db.AutoMigrate(&entity.Product{})
	//err = db.AutoMigrate(&entity.Transaction{})
	//err = db.AutoMigrate(&entity.TransactionDetail{})
	//err = db.AutoMigrate(&entity.User{})
	//err = db.AutoMigrate(&entity.UserRole{})
	exception.PanicLogging(err)
	return db
}
