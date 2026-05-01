-- name: CreateUser :exec
INSERT INTO users (
  name, email, password, phone, age, city
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, email = $3, phone = $4, age = $5, city = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
