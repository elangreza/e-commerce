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
	SeederFolder    string
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
		isSeeder := false
		err := config.migrator(db, config.MigrationFolder, isSeeder)
		if err != nil {
			return nil, err
		}
	}

	if config.SeederFolder != "" {
		isSeeder := true
		err := config.migrator(db, config.SeederFolder, isSeeder)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (c *Config) migrator(db *sql.DB, migratorFolder string, isSeeder bool) error {
	var driver database.Driver
	var err error
	switch c.DriverName {
	case "sqlite3":
		conf := &sqlite3.Config{}
		if isSeeder {
			conf.MigrationsTable = "schema_seeders"
		}
		driver, err = sqlite3.WithInstance(db, conf)
	case "postgres":
		conf := &postgres.Config{}
		if isSeeder {
			conf.MigrationsTable = "schema_seeders"
		}
		driver, err = postgres.WithInstance(db, conf)
	}
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migratorFolder,
		c.DriverName,
		driver)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
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

func WithAutoSeeder(seederFolder string) Option {
	return func(c *Config) error {
		c.SeederFolder = seederFolder
		return nil
	}
}
