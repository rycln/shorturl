package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rycln/shorturl/internal/db"
)

func main() {
	uri := flag.String("d", "", "Database connection address")
	flag.Parse()

	if *uri == "" {
		log.Fatal("dsn required")
	}

	database, err := sql.Open("pgx", *uri)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	goose.SetBaseFS(db.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(database, "migrations"); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		log.Fatal(err)
	}

	log.Print("Migrations applied successfully")
}
