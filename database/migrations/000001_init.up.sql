CREATE SCHEMA IF NOT EXISTS schema_astral;

CREATE TABLE IF NOT EXISTS schema_astral.users
(
    login TEXT PRIMARY KEY,
    password_hash TEXT
);

CREATE TABLE IF NOT EXISTS schema_astral.documents
(
    id UUID PRIMARY KEY,
    login TEXT NOT NULL REFERENCES schema_astral.users(login) ON DELETE CASCADE,
    name TEXT,
	mime TEXT,
	is_file BOOLEAN NOT NULL DEFAULT FALSE,
	is_public BOOLEAN NOT NULL DEFAULT FALSE,
	content BYTEA,
	json JSONB,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS schema_astral.documents_grants
(
    id SERIAL PRIMARY KEY,
    doc_id UUID NOT NULL REFERENCES schema_astral.documents(id) ON DELETE CASCADE,
    grantee_login TEXT NOT NULL REFERENCES schema_astral.users(login) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_documents_owner ON schema_astral.documents(login);
CREATE INDEX IF NOT EXISTS idx_doc_grants_doc_id ON schema_astral.documents_grants(doc_id);
CREATE INDEX IF NOT EXISTS idx_doc_grants_grantee ON schema_astral.documents_grants(grantee_login);


