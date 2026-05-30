package database

import (
	"database/sql"
	"fmt"
	"go-gaurd/core/config"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.DBUser + ":" + cfg.DBPass + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to MySQL successfully")

	boil.SetDB(db)
	return &Database{DB: db}, nil
}
