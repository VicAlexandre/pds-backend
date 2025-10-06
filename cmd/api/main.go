package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/VicAlexandre/pds-backend/internal/app"
	_ "github.com/lib/pq"
	// "github.com/VicAlexandre/pds-backend/internal/db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // rodar localmente
	}
	addr := ":" + port

	app := app.NewApplication(app.NewConfig(addr))

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// então, rodar localmente
		dsn = "postgres://pds:secret@localhost:5432/pds?sslmode=disable"
		log.Println("DATABASE_URL not set, using local fallback")
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	mux := app.Mount(dbConn)

	log.Println("Connected to database successfully!")
	log.Fatal(app.Run(mux))
}
