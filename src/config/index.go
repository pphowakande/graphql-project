package configuration

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// GetConfig ... This function gets configuration from .yaml files
func GetConfig() Configuration {
	env := os.Getenv("TIER")
	log.Println("Env : ", env)
	if env == "" {
		env = "local"
	}
	if env != "dev" && env != "ote" && env != "production" && env != "sit" && env != "local" {
		log.Fatal("invalid env for loading the configuration set TIER value as local, dev or production")
	}
	viper.SetConfigName(env)
	viper.AddConfigPath("./src/config/tier")
	var configuration Configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configuration
}
