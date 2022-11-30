package database

import (
	"fmt"
	"log"

	"assignment-golang-backend/internal/config"
	"assignment-golang-backend/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type DBConfig struct {
	HOST string
	PORT string
	USER string
	PASS string
	NAME string
}

func getConfig() *DBConfig {
	return &DBConfig{
		HOST: config.GetEnv("DB_HOST"),
		PORT: config.GetEnv("DB_PORT"),
		USER: config.GetEnv("DB_USER"),
		PASS: config.GetEnv("DB_PASS"),
		NAME: config.GetEnv("DB_NAME"),
	}
}

func Connect() {
	dbConfig := getConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.HOST,
		dbConfig.PORT,
		dbConfig.USER,
		dbConfig.PASS,
		dbConfig.NAME,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Wallet{},
		&entity.Transaction{},
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func Get() *gorm.DB {
	return db
}
