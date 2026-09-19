package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL string 			`json:"db_url"`
	CurrentUserName string 	`json:"current_user_name"`
}

const confName = ".gatorconfig.json"

func Read() (*Config, error)  {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, fmt.Errorf("Couldn't find Home directory! D:")
	}

	data, err := os.ReadFile(filepath.Join(homedir, confName))
	if err != nil {
		return &Config{}, fmt.Errorf("Couldn't find Home directory! D:")
	}

	var result = Config{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return &Config{}, fmt.Errorf("Couldn't unmarshal the config data! D:")
	}
	return &result, nil
}

func SetUser(cfg *Config, name string) error {
	cfg.CurrentUserName = name
	yeison, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	err = writeFile(yeison)
	if err != nil {
		return err
	}
	return nil
}

func writeFile(data []byte) error {
	dir, err := getConfDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(dir, data, 0644)
	if err != nil {
		return fmt.Errorf("An error ocurred while writing the config file")
	}

	return nil
}

func getConfDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Couldn't get the config Path!")
	}
	return filepath.Join(homedir, confName), nil
}