package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/imilchev/rpi-feeder/pkg/mqtt/model"
	"github.com/imilchev/rpi-feeder/pkg/service/mqtt"
	"gorm.io/gorm"
)

type UsersController struct {
	//repo repos.UsersRepository
	mqtt mqtt.MqttManager
}

func NewUsersController(db *gorm.DB, mqtt mqtt.MqttManager) *UsersController {
	return &UsersController{mqtt: mqtt} //repo: repos.NewUsersRepository(db)}
}

func (c *UsersController) RegisterHandlers(a *fiber.App) {
	route := a.Group(apiGroup)
	route.Post("/feed", c.FeedPortions)
}

func (c *UsersController) FeedPortions(ctx *fiber.Ctx) error {
	msg := model.FeedMessage{
		Portions: 1,
	}
	if err := c.mqtt.SendFeedCommand("rpidev", msg); err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"test": "ok"})
}
