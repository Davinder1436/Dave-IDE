package main

import (
	"log"

	"github.com/Davinder1436/Dave-IDE/terminal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// Serve static files
	app.Static("/", "./public")

	// WebSocket route to handle terminal requests
	app.Get("/terminal", websocket.New(terminal.TerminalHandler))

	// Start server on port 3000
	log.Fatal(app.Listen(":3000"))
}
