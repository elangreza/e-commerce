package dbsql

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func WithSqliteDB(fileName string) Option {
	if !fileExists(fileName) {
		if err := createDBFile(fileName); err != nil {
			return func(c *Config) error {
				return err
			}
		}
	}

	return func(c *Config) error {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func createDBFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
