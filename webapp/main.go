package main

import (
	"embed"
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"golang.org/x/exp/slices"
)

//go:embed views/*
var viewsfs embed.FS

//go:embed 10-million-password-list-top-1000.txt
var passwordFile string
var passwordList = strings.Split(passwordFile, "\n")

// verifyPassword checks if the password passes OWASP Top 10 proactive controls c6: implement digital identity
func verifyPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// if password appears in the top 1000 passwords, it is not secure
	if slices.Contains(passwordList, password) {
		return errors.New("password is insecure")
	}

	return nil
}

type LoginRequest struct {
	Password string `form:"password"`
}

func runApp() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".gohtml")
	store := session.New()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/login", fiber.Map{})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Render("views/login", fiber.Map{
				"Error": err.Error(),
			})
		}

		if err := verifyPassword(req.Password); err != nil {
			return c.Render("views/login", fiber.Map{
				"Error": err.Error(),
			})
		}

		sess, err := store.Get(c)
		if err != nil {
			return c.Render("views/login", fiber.Map{
				"Error": err.Error(),
			})
		}

		sess.Set("password", req.Password)
		sess.Save()

		return c.Redirect("/welcome")
	})

	app.Get("/welcome", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Render("views/login", fiber.Map{
				"Error": err.Error(),
			})
		}

		return c.Render("views/welcome", fiber.Map{
			"Password": sess.Get("password"),
		})
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Render("views/login", fiber.Map{
				"Error": err.Error(),
			})
		}

		sess.Destroy()

		return c.Redirect("/")
	})

	app.Listen(":3000")
}

func main() {
	runApp()
}
