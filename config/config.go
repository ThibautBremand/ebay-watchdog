package config

import (
	"fmt"
	"github.com/spf13/viper"
	"text/template"
)

type SearchItem struct {
	URL     string
	Domains []string
}

// Load loads both the toml config and the .env file.
func Load() error {
	err := loadConfig()
	if err != nil {
		return err
	}

	err = loadEnv(err)
	if err != nil {
		return err
	}

	return nil
}

// LoadSearchItems loads a list of SearchItem items from the config file.
func LoadSearchItems() ([]SearchItem, error) {
	var searchItems []SearchItem

	err := viper.UnmarshalKey("searches", &searchItems)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal searchItems key of config: %v", err)
	}

	return searchItems, nil
}

// LoadTemplate loads the message template from the config file, used for the Telegram messages format.
func LoadTemplate() (*template.Template, error) {
	return template.New("message").Parse(viper.GetString("message"))
}

// loadConfig loads the toml file.
func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not load config.toml: %v", err)
	}

	return nil
}

// loadEnv loads the env file.
func loadEnv(err error) error {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	err = viper.MergeInConfig()
	if err != nil {
		return fmt.Errorf("could not load .env: %v", err)
	}

	return nil
}
