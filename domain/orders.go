package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID   `json:"id"`
	Amount    float64     `json:"amount"`
	Status    string      `json:"status"`
	Items     []OrderItem `json:"items,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID      uuid.UUID `json:"id"`
	OrderID uuid.UUID `json:"order_id"`
	Name    string    `json:"name"`
	Qty     int       `json:"qty"`
}