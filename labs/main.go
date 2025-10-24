package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:@localhost:5432/labsDB"
	}
	db, err := ConnectDB(dbURL)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	if err := InitDB(db); err != nil {
		log.Fatal("DB init error:", err)
	}

	// Load cache from DB
	LoadCache(db)

	// Connect NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	sc, err := stan.Connect("test-cluster", "order-service-1", stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal("NATS connection error:", err)
	}
	defer sc.Close()

	// Subscribe to orders
	if err := SubscribeAndStore(sc, db); err != nil {
		log.Fatal("NATS subscribe error:", err)
	}

	// Setup HTTP server
	SetupRoutes()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
