package controllers

import (
	"log"
	"net/http"

	"github.com/21TechLabs/factory-backend/utils"
)

type HealthCheckController struct {
	Logger *log.Logger
}

func NewHealthCheckController(logger *log.Logger) *HealthCheckController {
	return &HealthCheckController{
		Logger: logger,
	}
}

func (hc HealthCheckController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.Response(hc.Logger, w, http.StatusOK, []byte("{\"status\": \"Ok!\"}"), "application/json")
}
