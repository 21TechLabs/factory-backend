package utils

import "github.com/gofiber/fiber/v2"

type RaiseError struct {
	Message string
}

func (err RaiseError) Error() string {
	return err.Message
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error":   message,
		"success": false,
	})
}
