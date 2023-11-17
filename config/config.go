package config

import (
	"log"

	"github.com/spf13/viper"
)

// Struct to map env values.
type Config struct {
	DBUser string `mapstructure:"POSTGRES_USER"`
	DBPass string `mapstructure:"POSTGRES_PASSWORD"`
	DBHost string `mapstructure:"POSTGRES_HOST"`
	DBPort int    `mapstructure:"POSTGRES_PORT"`
	DBName string `mapstructure:"POSTGRES_DB"`
}

// Call to get a new instance of config with .env variables.
func LoadConfig(path string) (config Config, err error) {
	// Set path/location of env file.
	viper.AddConfigPath(path)

	// Tell viper the name of config file.
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Viper reads all the variables.
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Unmarshal into our struct.
	err = viper.Unmarshal(&config)

	return
}
