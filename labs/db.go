package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var cache = make(map[string]Order)

func ConnectDB(dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func InitDB(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(50) PRIMARY KEY,
		track_number VARCHAR(50),
		entry VARCHAR(50),
		locale VARCHAR(10),
		internal_signature TEXT,
		customer_id VARCHAR(50),
		delivery_service VARCHAR(50),
		shardkey VARCHAR(10),
		sm_id INT,
		date_created TIMESTAMP,
		oof_shard VARCHAR(10)
	);

	CREATE TABLE IF NOT EXISTS delivery (
		order_uid VARCHAR(50) PRIMARY KEY REFERENCES orders(order_uid),
		name VARCHAR(100),
		phone VARCHAR(20),
		zip VARCHAR(20),
		city VARCHAR(50),
		address TEXT,
		region VARCHAR(50),
		email VARCHAR(100)
	);

	CREATE TABLE IF NOT EXISTS payment (
		transaction VARCHAR(50) PRIMARY KEY,
		order_uid VARCHAR(50) REFERENCES orders(order_uid),
		request_id VARCHAR(50),
		currency VARCHAR(10),
		provider VARCHAR(50),
		amount NUMERIC,
		payment_dt BIGINT,
		bank VARCHAR(50),
		delivery_cost NUMERIC,
		goods_total NUMERIC,
		custom_fee NUMERIC
	);

	CREATE TABLE IF NOT EXISTS items (
		rid VARCHAR(50) PRIMARY KEY,
		order_uid VARCHAR(50) REFERENCES orders(order_uid),
		chrt_id INT,
		track_number VARCHAR(50),
		price NUMERIC,
		name VARCHAR(100),
		sale INT,
		size VARCHAR(10),
		total_price NUMERIC,
		nm_id INT,
		brand VARCHAR(50),
		status INT
	);
	`)
	return err
}

func LoadCache(db *pgxpool.Pool) {
	rows, err := db.Query(context.Background(), "SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders")
	if err != nil {
		log.Println("Error loading cache:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerID,
			&o.DeliveryService,
			&o.ShardKey,
			&o.SmID,
			&o.DateCreated,
			&o.OofShard,
		)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		cache[o.OrderUID] = o
	}
	fmt.Println("Cache loaded:", len(cache), "orders")
}
