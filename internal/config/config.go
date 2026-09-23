package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"github.com/google/uuid"
)

type Config struct {
	DbURL string 			`json:"db_url"`
	CurrentUserName string 	`json:"current_user_name"`
	CurrentUserId uuid.UUID 	`json:"current_user_id"`
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

func SetUser(cfg *Config, name string, id uuid.UUID) error {
	cfg.CurrentUserName = name
	cfg.CurrentUserId = id
	
	err := writeConf(cfg)
	if err != nil {
		return err
	}
	return nil
}

func writeConf(conf *Config) error {
	dir, err := getConfDir()
	if err != nil {
		return err
	}
	data, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("An error ocurred while writing the config file")
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

func ResetConf(cfg *Config) error {
	result := Config{DbURL: cfg.DbURL, CurrentUserName: "", CurrentUserId: uuid.UUID{}}
	err := writeConf(&result)
	if err != nil {
		return fmt.Errorf("An error ocurred while reseting the config file")
	}
	return nil
}