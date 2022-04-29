-- name: CreateUser :execresult
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  ?, ? , ?, ?
);

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;
