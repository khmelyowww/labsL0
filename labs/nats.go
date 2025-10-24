package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SubscribeAndStore(sc stan.Conn, db *pgxpool.Pool) error {
	_, err := sc.Subscribe("orders", func(m *stan.Msg) {
		var order Order
		if err := json.Unmarshal(m.Data, &order); err != nil {
			fmt.Println("Invalid message:", err)
			return
		}

		tx, err := db.Begin(context.Background())
		if err != nil {
			fmt.Println("DB transaction error:", err)
			return
		}

		// Orders
		_, err = tx.Exec(context.Background(),
			`INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT(order_uid) DO NOTHING`,
			order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
		if err != nil {
			fmt.Println("DB insert order error:", err)
			tx.Rollback(context.Background())
			return
		}

		// Delivery
		_, err = tx.Exec(context.Background(),
			`INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT(order_uid) DO NOTHING`,
			order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
		if err != nil {
			fmt.Println("DB insert delivery error:", err)
			tx.Rollback(context.Background())
			return
		}

		// Payment
		_, err = tx.Exec(context.Background(),
			`INSERT INTO payment(transaction, order_uid, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT(transaction) DO NOTHING`,
			order.Payment.Transaction, order.OrderUID, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
		if err != nil {
			fmt.Println("DB insert payment error:", err)
			tx.Rollback(context.Background())
			return
		}

		// Items
		for _, item := range order.Items {
			_, err = tx.Exec(context.Background(),
				`INSERT INTO items(rid, order_uid, chrt_id, track_number, price, name, sale, size, total_price, nm_id, brand, status)
				VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) ON CONFLICT(rid) DO NOTHING`,
				item.Rid, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
			if err != nil {
				fmt.Println("DB insert item error:", err)
				tx.Rollback(context.Background())
				return
			}
		}

		tx.Commit(context.Background())

		// Cache
		cache[order.OrderUID] = order
		fmt.Println("Order stored:", order.OrderUID)

	}, stan.DurableName("orders-durable"))
	return err
}
