-- name: GetUser :one
select username, secret
from users
where username = ?;