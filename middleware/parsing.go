package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/utils"
)

func (m *Middleware) SchemaValidatorMiddleware(schemaKey dto.DtoMapKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get content type
			contentType := r.Header.Get("Content-Type")

			schemaFunc, ok := dto.DTOMap[schemaKey]

			if !ok {
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Schema not found"))
				return
			}

			body := schemaFunc()

			if body == nil {
				utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to create schema instance"))
				return
			}

			switch contentType {
			case "application/json":
				if err := json.NewDecoder(r.Body).Decode(body); err != nil {
					utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("Failed to parse JSON"))
					return
				}
			case "application/x-www-form-urlencoded":
				if err := r.ParseForm(); err != nil {
					utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("Failed to parse form"))
					return
				}
				jsonData, err := json.Marshal(r.Form)
				if err != nil {
					utils.ErrorResponse(m.Logger, w, http.StatusInternalServerError, []byte("Failed to marshal form data"))
					return
				}
				if err := json.Unmarshal(jsonData, body); err != nil {
					utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("Failed to unmarshal form data"))
					return
				}
			default:
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

			if err := utils.ValidateStruct(body); err != nil {
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte(err.Error()))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.SchemaValidatorContextKey, body)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
