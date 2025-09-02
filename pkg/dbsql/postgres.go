package dbsql

import (
	"fmt"

	_ "github.com/lib/pq"
)

func WithPostgresDB(host, port, user, password, dbname, sslmode string) Option {
	return func(c *Config) error {
		c.DriverName = "postgres"
		c.DataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			host, port, user, password, dbname)
		return nil
	}
}

func WithPostgresDBSSLMode(sslmode string) Option {
	return func(c *Config) error {
		if c.DriverName != "postgres" {
			return fmt.Errorf("driver is not postgres")
		}
		c.DataSourceName = fmt.Sprintf("%s sslmode=%s", c.DataSourceName, sslmode)
		return nil
	}
}
