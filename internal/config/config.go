package config

import "os"
type Config struct {}

func Read() (*Config, error) {
	config_link := os.UserHomeDir() + ".gatorconfig.json"
	fmt.Println("Reading config from ", config_link)
	return &Config{}, nil
}