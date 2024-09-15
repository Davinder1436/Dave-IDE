package main

import (
	"log"

	"os/exec"

	"github.com/creack/pty"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	// Serve static files for frontend
	app.Static("/", "./public")

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Start a bash process
		cmd := exec.Command("bash")
		pty, err := pty.Start(cmd)
		if err != nil {
			log.Printf("Failed to start pty: %v", err)
			return
		}
		defer func() {
			_ = pty.Close()
			_ = cmd.Wait()
		}()

		// Read data from the pty and send to frontend via WebSocket
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := pty.Read(buf)
				if err != nil {
					log.Println("Error reading from pty:", err)
					return
				}
				if err := c.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					log.Println("Error sending message:", err)
					return
				}
			}
		}()

		// Read data from WebSocket and write to the pty (terminal)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("WebSocket closed:", err)
				return
			}
			if _, err := pty.Write(msg); err != nil {
				log.Println("Error writing to pty:", err)
				return
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
