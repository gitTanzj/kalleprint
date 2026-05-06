package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	UploadPath string

	PrusaSlicerPath   string
	PrusaSlicerConfig string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://127.0.0.1"),
		Port:       getEnv("PORT", "8080"),

		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "kalleprint"),

		UploadPath: getEnv("UPLOAD_PATH", "/kalleprint/uploads/quotes"),

		PrusaSlicerPath:   getEnv("PRUSA_SLICER_PATH", "/prusa-slicer/AppRun"),
		PrusaSlicerConfig: getEnv("PRUSA_SLICER_CONFIG", "/kalleprint/config/prusa-slicer.ini"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
