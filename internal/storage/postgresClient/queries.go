package postgresClient

const (
	querySaveUser = `INSERT INTO schema_astral.users (login, password_hash) VALUES ($1, $2)`

	queryGetPasswordHash = `SELECT password_hash FROM schema_astral.users WHERE login = $1`
)
