package queries

const (
	PlaceHolder = "$"
)

const (
	InsertNewUser = `INSERT INTO users (name, email, password, role, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, NOW(), NOW()) 
			 RETURNING id`
	SelectById = `SELECT name, email, role, created_at, updated_at FROM users WHERE id=($1)`
	UpdateById = `UPDATE users SET %s, updated_at = now() WHERE id = $%d`
	DeleteById = `DELETE FROM users WHERE id = $1`
)
