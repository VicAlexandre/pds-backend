package main

import (
	"log"

	"github.com/VicAlexandre/pds-backend/internal/app"
	"github.com/VicAlexandre/pds-backend/internal/db"
)

func main() {
	app := app.NewApplication(app.NewConfig(":8080"))

	// TODO: load db config from .env
	cfg := db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		DBName:   "pdsdb",
		SSLMode:  "disable",
	}

	conn, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	mux := app.Mount(conn)

	log.Fatal(app.Run(mux))
}
