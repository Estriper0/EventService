package main

import (
	"database/sql"
	"fmt"

	"github.com/Estriper0/EventService/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	config := config.New()

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.DB.DbUser,
			config.DB.DbPassword,
			config.DB.DbHost,
			config.DB.DbPort,
			config.DB.DbName,
			config.DB.SSLMode,
		),
	)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS event")
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "event.migrations",
		SchemaName:      "event",
	})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations complete!")
}
