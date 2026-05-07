package upload

import (
    "github.com/spf13/viper"
    "log"
)

type UploadConfig struct {
    BasePath  string `mapstructure:"base_path"`
    ReturnUrl string `mapstructure:"return_url"`
}

func LoadUploadConfig() UploadConfig {
    // chỉ định file config trực tiếp
    viper.SetConfigFile("internal/configs/config.yaml")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

    var cfg UploadConfig
    if err := viper.UnmarshalKey("upload", &cfg); err != nil {
        log.Fatalf("Error unmarshalling upload config: %v", err)
    }

    return cfg
}
