package app

import (
	"fmt"
	"log"
	"os"

	"github.com/21TechLabs/factory-backend/controllers"
	oauth_controller "github.com/21TechLabs/factory-backend/controllers/oauth"
	payments_controller "github.com/21TechLabs/factory-backend/controllers/payments"
	"github.com/21TechLabs/factory-backend/database"
	"github.com/21TechLabs/factory-backend/middleware"
	"github.com/21TechLabs/factory-backend/models"
	"github.com/21TechLabs/factory-backend/utils"
	"gorm.io/gorm"
)

type Application struct {
	Logger                *log.Logger
	DB                    *gorm.DB
	Middleware            *middleware.Middleware
	UserController        *controllers.UserController
	FileController        *controllers.FileController
	OAuthController       *oauth_controller.OAuthController
	HealthCheckController *controllers.HealthCheckController
	PaymentPlanController *payments_controller.PaymentPlanController
}

func NewApplication() (*Application, error) {
	if err := utils.LoadEnv(); err != nil {
		return nil, err
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbCredentials := database.DatabaseCredentials{
		Username:    utils.GetEnv("DB_USERNAME", true),
		Password:    utils.GetEnv("DB_PASSWORD", true),
		Database:    utils.GetEnv("DB_NAME", true),
		Host:        utils.GetEnv("DB_HOST", true),
		Port:        utils.GetEnv("DB_PORT", true),
		SSLDisabled: utils.GetEnv("DB_SSL_DISABLED", true),
		TimeZone:    utils.GetEnv("DB_TIMEZONE", true),
	}

	db, err := database.Open(&dbCredentials)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	if db == nil {
		logger.Fatal("Database connection is nil")
	}
	logger.Println("✅ Database connection established successfully")

	var modelsToMigrate = []interface{}{
		models.User{},
		models.File{},
		models.ProductPlan{},
		models.Transaction{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("failed to auto migrate model %T: %w", model, err)
		}
		logger.Printf("✅ Model %T migrated successfully", model)
	}

	// store initialization
	fileStore := models.NewFileStore(db)
	userStore := models.NewUserStore(db, fileStore)
	paymentPlanStore := models.NewProductPlanStore(db, userStore)

	// middleware initialization
	middleware := middleware.NewMiddleware(logger, userStore)

	// controller initialization
	userController := controllers.NewUserController(logger, userStore)
	fileController := controllers.NewFileController(logger, fileStore, userStore)
	oauthController := oauth_controller.NewOAuthController(logger, userStore)
	healthCheckController := controllers.NewHealthCheckController(logger)
	paymentPlanController := payments_controller.NewPaymentPlanController(logger, paymentPlanStore)

	app := &Application{
		Logger:                logger,
		DB:                    db,
		Middleware:            middleware,
		UserController:        userController,
		FileController:        fileController,
		OAuthController:       oauthController,
		HealthCheckController: healthCheckController,
		PaymentPlanController: paymentPlanController,
	}

	return app, nil

}
