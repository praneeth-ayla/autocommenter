package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

var skipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	".next":        true,
	"build":        true,
	"dist":         true,
	"migrations":   true,
	"prisma":       true,
}

var allowedExt = map[string]bool{
	".ts":  true,
	".tsx": true,
	".js":  true,
	".jsx": true,
	".go":  true,
	".py":  true,
}

func Scanner(root string) ([]string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var files []string
	err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Walk error:", err)
			return nil
		}

		// Skip directories by name
		if d.IsDir() {
			rel, _ := filepath.Rel(abs, path)
			for skip := range skipDirs {
				if strings.HasPrefix(rel, skip+"/") {
					return filepath.SkipDir
				}
			}
		}

		// Accept only files
		if !d.IsDir() {
			ext := filepath.Ext(path)
			// Skip files with non-allowed extensions
			if !allowedExt[ext] {
				return nil
			}

			// Convert to clean relative path
			rel, err := filepath.Rel(abs, path)
			if err == nil {
				files = append(files, rel)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort for consistent output every run
	sort.Strings(files)
	return files, nil
}

// FilterFilesNeedingComments removes files that typically don't need comments
func FilterFilesNeedingComments(files []string) []string {
	var result []string

	for _, file := range files {
		// Skip type definition files
		if strings.HasSuffix(file, ".d.ts") {
			continue
		}

		// Skip specific directories that don't need comments
		if strings.Contains(file, "/ui/") || // shadcn/ui components
			strings.Contains(file, "/types/") || // type definition files
			strings.Contains(file, "/__tests__/") || // test files
			strings.Contains(file, "/.storybook/") { // storybook config
			continue
		}

		// Skip specific config/generated files
		if strings.HasPrefix(file, "next-env.d.ts") ||
			strings.HasPrefix(file, "next.config") ||
			strings.Contains(file, "seed.") { // seed files are usually simple
			continue
		}

		// Include everything else
		result = append(result, file)
	}

	return result
}
