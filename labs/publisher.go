package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	// Читаем model.json
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		log.Fatal("Error reading model.json:", err)
	}

	// Подключаемся к NATS-Streaming
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	sc, err := stan.Connect("test-cluster", "publisher-1", stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal("NATS connection error:", err)
	}
	defer sc.Close()

	// Проверим, что JSON валиден
	var test interface{}
	if err := json.Unmarshal(data, &test); err != nil {
		log.Fatal("Invalid JSON in model.json:", err)
	}

	// Публикуем данные
	if err := sc.Publish("orders", data); err != nil {
		log.Fatal("Publish error:", err)
	}

	fmt.Println("Order published successfully!")
}
