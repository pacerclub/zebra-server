package database

import (
	"fmt"
	"log"

	"github.com/ZigaoWang/zebra-server/internal/config"
	"github.com/ZigaoWang/zebra-server/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.WorkLog{},
		&domain.LogEntry{},
		&domain.Project{},
		&domain.Session{},
		&domain.Record{},
		&domain.File{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database connected and migrated successfully")
	return db, nil
}
