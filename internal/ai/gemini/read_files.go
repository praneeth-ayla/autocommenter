package ai

import (
	"fmt"
	"os"
)

// FileContent holds the path and content of a file.
type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// ReadFileContent reads the content for a list of file paths.
func ReadFileContent(files []string) []FileContent {
	var fileContents []FileContent

	for _, filePath := range files {

		// Read the entire content of the file.
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		// Create a FileContent struct and append it to the slice.
		newFile := FileContent{
			Path:    filePath,
			Content: string(content),
		}
		fileContents = append(fileContents, newFile)
	}

	return fileContents
}
