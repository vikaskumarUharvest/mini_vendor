package service

import (
	"context"
	"fmt"

	"pgxpostgress/domain"
	"pgxpostgress/repository"
)

type OrderService interface {
	Create(ctx context.Context, o *domain.Order) error
	Get(ctx context.Context, id string) (*domain.Order, error)
	List(ctx context.Context) ([]domain.Order, error)
	Update(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repo: r}
}

func (s *orderService) Create(ctx context.Context, o *domain.Order) error {
	if o.Amount <= 0 {
		return fmt.Errorf("amount must be > 0")
	}
	o.Status = "created"
	return s.repo.Create(ctx, o)
}

func (s *orderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *orderService) List(ctx context.Context) ([]domain.Order, error) {
	return s.repo.List(ctx)
}

func (s *orderService) Update(ctx context.Context, id string, status string) error {
	return s.repo.Update(ctx, id, status)
}

func (s *orderService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}