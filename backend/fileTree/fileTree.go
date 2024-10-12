package filetree

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type FileNode struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	Children []FileNode `json:"children,omitempty"`
}

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

func GetFileTreeHandler(c *fiber.Ctx) error {
	root := "./User"
	tree, err := buildFileTree(root)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(tree)
}
