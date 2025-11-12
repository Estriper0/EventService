package main

import (
	"fmt"

	"github.com/Estriper0/EventService/internal/config"
	"github.com/Estriper0/EventService/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := config.New()

	db := database.GetDB(&config.DB)
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()
	fmt.Println("Migrations complete!")
}
