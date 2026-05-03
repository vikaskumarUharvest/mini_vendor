package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Phone     string    `json:"phone" db:"phone"`
	Age       int       `json:"age" db:"age"`
	Country   string    `json:"country"` // [NEW] added
	City      string    `json:"city" db:"city"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type PageResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}
