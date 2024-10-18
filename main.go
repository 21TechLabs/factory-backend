package main

import (
	"log"
	"os"

	"github.com/21TechLabs/factory-be/database"
	"github.com/21TechLabs/factory-be/routes"
	"github.com/21TechLabs/factory-be/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	utils.LoadEnv()

	dbName := os.Getenv("MONGODB_DATABASE_NAME")

	if dbName == "" {
		log.Fatal("Please provide DB Name")
	}
	err := database.Connect(dbName)

	if err != nil {
		panic(err)
	}
}

func main() {
	var App *fiber.App = fiber.New()

	App.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetEnv("CORS_URLS", false),
		AllowCredentials: true,
	}))

	routes.SetupUser(App)

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Please provide PORT")
	}
	err := App.Listen(":" + PORT)

	if err != nil {
		log.Fatal(err)
	}

}
