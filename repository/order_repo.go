package repository

import (
	"context"
	"fmt"

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
	// 1. `tx, err := r.db.Begin(ctx)` -> This tells the database, *"Hey, start a new transaction. Don't permanently save anything until I say so!"*

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// 2. `defer tx.Rollback(ctx)` -> This is a safety net. It says, *"If this function crashes or returns an error before committing, undo everything we just did."*

	o.ID = uuid.New()
	// 3. `tx.QueryRow(...)` -> This successfully inserts the `order`. (At this point, the order is in a temporary state in the database).

	err = tx.QueryRow(ctx,
		`INSERT INTO orders (id, amount, status, created_at)
		 VALUES ($1,$2,$3,NOW()) RETURNING created_at`,
		o.ID, o.Amount, o.Status,
	).Scan(&o.CreatedAt)

	if err != nil {
		return err
	}
	// Insert items
	for i := range o.Items {
		// 1. MODIFIED: Only generate a new UUID if Postman didn't provide one.
		// This allows the "Duplicate ID" test to actually reach the database.
		if o.Items[i].ID == uuid.Nil {
			o.Items[i].ID = uuid.New()
		}

		o.Items[i].OrderID = o.ID

		// 2. Validation: Catch empty names before hitting the DB
		if o.Items[i].Name == "" {
			return fmt.Errorf("item name cannot be empty")
		}

		// 3. Test Trigger: Manual fail for testing rollback
		if o.Items[i].Name == "FAIL" {
			return fmt.Errorf("simulated database crash for testing")
		}

		// 4. DB Execution
		_, err := tx.Exec(ctx,
			`INSERT INTO order_items (id, order_id, name, qty)
         VALUES ($1,$2,$3,$4)`,
			o.Items[i].ID, o.Items[i].OrderID, o.Items[i].Name, o.Items[i].Qty,
		)

		if err != nil {
			// If SQL throws a "Duplicate Key" or "Constraint" error,
			// we return here, and 'defer tx.Rollback(ctx)' cleans up the Order.
			return fmt.Errorf("database error: %w", err)
		}
	}
	// 5. Because the function returned early, it never reaches `tx.Commit(ctx)`, and the `defer tx.Rollback(ctx)` triggers! The database completely erases the temporary `order` it just inserted.

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
