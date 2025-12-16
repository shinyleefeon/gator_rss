package config

import (
"os"
"fmt"
"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	config_link, err := os.UserHomeDir() 
	config_link += "/" + configFileName
	if err != nil {
		return "", err
	}
	return config_link, nil
}

func Read() (*Config, error) {
	config_link, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	
	jsonData, err := os.ReadFile(config_link)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func Write(cfg Config) error {
	
	config_link, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing config data: %v", err)
	}

	err = os.WriteFile(config_link, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

func (c Config) SetUser(userName string) {
	c.Current_user_name = userName
	Write(c)
}

	
