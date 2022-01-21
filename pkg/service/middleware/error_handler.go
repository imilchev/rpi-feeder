package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/service/models"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's an *models.ApiError
	if e, ok := err.(*models.ApiError); ok {
		code = e.Code()
	} else {
		err = models.NewApiError(code, err.Error())
	}

	// Send custom error page
	return ctx.Status(code).JSON(err)
}
