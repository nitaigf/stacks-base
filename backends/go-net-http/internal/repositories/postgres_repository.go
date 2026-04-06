package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt sql.NullTime
}

type Repository interface {
	CreateUser(ctx context.Context, name string, email string, passwordHash string) (User, error)
	FindUserByEmail(ctx context.Context, email string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
	SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	CreateAuditLog(ctx context.Context, actorUserID *string, action string, resource string, resourceID string, metadata map[string]any) error
	RunInTx(ctx context.Context, fn func(Repository) error) error
}

type PostgresRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) RunInTx(ctx context.Context, fn func(Repository) error) error {
	transaction, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	txRepo := &PostgresRepository{db: r.db, tx: transaction}
	if err := fn(txRepo); err != nil {
		if rollbackErr := transaction.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback transaction after %v: %w", err, rollbackErr)
		}
		return err
	}

	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, name string, email string, passwordHash string) (User, error) {
	const query = `
		insert into users (name, email, password_hash)
		values ($1, $2, $3)
		returning id, name, email, password_hash, role, status, created_at, updated_at
	`

	var user User
	if err := r.queryRowContext(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	const query = `
		select id, name, email, password_hash, role, status, created_at, updated_at
		from users
		where email = $1
	`

	var user User
	if err := r.queryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserByID(ctx context.Context, id string) (User, error) {
	const query = `
		select id, name, email, password_hash, role, status, created_at, updated_at
		from users
		where id = $1
	`

	var user User
	if err := r.queryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	const query = `
		insert into refresh_tokens (user_id, token_hash, expires_at)
		values ($1, $2, $3)
	`

	if _, err := r.execContext(ctx, query, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	const query = `
		update refresh_tokens
		set revoked_at = now()
		where token_hash = $1 and revoked_at is null
	`

	result, err := r.execContext(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected on revoke refresh token: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepository) CreateAuditLog(ctx context.Context, actorUserID *string, action string, resource string, resourceID string, metadata map[string]any) error {
	payload, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal audit metadata: %w", err)
	}

	const query = `
		insert into audit_logs (actor_user_id, action, resource, resource_id, metadata)
		values ($1, $2, $3, $4, $5)
	`

	if _, err := r.execContext(ctx, query, actorUserID, action, resource, resourceID, payload); err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (r *PostgresRepository) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if r.tx != nil {
		return r.tx.ExecContext(ctx, query, args...)
	}

	return r.db.ExecContext(ctx, query, args...)
}

func (r *PostgresRepository) queryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if r.tx != nil {
		return r.tx.QueryRowContext(ctx, query, args...)
	}

	return r.db.QueryRowContext(ctx, query, args...)
}