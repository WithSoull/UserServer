package queries

const (
	InsertNewUser = `INSERT INTO users (name, email, password_hash, role, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, NOW(), NOW()) 
			 RETURNING id`
)
