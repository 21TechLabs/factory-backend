package middleware

import (
	"github.com/21TechLabs/musiclms-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) SchemaValidatorMiddleware(schemaFunc func() interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {

		schema := schemaFunc()
		var err error

		if err = c.BodyParser(schema); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body: "+err.Error())
		}

		validate := validator.New()

		err = validate.Struct(schema)

		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed: "+err.Error())
		}

		return c.Next()
	}
}
