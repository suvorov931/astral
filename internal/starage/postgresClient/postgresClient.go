package postgresClient

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func New(ctx context.Context, config *Config, logger *zap.Logger, migrationsPath string) (*PostgresService, error) {
	if config.Timeout == 0 {
		config.Timeout = DefaultPostgresTimeout
	}

	url := buildURL(config)
	dsn := buildDSN(config)

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	err = upMigration(url, migrationsPath)
	if err != nil {
		return nil, err
	}

	return &PostgresService{
		pool:    pool,
		logger:  logger,
		timeout: config.Timeout,
	}, nil
}

func (ps *PostgresService) SaveUser(ctx context.Context, login string, passwordHash string) error {
	ctx, cancel := context.WithTimeout(ctx, ps.timeout)
	defer cancel()

	q := `INSERT INTO schema_users.users (login, password_hash) VALUES ($1, $2)`
	tag, err := ps.pool.Exec(ctx, q, login, passwordHash)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				ps.logger.Warn("SaveUser: duplicate login")
				return ErrDuplicateLogin
			}
		}

		ps.logger.Error("SaveUser: failed to save user into database:", zap.Error(err))
		return fmt.Errorf("SaveUser: failed to save user into database: %w", err)
	}

	if tag.RowsAffected() == 0 {
		ps.logger.Error("SaveUser: no rows affected")
		return fmt.Errorf("SaveUser: no rows affected")
	}

	ps.logger.Info("SaveUser: successfully save user")

	return nil
}

func (ps *PostgresService) Close() {
	ps.pool.Close()
}

func buildURL(config *Config) string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	return url
}

func buildDSN(config *Config) string {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%d pool_min_conns=%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.MaxConns,
		config.MinConns,
	)

	return dsn
}

func upMigration(url string, path string) error {
	migration, err := migrate.New(path, url)
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
}
