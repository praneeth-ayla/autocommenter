package contextstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/praneeth-ayla/autocommenter/internal/scanner"
)

// Save stores the provided FileDetails map to a JSON file.
func Save(all map[string]FileDetails) error {
	configPath, err := getConfigFilePath() // Determine the path for the configuration file.
	if err != nil {
		return err
	}
	fmt.Println(configPath) // Log the configuration file path.

	err = ensureDir(filepath.Dir(configPath)) // Ensure the directory for the config file exists.
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(all, "", "  ") // Marshal the map into a pretty-printed JSON byte slice.
	if err != nil {
		return err
	}

	return scanner.WriteFile(configPath, string(data)) // Write the JSON data to the file using scanner.WriteFile.
}

// Load retrieves FileDetails from the JSON configuration file.
func Load() (map[string]FileDetails, error) {
	configPath, err := getConfigFilePath() // Get the configuration file path.
	if err != nil {
		return nil, err
	}

	// Check if the configuration file exists, return an error if it doesn't.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("context file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath) // Read the content of the configuration file.
	if err != nil {
		return nil, err
	}

	var all map[string]FileDetails
	err = json.Unmarshal(data, &all) // Unmarshal the JSON data into the 'all' map.
	if err != nil {
		return nil, err
	}

	return all, nil
}

// getConfigFilePath determines the full path to the context JSON file.
func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir() // Get the user's home directory.
	if err != nil {
		return "", err
	}

	projectRoot := scanner.GetProjectRoot()           // Get the project's root directory.
	goModPath := filepath.Join(projectRoot, "go.mod") // Construct the path to the go.mod file.

	data, err := os.ReadFile(goModPath) // Read the content of the go.mod file.
	if err != nil {
		return "", err
	}

	var moduleName string
	// Parse the go.mod file to extract the module name.
	_, err = fmt.Sscanf(string(data), "module %s", &moduleName)
	if err != nil {
		return "", fmt.Errorf("failed to read module name")
	}

	projectName := filepath.Base(moduleName) // Extract the project name from the module name.

	fileName := projectName + ".json" // Construct the filename for the context file.

	// Combine home directory, a fixed directory name, and the generated filename.
	return filepath.Join(home, "AutoCommenter", fileName), nil
}

// ensureDir creates a directory and any necessary parent directories if they don't exist.
func ensureDir(dir string) error {
	// os.MkdirAll creates the directory named dir, along with any necessary parents.
	// The permission bits are 0755, which means read/execute for all, write for owner.
	return os.MkdirAll(dir, 0755)
}
