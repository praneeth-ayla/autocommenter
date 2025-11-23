package scanner

import (
	"fmt"
	"os"
)

func Load(files []Info) []Data {

	filesContent := []Data{}
	for _, file := range files {
		fileContent, err := os.ReadFile(file.Path)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		filesContent = append(filesContent, Data{
			Path:    file.Path,
			Content: string(fileContent)},
		)
	}
	return filesContent

}

func LoadSingle(file Info) Data {

	fileContent, err := os.ReadFile(file.Path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Data{}
	}

	fileData := Data{
		Path:    file.Path,
		Content: string(fileContent),
	}

	return fileData
}
