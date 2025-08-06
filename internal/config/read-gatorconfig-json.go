package config

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
)

const configFileName = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullURL := homeDir + configFileName
	return fullURL, nil
}

func Read() Config {
	var config Config
	fullURL, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
		return config
	}
	file, err2 := os.Open(fullURL)
	if err2 != nil {
		log.Fatal(err2)
		return config
	}
	defer file.Close()

	data := make([]byte, 1000)
	count, err3 := file.Read(data)
	if err3 != nil {
		log.Fatal(err3)
		return config
	}
	
	fileData := data[:count]
	if err4 := json.Unmarshal(fileData, &config); err4 != nil {
		fmt.Println("failed while decoding data", err4)
		return config
	}
	return config
}
