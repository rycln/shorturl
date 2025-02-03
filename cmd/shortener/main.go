package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	config "github.com/rycln/shorturl/configs"
	"github.com/rycln/shorturl/internal/app/handlers"
	"github.com/rycln/shorturl/internal/app/mem"
)

func main() {
	config.Init()

	store := mem.NewSimpleMemStorage()
	hv := handlers.NewHandlerVariables(store)

	webApp := fiber.New()
	webApp.Use(func(c *fiber.Ctx) error {
		c.Status(http.StatusBadRequest)
		return c.Next()
	})
	webApp.All("/", hv.ShortenURL)
	webApp.All("/:short", hv.ReturnURL)
	err := webApp.Listen(config.GetServerAddr())
	if err != nil {
		panic(err)
	}
}
