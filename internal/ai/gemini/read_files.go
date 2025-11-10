package ai

import (
	"fmt"
	"os"
)

type FileContent struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func ReadFileContent(files FileResponse) []FileContent {
	var fileContents []FileContent

	for _, filePath := range files.Files {

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return nil
		}
		newFile := FileContent{
			Path:    filePath,
			Content: string(content),
		}
		fileContents = append(fileContents, newFile)
	}

	return fileContents
}
