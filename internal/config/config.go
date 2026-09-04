package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	fullPath := filepath.Join(homePath, ".gatorconfig.json")
	data, err2 := os.ReadFile(fullPath)
	if err2 != nil {
		return Config{}, err2
	}

	var config Config
	if err3 := json.Unmarshal(data, &config); err3 != nil {
		return Config{}, err3
	}
	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	homepath, err2 := os.UserHomeDir()
	if err2 != nil {
		return err2
	}

	fullPath := filepath.Join(homepath, ".gatorconfig.json")
	if err3 := os.WriteFile(fullPath, bytes, 0644); err3 != nil {
		return err3
	}
	return nil
}
