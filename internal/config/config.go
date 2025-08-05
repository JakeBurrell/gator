package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const filename = ".gatorconfig.json"

type Config struct {
	DataBaseURL     string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	fileLocation, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonData, err := os.ReadFile(fileLocation)

	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = json.Unmarshal(jsonData, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username

	return write(*cfg)
}

func getConfigFilePath() (string, error) {

	dirName, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(dirName, filename)
	return fullPath, nil

}

func write(cfg Config) error {
	fileLocation, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil

}
