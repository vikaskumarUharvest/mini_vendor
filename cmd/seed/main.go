// Run anytime to reset + reseed the entire database.
// Usage: cd vendor-service && go run cmd/seed/main.go
package main

import (
	"context"
	"log"
	"os"

	"pgxpostgress/seed"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:573636@localhost:5432/go_postgres"
	}

	/* the context package is used to control timeouts, cancellation, and request-scoped values.

	   context.Background() is the most basic, empty context:

	   never cancels
	   has no timeout
	   carries no values
	*/
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ DB connect error: %v", err)
	}
	defer pool.Close()
	log.Println("🌱 Starting database seed...")
	log.Println("   50 users | 50 orders | 100–150 order items")
	log.Println("   Random statuses (created/pending/completed/cancelled)")
	log.Println("")

	if err := seed.Run(pool); err != nil {
		log.Fatalf("❌ Seed failed: %v", err)
	}

	log.Println("\n✅ Seed complete! Test data summary:")

	log.Println("   👤 Users:")
	log.Println("      user1@example.com   / password123")
	log.Println("      user2@example.com   / password123")
	log.Println("      ... up to user50@example.com")

	log.Println("")
	log.Println("   📦 Orders:")
	log.Println("      50 orders created with random amounts")
	log.Println("      statuses → created | pending | completed | cancelled")

	log.Println("")
	log.Println("   🧾 Order Items:")
	log.Println("      Each order has 1–3 items")
	log.Println("      Products → Laptop Bag, Mouse, Keyboard, USB Cable, etc.")

	log.Println("")
	log.Println("🚀 Ready for API testing!")
}
