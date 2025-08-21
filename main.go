package main

/**
 * Notes of Abbreviations & Terms:
 *
 * OFF | OpenFoodFacts - Main source for UPC data & Food Data
 */

import (
	shiutils "CrystalENVT/SimpleHomeInventory/SHI_Utils"

	"github.com/gofiber/fiber/v2"         // web server
	"github.com/gofiber/template/html/v2" // web server
)

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Serve static files (HTML templates and stylesheets).
	app.Static("/", "./static")

	// Define routes.
	app.Get("/", shiutils.RenderForm)
	app.Post("/submit", shiutils.ProcessForm)

	// Start the Fiber app on port 7070.
	app.Listen(":7070")
}
