package config

import (
	"log"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")

	config.AddConfigPath("../")

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %v", err)
	}

	return config
}
