package config

import (
	"log"

	"github.com/spf13/viper"
)


func viperConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		panic(err)
	}
}
