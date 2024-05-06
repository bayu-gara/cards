package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/bayu-gara/cards/pkg/config"

	//external
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteConnector struct{}

func (sc SQLiteConnector) Connect(cfg config.DatabaseConfig) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "cards.db")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConnection)
	db.SetMaxOpenConns(cfg.MaxOpenConnection)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnectionMaxIdleTimeMinutes) * time.Minute)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnectionMaxLifeTimeMinutes) * time.Minute)

	return db, err
}
