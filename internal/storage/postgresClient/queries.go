package postgresClient

const (
	querySaveUser = `INSERT INTO schema_users.users (login, password_hash) VALUES ($1, $2)`

	queryIsUserExists = `SELECT EXISTS (SELECT 1 FROM schema_users.users WHERE login = $1)`

	queryGetPasswordHash = `SELECT password_hash FROM schema_users.users WHERE login = $1`
)
