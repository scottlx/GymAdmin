package database

import (
	"fmt"
	"gym-admin/internal/config"
	"gym-admin/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DB = db

	// Auto migrate tables
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.UserTrainingStats{},
		&models.CardType{},
		&models.MembershipCard{},
		&models.Coach{},
		&models.Course{},
		&models.Booking{},
		&models.CheckIn{},
		&models.FaceRecord{},
		&models.VoucherRecord{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
