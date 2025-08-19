package main

import (
	"log"
	"time"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
