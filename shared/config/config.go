package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	AppDSN       string
	CloudtimeDSN string
}

var AppConfig *Config

func LoadConfig() {
	// โหลดตัวแปรจากไฟล์ .env และให้ค่าจากไฟล์ override env เดิม (เช่น DB_USER, DB_NAME ที่อาจตั้งไว้ในระบบ)
	_ = godotenv.Overload()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		log.Fatal("DB_HOST not set")
	}

	// 🔥 build postgres dsn
	appDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	AppConfig = &Config{
		AppPort:      getEnv("PORT", "8080"),
		AppDSN:       appDSN,
		CloudtimeDSN: os.Getenv("CLOUDTIME_DSN"),
	}

	if AppConfig.CloudtimeDSN == "" {
		log.Fatal("CLOUDTIME_DSN not set")
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
