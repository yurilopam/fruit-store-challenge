package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBOption func(*gorm.Config)

func NewPostgres(dsn string, opts ...DBOption) (*gorm.DB, error) {
	cfg := &gorm.Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return gorm.Open(postgres.Open(dsn), cfg)
}

func DSN(host, port, user, pass, name string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
}
