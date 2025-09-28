package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type Map map[string]interface{}

type RaiseError struct {
	Message string
}

func (err RaiseError) Error() string {
	return err.Message
}

func ErrorResponse(logger *log.Logger, w http.ResponseWriter, status int, message []byte) {
	w.WriteHeader(status)
	if message != nil {
		w.Header().Set("Content-Type", "application/json")
		msg := map[string]string{"error": string(message)}
		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			logger.Printf("Error marshaling error message: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonMessage)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, http.StatusText(status), status)
	}

}

func ErrorResponseWithJSON(logger *log.Logger, w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Printf("Error marshaling error response data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		ErrorResponse(logger, w, status, jsonData)
	} else {
		http.Error(w, http.StatusText(status), status)
	}
}

func Response(logger *log.Logger, w http.ResponseWriter, status int, message []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	if message != nil {
		_, err := w.Write(message)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, http.StatusText(status), status)
	}
}

func ResponseWithJSON(logger *log.Logger, w http.ResponseWriter, status int, data interface{}) {

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Printf("Error marshaling response data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		Response(logger, w, status, jsonData, "application/json")
	}
}
