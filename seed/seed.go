package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var statuses = []string{"created", "pending", "completed", "cancelled"}
var products = []string{"Laptop Bag", "Mouse", "Keyboard", "USB Cable", "Monitor", "Headphones", "Desk Mat", "Webcam"}

func Run(pool *pgxpool.Pool) error {
	ctx := context.Background()

	log.Println("Truncating existing data...")
	_, err := pool.Exec(ctx, "TRUNCATE TABLE order_items, orders, users CASCADE")
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	log.Println("Inserting 50 users...")
	for i := 1; i <= 50; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		name := fmt.Sprintf("Test User %d", i)
		phone := fmt.Sprintf("98765%05d", i)
		age := 20 + rand.Intn(40) // Age between 20 and 59
		city := "Delhi"
		if i%2 == 0 {
			city = "Mumbai"
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO users (name, email, password, phone, age, city)
			VALUES ($1, $2, 'password123', $3, $4, $5)
		`, name, email, phone, age, city)
		if err != nil {
			return fmt.Errorf("failed to insert user %d: %w", i, err)
		}
	}

	log.Println("Inserting 50 orders and order items...")
	for i := 1; i <= 50; i++ {
		amount := float64(100 + rand.Intn(900)) + rand.Float64()
		status := statuses[rand.Intn(len(statuses))]

		var orderID string
		err := pool.QueryRow(ctx, `
			INSERT INTO orders (amount, status)
			VALUES ($1, $2)
			RETURNING id
		`, amount, status).Scan(&orderID)
		if err != nil {
			return fmt.Errorf("failed to insert order %d: %w", i, err)
		}

		// Insert 1-3 items for this order
		numItems := 1 + rand.Intn(3)
		for j := 0; j < numItems; j++ {
			productName := products[rand.Intn(len(products))]
			qty := 1 + rand.Intn(5)

			_, err := pool.Exec(ctx, `
				INSERT INTO order_items (order_id, name, qty)
				VALUES ($1, $2, $3)
			`, orderID, productName, qty)
			if err != nil {
				return fmt.Errorf("failed to insert order item for order %s: %w", orderID, err)
			}
		}
	}

	return nil
}
