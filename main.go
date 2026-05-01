package main

import (
	"context"
	"fmt"
	"log"
	"pgxpostgress/route"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	fmt.Println("---------pgx Poool Postgeess--------")
	dbURL := "postgres://postgres:573636@localhost:5432/librarydb"

	pool, err := pgxpool.New(
		context.Background(),
		dbURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	r := gin.Default()

	route.SetupRoutes(r, pool)

	r.Run(":8080")

}
