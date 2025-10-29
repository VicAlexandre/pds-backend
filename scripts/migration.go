package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	// Ensure migrations table exists
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	migrations := []struct {
		version string
		query   string
	}{
		{
			version: "000_uuid_extension",
			query: `
				CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
				`,
		},
		{
			version: "001_create_users",
			query: `
				CREATE TABLE IF NOT EXISTS users (
					id SERIAL PRIMARY KEY,
					name TEXT NOT NULL,
					email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				)
					
			`,
		},
		{
			version: "002_create_apostilas",
			query: `
				CREATE TABLE IF NOT EXISTS apostilas (
					id UUID PRIMARY KEY,
					user_id INTEGER NOT NULL REFERENCES users(id),
					edited_html TEXT,
					pdf_raw BYTEA,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			)
			`,
		},
	}

	for _, m := range migrations {
		var exists bool
		err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", m.version).Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		if exists {
			fmt.Printf("Migration %s already applied\n", m.version)
			continue
		}

		fmt.Printf("Applying migration %s...\n", m.version)
		_, err = db.ExecContext(ctx, m.query)
		if err != nil {
			log.Fatalf("Failed migration %s: %v\n", m.version, err)
		}

		_, err = db.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", m.version)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Migration %s applied\n", m.version)
	}
}
