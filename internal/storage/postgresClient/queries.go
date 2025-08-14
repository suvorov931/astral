package postgresClient

const (
	querySaveUser = `INSERT INTO schema_astral.users (login, password_hash) VALUES ($1, $2)`

	queryGetPasswordHash = `SELECT password_hash FROM schema_astral.users WHERE login = $1`

	querySaveDocument = `INSERT INTO schema_astral.documents
    (id, login, name, mime, is_file, is_public, content, json, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	querySaveDocumentGrant = `INSERT INTO schema_astral.documents_grants (doc_id, grantee_login) VALUES ($1,$2)`
)
