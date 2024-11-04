package config

import (
	"github.com/spf13/viper"
	"log"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort int
}

func LoadConfig() *Config {
	// Set up Viper
	viper.SetConfigFile(".env") // or set your custom path
	viper.AutomaticEnv()        // Automatically use environment variables

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	return &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	value := viper.GetString(key)
	if parsedValue, err := strconv.Atoi(value); err == nil {
		return parsedValue
	}
	return defaultValue
}

