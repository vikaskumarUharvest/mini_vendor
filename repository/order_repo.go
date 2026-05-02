package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"pgxpostgress/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	Get(ctx context.Context, id string) (*domain.Order, error)
	List(ctx context.Context) ([]domain.Order, error)
	Update(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, o *domain.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	o.ID = uuid.New()

	err = tx.QueryRow(ctx,
		`INSERT INTO orders (id, amount, status, created_at)
		 VALUES ($1,$2,$3,NOW()) RETURNING created_at`,
		o.ID, o.Amount, o.Status,
	).Scan(&o.CreatedAt)

	if err != nil {
		return err
	}

	// Insert items
	for _, item := range o.Items {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_items (id, order_id, name, qty)
			 VALUES ($1,$2,$3,$4)`,
			uuid.New(), o.ID, item.Name, item.Qty,
		)
		if err != nil {
			return err // rollback triggered
		}
	}

	return tx.Commit(ctx)
}

func (r *orderRepo) Get(ctx context.Context, id string) (*domain.Order, error) {
	var o domain.Order

	err := r.db.QueryRow(ctx,
		`SELECT id, amount, status, created_at 
		 FROM orders WHERE id=$1`, id,
	).Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, name, qty FROM order_items WHERE order_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		rows.Scan(&item.ID, &item.Name, &item.Qty)
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *orderRepo) List(ctx context.Context) ([]domain.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, amount, status, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order

	for rows.Next() {
		var o domain.Order
		rows.Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *orderRepo) Update(ctx context.Context, id string, status string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE orders SET status=$1 WHERE id=$2`,
		status, id,
	)
	return err
}

func (r *orderRepo) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM order_items WHERE order_id=$1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}