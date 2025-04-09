-- name: GetUser :one
select username, password, secret
from users
where username = ?;