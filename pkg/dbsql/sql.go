package dbsql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

type Option func(*Config) error

type Config struct {
	DriverName      string
	DataSourceName  string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationFolder string
}

// new sqlite database
func NewDbSql(options ...Option) (*sql.DB, error) {
	// if opt is nil default strategy is using WithSqliteDB
	if len(options) == 0 {
		options = append(options, WithSqliteDB("default.db"))
	}

	var config Config
	for _, opt := range options {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if config.MigrationFolder != "" {
		var driver database.Driver
		switch config.DriverName {
		case "sqlite3":
			driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
		case "postgres":
			driver, err = postgres.WithInstance(db, &postgres.Config{})
		}
		if err != nil {
			return nil, err
		}

		m, err := migrate.NewWithDatabaseInstance(
			config.MigrationFolder,
			config.DriverName,
			driver)
		if err != nil && err != migrate.ErrNoChange {
			return nil, err
		}

		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return nil, err
		}
	}

	return db, nil
}

func WithTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("error occurred while rolling back transaction: %w", rbErr)
		}
		return err
	}
	return tx.Commit()
}

func WithDBConnectionPool(maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) Option {
	return func(c *Config) error {
		c.MaxOpenConns = maxOpenConns
		c.MaxIdleConns = maxIdleConns
		c.ConnMaxLifetime = connMaxLifetime
		return nil
	}
}

func WithAutoMigrate(migrationFolder string) Option {
	return func(c *Config) error {
		c.MigrationFolder = migrationFolder
		return nil
	}
}
