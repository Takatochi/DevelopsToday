package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string, opts ...GormOption) (*gorm.DB, error) {
	cfg := &gorm.Config{}

	for _, opt := range opts {
		opt(cfg)
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
