package postgres

import (
	"context"

	"pgxpostgress/domain"
	"pgxpostgress/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(
	pool *pgxpool.Pool,
) repository.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) List(
	ctx context.Context,
) ([]domain.User, int, error) {

	query := `
		SELECT id, email, name, city
		FROM users
	`
  // pool.Query for multiple Query
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []domain.User

	// Loop Through Query Results
	for rows.Next() {
		// Create One User Per Row
		var user domain.User

		// Scan Database Columns Into Struct
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.City,
		)
		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	totalUsers := len(users)

	return users, totalUsers, nil
}

func (r *userRepository) Create(
	ctx context.Context,
	user domain.User,
) error {

	query := `
		INSERT INTO users
		(id,name,email,password,phone,age,city)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
// Exec: For commands that don't return rows.
	_, err := r.pool.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.Age,
		user.City,
	)

	return err
}

func (r *userRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {

	query := `
		SELECT id,email,name,city
		FROM users
		WHERE id=$1
	`

	var user domain.User
  // QueryRow: For exactly one row.
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.City,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	user domain.User,
) error {

	query := `
		UPDATE users
		SET name=$2, email=$3, phone=$4, age=$5, city=$6
		WHERE id=$1
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.Age,
		user.City,
	)

	return err
}

func (r *userRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	query := `
		DELETE FROM users
		WHERE id=$1
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		id,
	)

	return err
}
