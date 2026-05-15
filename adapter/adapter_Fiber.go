package adapter

import (
	"music/service"
	"music/utils"

	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	// "golang.org/x/crypto/bcrypt"
	"os"

	"github.com/joho/godotenv"
)

type userHttpsHandler struct {
	services service.ServiceUser
}

func NewUserHttps(service service.ServiceUser) *userHttpsHandler {
	return &userHttpsHandler{services: service}
}

func (s *userHttpsHandler) CreatUser(c *fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	bio := c.FormValue("biaoname")

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image required"})
	}
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	path := "./uploads/images/" + filename

	if err := c.SaveFile(file, path); err != nil {
		return err
	}

	user := service.User{
		Username: username,
		Email:    email,
		Password: password,
		Biaoname: bio,
		Image:    "/images/" + filename,
	}
	err = s.services.Creat(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *userHttpsHandler) Login(c *fiber.Ctx) error {
	var req struct {
		ID       uint   `gorm:"primaryKey" json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	user, err := h.services.FindByUsername(req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// if user.Password != req.Password {
	// 	return c.Status(401).JSON(fiber.Map{
	// 		"error": "wrong password",
	// 	})
	// }
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return err
	}

	godotenv.Load()
	Name := os.Getenv("NCOOKIE")
	c.Cookie(&fiber.Cookie{
		Name:     Name,
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",

		Expires: time.Now().Add(72 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (u *userHttpsHandler) ShowUsserID(c *fiber.Ctx) error {
	idVal := c.Locals("user_id")
	id, ok := idVal.(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := u.services.First(int(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(user)
}

func (u *userHttpsHandler) ImageShow(c *fiber.Ctx) error {
	name := c.Params("name")
	idVal := c.Locals("user_id")

	_, ok := idVal.(uint)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.SendFile("./uploads/images/" + name)
}
