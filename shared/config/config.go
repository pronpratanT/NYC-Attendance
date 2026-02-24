package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppDSN              string
	ECONS_SQLSERVER_DSN string
	CloudtimeDSN        string
}

var AppConfig *Config

func LoadConfig() {
	// โหลดตัวแปรจากไฟล์ .env และให้ค่าจากไฟล์ override env เดิม (เช่น DB_USER, DB_NAME ที่อาจตั้งไว้ในระบบ)
	_ = godotenv.Overload()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// ดึงและ validate env ที่จำเป็นสำหรับ Postgres
	host := mustEnv("DB_HOST")
	port := mustEnv("DB_PORT")
	user := mustEnv("DB_USER")
	password := mustEnv("DB_PASSWORD")
	dbname := mustEnv("DB_NAME")

	// 🔥 build postgres dsn
	appDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// ECONS SQL Server DSN
	sqlUser := mustEnv("SQLSERVER_USER")
	sqlPass := mustEnv("SQLSERVER_PASSWORD")
	sqlHost := mustEnv("SQLSERVER_HOST")
	sqlPort := mustEnv("SQLSERVER_PORT")
	sqlDB := mustEnv("SQLSERVER_DB")

	ECONS_SQLSERVER_DSN := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		sqlUser,
		sqlPass,
		sqlHost,
		sqlPort,
		sqlDB,
	)

	// Cloudtime DSN
	cloudtimeDSN := mustEnv("CLOUDTIME_DSN")

	AppConfig = &Config{
		AppPort:             getEnv("PORT", "8080"),
		AppDSN:              appDSN,
		ECONS_SQLSERVER_DSN: ECONS_SQLSERVER_DSN,
		CloudtimeDSN:        cloudtimeDSN,
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

// mustEnv คืนค่า env ถ้ามีค่า และถ้าไม่มีจะ log.Fatal เพื่อหยุดโปรแกรมทันที
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s not set", key)
	}
	return val
}
