package dbsql

import (
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// func WithSqliteDB(fileName string) Option {
// 	if !fileExists(fileName) {
// 		if err := createDBFile(fileName); err != nil {
// 			return func(c *Config) error {
// 				return err
// 			}
// 		}
// 	}

// 	return func(c *Config) error {
// 		c.DriverName = "sqlite3"
// 		c.DataSourceName = fileName
// 		return nil
// 	}
// }

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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Directory does not exist, create it with permissions 0755
		err = os.MkdirAll(dirName, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDBFile(filename string) error {

	err := ensureDir(filename)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
