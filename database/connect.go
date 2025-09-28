package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseCredentials struct {
	Username    string
	Password    string
	Database    string
	Host        string
	Port        string
	SSLDisabled string
	TimeZone    string
}

func Open(dc *DatabaseCredentials) (*gorm.DB, error) {
	if dc == nil {
		return nil, fmt.Errorf("DatabaseCredentials cannot be nil")
	}
	if dc.Username == "" || dc.Password == "" || dc.Database == "" {
		return nil, fmt.Errorf("DatabaseCredentials fields cannot be empty")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dc.Host, dc.Username, dc.Password, dc.Database, dc.Port, dc.SSLDisabled, dc.TimeZone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
