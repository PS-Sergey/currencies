package postgres

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

const (
	maxOpenConnections = 60
	connMaxLifetime    = 120
	maxIdleConnections = 30
	connMaxIdleTime    = 20
)

func (c *Config) toPgConnection() string {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.DBName,
		c.Password,
		c.SSLMode,
	)

	return dataSourceName
}

func (c *Config) NewClient() (*sqlx.DB, error) {
	connectionString := c.toPgConnection()

	db, err := otelsqlx.Connect(c.PgDriver, connectionString,
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
			attribute.KeyValue{Key: "driver", Value: attribute.StringSliceValue(sql.Drivers())},
		),
		otelsql.WithDBName(c.DBName),
	)

	if err != nil {
		return nil, errors.Wrap(err, "database connection")
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "db ping")
	}

	return db, nil
}
