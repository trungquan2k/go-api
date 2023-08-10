package config

import "os"

type DBConfig struct {
	User     string
	Password string
	Driver   string
	Name     string
	Host     string
	Port     string
}

type HTTPConfig struct {
	Host       string
	Port       string
	ExposePort string
}

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
}

func LoadConfig() *Config {
	return &Config{
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		HTTP: HTTPConfig{
			Host:       os.Getenv("APP_HOST"),
			Port:       os.Getenv("APP_PORT"),
			ExposePort: os.Getenv("EXPOSE_PORT"),
		},
	}
}