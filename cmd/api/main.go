package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/resya202/S16TEST/internal/api"
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

	r := gin.Default()
	api.RegisterRoutes(r, conn)
	r.Run() // listens on :8080
}
