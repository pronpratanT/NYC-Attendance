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
	SQLExpressDSN       string
	CloudtimeDSN        string
	RedisHost           string
	RedisPort           string
	RedisPassword       string
	RedisDB             string
	JWTSecret           string
	JWTAccessTTLMinutes string
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

	// SQL EXPRESS (SQLEXPRESS_*) DSN (ถ้าใช้ B PLUS หรือ DB อื่นแยกจาก ECONS)
	sqxHost := os.Getenv("SQLEXPRESS_HOST")
	var sqlExpressDSN string
	if sqxHost != "" {
		sqxUser := mustEnv("SQLEXPRESS_USER")
		sqxPass := mustEnv("SQLEXPRESS_PASSWORD")
		// sqxPort := mustEnv("SQLEXPRESS_PORT")
		sqxDB := mustEnv("SQLEXPRESS_DB")
		sqxInstance := os.Getenv("SQLEXPRESS_INSTANCE")

		// ใส่ encrypt=disable เพื่อไม่ให้ Go TLS บังคับใช้ TLS1.2 กับ SQL Express เก่าที่รองรับแค่ TLS เก่า
		// รูปแบบ: sqlserver://user:pass@host:port?database=DB&encrypt=disable[&instance=SQLEXPRESS]
		sqlExpressDSN = fmt.Sprintf("sqlserver://%s:%s@%s/%s?database=%s&encrypt=disable",
			sqxUser,
			sqxPass,
			sqxHost,
			sqxInstance,
			sqxDB,
		)
		if sqxInstance != "" {
			sqlExpressDSN = fmt.Sprintf("%s&instance=%s", sqlExpressDSN, sqxInstance)
		}
	}

	// Cloudtime DSN
	cloudtimeDSN := mustEnv("CLOUDTIME_DSN")

	AppConfig = &Config{
		AppPort:             getEnv("PORT", "8080"),
		AppDSN:              appDSN,
		ECONS_SQLSERVER_DSN: ECONS_SQLSERVER_DSN,
		SQLExpressDSN:       sqlExpressDSN,
		CloudtimeDSN:        cloudtimeDSN,
		RedisHost:           mustEnv("REDIS_HOST"),
		RedisPort:           mustEnv("REDIS_PORT"),
		RedisPassword:       mustEnv("REDIS_PASSWORD"),
		RedisDB:             mustEnv("REDIS_DB"),
		JWTSecret:           mustEnv("JWT_SECRET"),
		JWTAccessTTLMinutes: mustEnv("JWT_ACCESS_TTL_MINUTES"),
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
