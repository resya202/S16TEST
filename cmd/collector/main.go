package main

import (
	"log"
	"time"

	"github.com/resya202/S16TEST/internal/collector"
	"github.com/resya202/S16TEST/internal/config"
	"github.com/resya202/S16TEST/internal/db"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("config load: %v", err)
	}
	conn, err := db.Connect(config.PostgresURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	if err := collector.FetchAndStoreAll(conn); err != nil {
		log.Printf("initial fetch error: %v", err)
	}

	ticker := time.NewTicker(time.Hour)
	for range ticker.C {
		if err := collector.FetchAndStoreAll(conn); err != nil {
			log.Printf("hourly fetch error: %v", err)
		}
	}
}
