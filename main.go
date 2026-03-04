package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/routes"
	"github.com/21TechLabs/factory-backend/utils"

	"github.com/rs/cors"
)

func init() {
	// docs
	utils.LoadEnv()
}

func main() {

	app, err := app.NewApplication()

	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Please provide PORT")
	}

	router := routes.SetupRoutes(app)

	corsUrls := strings.Split(utils.GetEnv("CORS_URLS", false), ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   corsUrls,
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: handler,
	}
	app.Logger.Printf("✅ We are running on port %s\n", PORT)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
