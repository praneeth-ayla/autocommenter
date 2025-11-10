package scanner

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// skipDirs defines a map of directory names to be ignored during scanning.
var skipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"venv":         true,
	".next":        true,
	"build":        true,
}

// skipFilePatterns defines a list of file name patterns to be ignored during scanning.
var skipFilePatterns = []string{
	".env",
	".env.*",
}

// Scanner walks the specified path, collects file information, and skips unwanted directories and files.
func Scanner(path string) ([]FileInfo, error) {
	files := []FileInfo{}
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		// Skip unwanted directories if identified.
		if d.IsDir() && skipDirs[d.Name()] {
			return fs.SkipDir
		}

		// Process files (not directories).
		if !d.IsDir() {
			base := filepath.Base(path)
			// Skip files matching defined patterns.
			if shouldSkipFile(base) {
				return nil
			}

			// Get file info and count lines.
			info, _ := d.Info()
			lines := countLines(path)

			// Create FileInfo struct and append to results.
			file := FileInfo{
				Path:  path,
				Name:  base,
				Size:  info.Size(),
				Lines: lines,
			}

			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return files, errors.New("something went wrong")
	}

	return files, nil
}

// shouldSkipFile checks if a file name matches any of the skip patterns.
func shouldSkipFile(name string) bool {
	for _, pattern := range skipFilePatterns {
		match, _ := filepath.Match(pattern, name)

		if match {
			return true
		}
	}
	return false
}

// countLines counts the number of lines in a given file.
func countLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], []byte{'\n'})
		if err != nil {
			break
		}
	}
	return count
}
