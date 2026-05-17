package adapter

import (
	"fmt"
	"music/service"
	"music/slugs"
	"music/utils"
	"path/filepath"
	"strconv"

	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

type songHttpsHandler struct {
	service service.ServiceSong
}

func NewSongHttps(service service.ServiceSong) *songHttpsHandler {
	return &songHttpsHandler{service: service}
}

func (s *songHttpsHandler) SaveSong(c *fiber.Ctx) error {
	title := c.FormValue("title")
	category := c.FormValue("category")
	sub := c.FormValue("sub")
	country := c.FormValue("country")

	DurationStr := c.FormValue("duration")
	Duration, err := strconv.Atoi(DurationStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Duration must be number",
		})
	}
	artistStr := c.FormValue("artist_id")
	fmt.Println("art", artistStr)
	aitistid, err := strconv.Atoi(artistStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "aitistID must be number",
		})
	}

	// upStr := c.FormValue("up")
	// up, err := time.Parse("2006-01-02", upStr)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"error": "invalid date",
	// 	})
	// }

	fileCover, err := c.FormFile("cover")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image required"})
	}
	fileMusic, err := c.FormFile("song")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Song required"})
	}
	// filenameCover := fmt.Sprintf("%d_%s", time.Now().Unix(), fileCover.Filename)
	// pathCover := "./uploads/cover/" + filenameCover

	// pathMusic := "./uploads/music/" + fileMusic.Filename
	coverURL, songURL, err := utils.UploadSong(fileCover, fileMusic)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	songs := service.Song{
		Title:    title,
		Duration: Duration,
		Category: category,
		ArtistID: uint(aitistid),
		// Up:       up,
		File:    songURL,
		Cover:   coverURL,
		Sub:     sub,
		Country: country,
	}

	err = s.service.Creat(songs)

	// if err := c.SaveFile(fileCover, pathCover); err != nil {
	// 	return err
	// }

	// if err := c.SaveFile(fileMusic, pathMusic); err != nil {
	// 	return err
	// }
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Data recording completed",
		"Data":    songs,
	})
}

func (s *songHttpsHandler) ShowSong(c *fiber.Ctx) error {
	songs, err := s.service.Find()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	return c.Status(fiber.StatusBadRequest).JSON(songs)
}

func (s *songHttpsHandler) ShowNew(c *fiber.Ctx) error {
	songsnew, songstop, err := s.service.FindNew()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	return c.JSON(fiber.Map{
		"new": songsnew,
		"top": songstop,
	})
}

// func (s *songHttpsHandler) ShowSongID(c *fiber.Ctx) error {
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
// 	}
// 	song, err := s.service.First(id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(fiber.StatusBadRequest).JSON(song)
// }

func (s *songHttpsHandler) PlaySong(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	song, err := s.service.First(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	fileName := filepath.Base(song.File)
	path := "./uploads/music/" + fileName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.Status(404).JSON(fiber.Map{"error": "file not found"})
	}

	// go func() {

	// }()
	go s.service.UpdateId(id)

	c.Type("mp3")

	return c.SendFile(path)
}

func (h *songHttpsHandler) GetSongBySlug(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	slug := c.Params("slug")
	decodedSlug, _ := url.QueryUnescape(slug)

	userID := c.Locals("user_id").(uint)
	song, err := h.service.FirstLibrary(id, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "song not found"})
	}

	expectedSlug := slugs.GenerateSlug(song.Title)

	if decodedSlug != expectedSlug {
		return c.Redirect("/song/" + strconv.Itoa(id) + "/" + expectedSlug)
	}

	return c.JSON(song)
}

func (u *songHttpsHandler) CoverImage(c *fiber.Ctx) error {
	name := c.Params("name")
	idVal := c.Locals("user_id")

	_, ok := idVal.(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendFile("./uploads/cover/" + name)
}

func (u *songHttpsHandler) PlaySongNow(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	userID := c.Locals("user_id").(uint)

	data, err := u.service.PlayQueue(id, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	go u.service.UpdateId(id)
	return c.JSON(data)
}
