package main

import (
	"log"

	filetree "github.com/Davinder1436/Dave-IDE/fileTree"
	"github.com/Davinder1436/Dave-IDE/terminal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/terminal", websocket.New(terminal.TerminalHandler))
	app.Get("/filetree", filetree.GetFileTreeHandler)

	log.Fatal(app.Listen(":3000"))
}
