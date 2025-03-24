package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"codepair-sinarmas/models"
	"codepair-sinarmas/pkg/serror"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (cfg *Config) InitPostgres() serror.SError {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname =%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatalf("failed connect to database %+v", err)
		return serror.NewFromError(err)
	}

	err = db.Debug().AutoMigrate(
		models.User{},
		models.OTPLog{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database %+v", err)
		return serror.NewFromError(err)
	}

	if db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			users := []models.User{
				{Name: "Admin 1", Email: "admin@gmail.com", Password: "admin123"},
			}
			if err := db.Create(&users).Error; err != nil {
				log.Printf("Error seeding users: %s", err)
			} else {
				log.Println("Users seeded successfully")
			}
		}
	}

	cfg.DB = db

	// GlobalShutdown.RegisterGracefullyShutdown("database/postgres", func(ctx context.Context) error {
	// 	return cfg.DB.Close()
	// })

	return nil
}
