package database

import (
	"database/sql"

	"github.com/bayu-gara/cards/pkg/config"
	sqlite "github.com/bayu-gara/cards/pkg/database/sqlite"
)

type SQLConnector interface {
	Connect(cfg config.DatabaseConfig) (*sql.DB, error)
}

var GetConnector = getConnectorFunc

func getConnectorFunc(engine string) SQLConnector {
	switch engine {
	case "sqlite":
		return sqlite.SQLiteConnector{}
	}

	return nil
}
