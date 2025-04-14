package config

import "os"

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type JWTConfig struct {
	Secret string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func LoadJWTConfig() JWTConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	return JWTConfig{Secret: secret}
}
