-- name: CreateUser :one
INSERT INTO users (username, password, full_name, email) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateUser :one
UPDATE users
SET
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    password = COALESCE(sqlc.narg(password), password),
    password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at)
WHERE username = sqlc.arg(username)
RETURNING *;


-- name: UpdateUserCaseExample :one
UPDATE users 
SET
    password = CASE 
        WHEN @set_password::boolean = TRUE THEN @password
        ELSE password 
    END,
    full_name = CASE 
        WHEN @set_full_name::boolean = TRUE THEN @full_name
        ELSE full_name 
    END,
    email = CASE 
        WHEN @set_email::boolean = TRUE THEN @email
        ELSE email 
    END
WHERE
    username = @username
RETURNING *;
