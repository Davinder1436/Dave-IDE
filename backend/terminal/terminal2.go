package terminal

import (
	"log"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/gofiber/websocket/v2"
)

func TerminalHandler2(c *websocket.Conn) {
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

	var commandBuffer strings.Builder

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket closed:", err)
			return
		}

		commandBuffer.WriteString(string(msg))

		if strings.Contains(commandBuffer.String(), "\n") {
			command := strings.TrimSpace(commandBuffer.String())

			if isCommandHarmful(command) {
				warningMessage := "Warning: Command \"" + command + "\" is not allowed."
				if err := c.WriteMessage(websocket.TextMessage, []byte(warningMessage)); err != nil {
					log.Println("Error sending warning message:", err)
					return
				}

				commandBuffer.Reset()
				continue
			}

			if _, err := pty.Write([]byte(command + "\n")); err != nil {
				log.Println("Error writing to pty:", err)
				return
			}

			commandBuffer.Reset()
		}
	}
}

func isCommandHarmful(command string) bool {
	harmfulCommands := []string{
		"sudo",     // Superuser command
		"rm -rf",   // Forceful deletion
		"shutdown", // System shutdown
		"reboot",   // System reboot
		"poweroff", // System power off
		"passwd",   // Password change
		"chown",    // Permission change
		"mkfs",     // Disk formatting
		"dd",       // Disk overwriting
		"killall",  // Terminate all processes
		"kill -9",  // Force kill process
		"halt",     // Halt the system
		"init 0",   // Shut down the system
		"init 6",   // Reboot the system
	}

	for _, harmfulCmd := range harmfulCommands {
		if strings.Contains(command, harmfulCmd) {
			return true
		}
	}

	return false
}
