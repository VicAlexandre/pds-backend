package main

import (
	"log"

	"github.com/VicAlexandre/pds-backend/internal/app"
)

func main() {
	app := app.NewApplication(app.NewConfig(":8080"))

	mux := app.Mount()

	log.Fatal(app.Run(mux))
}
