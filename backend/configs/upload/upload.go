package upload

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type UploadConfig struct {
	BasePath  string `mapstructure:"base_path"`
	ReturnUrl string `mapstructure:"return_url"`
}

func LoadUploadConfig() UploadConfig {
	// Try to find config file in common locations
	configPaths := []string{
		"configs/config.yaml",
		"../../configs/config.yaml",
		"../configs/config.yaml",
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		log.Fatal("config.yaml not found in any expected location")
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg UploadConfig
	if err := viper.UnmarshalKey("upload", &cfg); err != nil {
		log.Fatalf("Error unmarshalling upload config: %v", err)
	}

	return cfg
}
