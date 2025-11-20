package database

import (
	"fmt"
	"gym-admin/internal/config"
	"gym-admin/internal/models"
	"log"

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

	// Seed initial data
	if err := seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
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

// seedData creates initial test data if database is empty
func seedData() error {
	// Check if admin user exists
	var count int64
	DB.Model(&models.User{}).Where("phone = ?", "admin").Count(&count)
	if count > 0 {
		return nil // Admin already exists
	}

	// Create admin user
	adminUser := &models.User{
		Name:   "管理员",
		Phone:  "admin",
		Gender: 1,
		Status: 1,
	}

	if err := DB.Create(adminUser).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Println("Created admin user: phone=admin, password=admin123")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
