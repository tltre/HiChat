package initialize

import (
	"HiChat/src/config"
	"HiChat/src/global"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// InitDB  initial the connection to the MySQL DB
func InitDB() {

	// declare the connection to the DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	// set the config of logger
	newLogger := logger.New(
		// writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		// config
		logger.Config{
			// set the threshold of slow query, and the query cost more than 1s will be record
			SlowThreshold: time.Second,
			// print different kind of information by different color
			Colorful: true,
			// ignore the error that didn't exist the record
			IgnoreRecordNotFoundError: true,
			// we record all kind of log, can set as Silent->Error->Warn->Info
			LogLevel: logger.Info,
		},
	)

	// connect to DB
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}

// InitRedis initial the connection to the Redis DB
func InitRedis() {
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: "",
		DB:       0,
	}
	global.RedisDB = redis.NewClient(&opt)
}
