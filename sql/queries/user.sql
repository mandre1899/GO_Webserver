-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
    $1, $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetUserPasswordHashByMail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserById :one

SELECT id, created_at, updated_at, email
FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
