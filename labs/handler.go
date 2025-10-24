package main

import (
	"encoding/json"
	"net/http"
)

// SetupRoutes настраивает маршруты HTTP сервера
func SetupRoutes() {
	// API для получения заказа по ID
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Query().Get("id")

		// Всегда возвращаем JSON
		w.Header().Set("Content-Type", "application/json")

		if orderID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "id parameter required"})
			return
		}

		order, ok := cache[orderID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "order not found"})
			return
		}

		// Возвращаем заказ
		json.NewEncoder(w).Encode(order)
	})
}
