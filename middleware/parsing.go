package middleware

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func SchemaValidatorMiddleware(schemaFunc func() interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Body()
		schema := schemaFunc()

		err := json.Unmarshal(body, &schema)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		validate := validator.New()

		err = validate.Struct(schema)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		return c.Next()
	}
}
