package ai

import (
	"fmt"

	"github.com/praneeth-ayla/AutoCommenter/internal/prompt"
)

func BuildAnalyzeFilesForCommentsPrompt(files []string) string {
	return fmt.Sprintf(prompt.AnalyzeFilesForCommentsResponse, files)
}
