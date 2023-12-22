package http

import "github.com/gofiber/fiber/v2"

func MapLinksRoutes(router fiber.Router, h *LinksHandler) {
	router.Get("/get_short_url", h.GetOriginalUrl())
	router.Post("/create_short_url", h.CreateShortUrl())
}
