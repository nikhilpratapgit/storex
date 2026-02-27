package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var (
	Store *sqlx.DB
)

const (
	SSLModeDisable SSLMode = "disable"
)

type SSLMode string

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {
	connectionStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, databaseName, user, password, sslMode)

	DB, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("database ping failed %w", err)
	}
	fmt.Println("Database connected successfully")
	Store = DB
	return migrateUp(DB)
}
func migrateUp(db *sqlx.DB) error {
	fmt.Println("Starting database migrations")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no new migration to apply")
			return nil
		}
		return fmt.Errorf("migration failed:%w", err)
	}
	fmt.Println("migrations applied successfully")
	return nil
}
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := Store.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				fmt.Println("failed to rollback tx : %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			fmt.Println("failed to commit: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
