package scanner

import (
	"fmt"
	"os"
)

func WriteFile(filePath, content string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
