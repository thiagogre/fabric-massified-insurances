package db

import "errors"

type DatabaseType string

const (
	SQLite DatabaseType = "sqlite"
)

func NewDatabase(dbType DatabaseType, dataSourceName string) (Database, error) {
	switch dbType {
	case SQLite:
		return NewSQLiteDB(dataSourceName)
	default:
		return nil, errors.New("unsupported database type")
	}
}
