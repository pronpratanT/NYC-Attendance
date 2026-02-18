package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pronpratanT/leave-system/shared/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	CloudtimeDB *gorm.DB
	AppDB       *gorm.DB
)

func ConnectDB() *gorm.DB {
	if AppDB != nil {
		return AppDB
	}

	config.LoadConfig()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")

	log.Printf("DB env: host=%s port=%s user=%s password=%s db=%s",
		host, port, user, password, dbname)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		dbname,
		port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect postgres: ", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	AppDB = db
	return AppDB
}

func ConnectCloudtime() *gorm.DB {
	if CloudtimeDB != nil {
		return CloudtimeDB
	}

	config.LoadConfig()

	db, err := gorm.Open(mysql.Open(config.AppConfig.CloudtimeDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Error connecting Cloudtime: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("Error getting sqlDB: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	CloudtimeDB = db
	log.Println("Cloudtime Connected Successfully")
	return CloudtimeDB
}
