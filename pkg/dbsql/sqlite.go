package dbsql

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func WithSqliteDB(fileName string) Option {
	return func(c *Config) error {
		// Ensure the directory exists with 0777 permissions
		// This allows both Docker containers and local development to write
		dir := filepath.Dir(fileName)
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("failed to create db dir: %w", err)
		}
		c.DriverName = "sqlite3"
		c.DataSourceName = fileName
		return nil
	}
}

func WithSqliteDBWalMode() Option {
	return func(c *Config) error {
		if c.DriverName != "sqlite3" {
			return fmt.Errorf("driver is not sqlite3")
		}
		c.DataSourceName = fmt.Sprintf("%s?_journal_mode=WAL", c.DataSourceName)
		return nil
	}
}
