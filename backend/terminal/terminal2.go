package terminal

import (
	"log"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"github.com/gofiber/websocket/v2"
)

// TerminalHandler handles the WebSocket connection and terminal execution with command security
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

	// Goroutine to read from pty and send output to WebSocket
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

	// Buffer to store the user's command
	var commandBuffer strings.Builder

	// Main loop to read messages from WebSocket
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket closed:", err)
			return
		}

		// Append the message to the command buffer
		commandBuffer.WriteString(string(msg))

		// Check if the message contains a newline, indicating the command is complete
		if strings.Contains(commandBuffer.String(), "\n") {
			command := strings.TrimSpace(commandBuffer.String())

			// Check if the command is harmful before executing
			if isCommandHarmful(command) {
				warningMessage := "Warning: Command \"" + command + "\" is not allowed."
				if err := c.WriteMessage(websocket.TextMessage, []byte(warningMessage)); err != nil {
					log.Println("Error sending warning message:", err)
					return
				}
				// Reset the buffer and continue without executing
				commandBuffer.Reset()
				continue
			}

			// Write the full command to the pty
			if _, err := pty.Write([]byte(command + "\n")); err != nil {
				log.Println("Error writing to pty:", err)
				return
			}

			// Clear the command buffer after executing the command
			commandBuffer.Reset()
		}
	}
}

// isCommandHarmful checks if the command is harmful and should not be executed
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

	// Check if the command contains any harmful keywords
	for _, harmfulCmd := range harmfulCommands {
		if strings.Contains(command, harmfulCmd) {
			return true
		}
	}

	// Return false if no harmful command found
	return false
}
