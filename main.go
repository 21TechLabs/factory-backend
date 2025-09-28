package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/21TechLabs/factory-backend/app"
	"github.com/21TechLabs/factory-backend/routes"
	"github.com/21TechLabs/factory-backend/utils"
)

func init() {
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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: router,
	}
	app.Logger.Printf("✅ We are running on port %s\n", PORT)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
