package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb() *gorm.DB {
	dsn := "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	// GORMを使用してMySQLに接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// 接続確認
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database connection")
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	err = sqlDB.Ping()
	if err != nil {
		panic("failed to ping database")
	}

	return db
}
