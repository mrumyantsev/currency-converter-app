package database

import (
	"fmt"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"

	"database/sql"

	_ "github.com/lib/pq" // necessary for Postgres driver
)

// A Database is used to control the connection to a database.
type Database struct {
	config *config.Config
	*sql.DB
}

func New(cfg *config.Config) *Database {
	return &Database{
		config: cfg,
	}
}

// Connect connects to the database.
func (d *Database) Connect() error {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		d.config.DbHostname,
		d.config.DbPort,
		d.config.DbUsername,
		d.config.DbPassword,
		d.config.DbDatabase,
		d.config.DbSSLMode,
	)

	db, err := sql.Open(d.config.DbDriver, dataSourceName)
	if err != nil {
		return errlib.Wrap("could not connect to db", err)
	}

	d.DB = db

	return nil
}

// Disconnect disconnects from the database.
func (d *Database) Disconnect() error {
	if err := d.DB.Close(); err != nil {
		return errlib.Wrap("could not disconnect from db", err)
	}

	return nil
}
