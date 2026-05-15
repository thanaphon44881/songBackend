package adapter

import (
	"music/service"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type userLibraryHttpsHandler struct {
	service service.ServiceuserLibrary
}

func NewuserLibraryHttps(service service.ServiceuserLibrary) *userLibraryHttpsHandler {
	return &userLibraryHttpsHandler{service: service}
}

func (u *userLibraryHttpsHandler) SaveLibrary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var body struct {
		SongID uint `json:"song_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	if body.SongID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "song_id is required",
		})
	}

	isLiked, err := u.service.Creat(userID, body.SongID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"is_liked": isLiked})
}

func (u *userLibraryHttpsHandler) DeleteUserLibrary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id, err := strconv.Atoi(c.Params("sonId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}

	if err = u.service.DeleteUser(id, int(userID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON("Delete This OK")
}

func (u *userLibraryHttpsHandler) ShowSong(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	id, err := strconv.Atoi(c.Params("songId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}
	library, err := u.service.FingUser(id, int(userID))

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(library)
}
