package contextstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Save(all map[string]FileDetails) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	fmt.Println(configPath)

	err = ensureDir(filepath.Dir(configPath))
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func Load() (map[string]FileDetails, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return make(map[string]FileDetails), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var all map[string]FileDetails
	err = json.Unmarshal(data, &all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AutoCommenter", "contextstore.json"), nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
