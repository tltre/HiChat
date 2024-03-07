package main

import (
	"HiChat/global"
	"HiChat/initialize"
	"HiChat/models"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	initialize.InitConfig("debug")
	sqlConfig := global.ServiceConfig.DB
	// declare the connection to the DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		sqlConfig.User, sqlConfig.Password, sqlConfig.Host, sqlConfig.Port, sqlConfig.Name)
	// connect to the Database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func CreateTables(db *gorm.DB) {
	createUserTable(db)
	createRelationTable(db)
	createCommunityTable(db)
}

func createUserTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.UserBasic{})
	if err != nil {
		panic(err)
	}
}

func createRelationTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.Relation{})
	if err != nil {
		panic(err)
	}
}

func createCommunityTable(db *gorm.DB) {
	err := db.AutoMigrate(&models.Community{})
	if err != nil {
		panic(err)
	}
}

func ConnectToRedis() *redis.Client {
	redisConfig := global.ServiceConfig.RedisDB
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	}
	return redis.NewClient(&opt)
}

func CheckIfConnectRedis(client *redis.Client) {
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to Connect to Redis")
		return
	}
	fmt.Println("Success to Connect to Redis")
}
