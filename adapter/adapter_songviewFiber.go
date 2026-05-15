package adapter

import (
	"music/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type songViewHttpsHandler struct {
	service service.ServiceSongView
}

func NewSonViewgHttps(service service.ServiceSongView) *songViewHttpsHandler {
	return &songViewHttpsHandler{service: service}
}

func (h *songViewHttpsHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var body struct {
		SongID uint `json:"song_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	err := h.service.Creat(userID, body.SongID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "saved"})
}

func (h *songViewHttpsHandler) GetMyHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	data, err := h.service.GetByUser(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

func (h *songViewHttpsHandler) GetMyHistoryAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	data, err := h.service.GetByUserAll(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

func (h *songViewHttpsHandler) DeleteSongView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.service.Delete(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON("Delete This OK")
}
