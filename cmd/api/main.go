package main

import (
	"log"

	"github.com/VicAlexandre/pds-backend/internal/app"
)

func main() {
	cfg := app.NewConfig(":8080")

	app := app.NewApplication(cfg)

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
