package main

import (
	"HiChat/src/config"
	"HiChat/src/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	// declare the connection to the DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	// connect to the Database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func CreateUserTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.UserBasic{})
	if err != nil {
		panic(err)
	}
}

func CreateRelationTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.Relation{})
	if err != nil {
		panic(err)
	}
}
