package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/21TechLabs/factory-backend/utils"
	"github.com/go-playground/validator/v10"
)

const SchemaValidatorContextKey utils.ContextKey = "parsedBody"

func (m *Middleware) SchemaValidatorMiddleware(schemaFunc func() interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get content type
			contentType := r.Header.Get("Content-Type")

			body := schemaFunc()

			switch contentType {
			case "application/json":
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
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
				if err := json.Unmarshal(jsonData, &body); err != nil {
					utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte("Failed to unmarshal form data"))
					return
				}
			default:
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

			validate := validator.New()
			err := validate.Struct(body)

			if err != nil {
				utils.ErrorResponse(m.Logger, w, http.StatusBadRequest, []byte(err.Error()))
				return
			}

			// attach body to request context
			ctx := context.WithValue(r.Context(), SchemaValidatorContextKey, body)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
