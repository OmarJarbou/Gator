package config

import (
	"os"
	"encoding/json"
)

func SetUser(current_user_name string, config *Config) error {
	config.CurrentUserName = current_user_name
	fullURL, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err2 := json.Marshal(config)
	if err2 != nil {
		return err2
	}

	err3 := os.WriteFile(fullURL, jsonData, 0666)
	if err3 != nil {
		return err3
	}
	return nil
} 