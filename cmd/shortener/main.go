package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/handlers"
	"github.com/rycln/shorturl/internal/app/mem"
)

func main() {
	store := mem.NewSimpleMemStorage()
	hv := handlers.NewHandlerVariables(store)

	webApp := fiber.New()
	webApp.Use(func(c *fiber.Ctx) error {
		c.Status(http.StatusBadRequest)
		return c.Next()
	})
	webApp.All("/", hv.ShortenURL)
	webApp.All("/:short", hv.ReturnURL)
	err := webApp.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
