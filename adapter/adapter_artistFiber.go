package adapter

import (
	"fmt"
	"music/service"
	"strconv"
	"time"
	"net/url"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

type artistHttpsHandler struct {
	service service.ServiceArtist
}

func NewArtistHttps(s service.ServiceArtist) *artistHttpsHandler {
	return &artistHttpsHandler{service: s}
}

func (a *artistHttpsHandler) SaveArtist(c *fiber.Ctx) error {
	Name := c.FormValue("name")
	Bio := c.FormValue("bio")

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image required"})
	}
	filenames := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	path := "./uploads/artists/" + filenames

	artists := service.Artist{
		Name:  Name,
		Bio:   Bio,
		Image: "/artists/" + filenames,
	}

	err = a.service.Create(artists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = c.SaveFile(file, path)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(artists)
}

func (a *artistHttpsHandler) ShowArtist(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	artis, err := a.service.First(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(artis)
}

func (a *artistHttpsHandler) ShowArtistAll(c *fiber.Ctx) error {
	artis, err := a.service.Find()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid reauest"})
	}
	return c.JSON(artis)
}

func (a *artistHttpsHandler) ArtistsImage(c *fiber.Ctx) error {
	name := c.Params("name")

	decodedName, err := url.QueryUnescape(name)
	if err != nil {
		return c.Status(400).SendString("invalid filename")
	}

	idVal := c.Locals("user_id")
	_, ok := idVal.(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	filePath := filepath.Join("./uploads/artists/" + decodedName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(404).SendString("file not found")
	}

	return c.SendFile(filePath)
}
