package repository

import (
	"context"

	"pgxpostgress/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	List(ctx context.Context) ([]domain.User, int, error)
	Create(ctx context.Context, user domain.User) error
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
