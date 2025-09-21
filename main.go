package main

import (
	"log"
	"os"

	"github.com/21TechLabs/musiclms-backend/app"
	"github.com/21TechLabs/musiclms-backend/routes"
	"github.com/21TechLabs/musiclms-backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	utils.LoadEnv()
}

func main() {

	app, err := app.NewApplication()

	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	var FiberApp = fiber.New()

	FiberApp.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			whitelistOrigins := utils.GetEnv("CORS_URLS", false)
			return utils.IsValidOrigin(origin, whitelistOrigins)
		},
	}))

	routes.SetupRoutes(FiberApp, app)

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Please provide PORT")
	}

	app.Logger.Printf("✅ We are running on port %s\n", PORT)
	err = FiberApp.Listen(":" + PORT)

	if err != nil {
		log.Fatal(err)
	}

}
