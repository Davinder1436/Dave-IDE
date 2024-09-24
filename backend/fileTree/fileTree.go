package filetree

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// FileNode represents a file or directory in the file tree
type FileNode struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"` // "file" or "directory"
	Children []FileNode `json:"children,omitempty"`
}

// buildFileTree recursively builds the file tree structure starting from the provided root directory
func buildFileTree(root string) ([]FileNode, error) {
	var tree []FileNode
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		node := FileNode{Name: entry.Name()}
		if entry.IsDir() {
			node.Type = "directory"
			children, err := buildFileTree(filepath.Join(root, entry.Name()))
			if err != nil {
				return nil, err
			}
			node.Children = children
		} else {
			node.Type = "file"
		}
		tree = append(tree, node)
	}
	return tree, nil
}

// GetFileTreeHandler is a Fiber handler that returns the file tree as JSON
func GetFileTreeHandler(c *fiber.Ctx) error {
	root := "./User" // Define the root directory inside the User folder
	tree, err := buildFileTree(root)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(tree)
}
