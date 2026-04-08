package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	DeletedBy    *string    `json:"deletedBy,omitempty"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt sql.NullTime
}

type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	Role         string
	Status       string
}

type UpdateUserParams struct {
	Name  string
	Email string
	Role  string
}

type UserListParams struct {
	Page           int
	PerPage        int
	Query          string
	Role           string
	Status         string
	IncludeDeleted bool
	IncludeAll     bool
}

type AuditLog struct {
	ID          string         `json:"id"`
	ActorUserID *string        `json:"actorUserId,omitempty"`
	ActorName   string         `json:"actorName,omitempty"`
	ActorEmail  string         `json:"actorEmail,omitempty"`
	Action      string         `json:"action"`
	Resource    string         `json:"resource"`
	ResourceID  string         `json:"resourceId,omitempty"`
	Route       string         `json:"route"`
	Method      string         `json:"method"`
	IPAddress   string         `json:"ipAddress,omitempty"`
	UserAgent   string         `json:"userAgent,omitempty"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   time.Time      `json:"createdAt"`
}

type AuditLogInput struct {
	ActorUserID *string
	ActorName   string
	ActorEmail  string
	Action      string
	Resource    string
	ResourceID  string
	Route       string
	Method      string
	IPAddress   string
	UserAgent   string
	Metadata    map[string]any
}

type AuditLogListParams struct {
	Page     int
	PerPage  int
	Query    string
	Action   string
	Resource string
}

type Repository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	UpdateUser(ctx context.Context, id string, params UpdateUserParams) (User, error)
	SetUserStatus(ctx context.Context, id string, status string) (User, error)
	SoftDeleteUser(ctx context.Context, id string, deletedBy *string) (User, error)
	RestoreUser(ctx context.Context, id string) (User, error)
	HardDeleteUser(ctx context.Context, id string) error
	UpdateUserPassword(ctx context.Context, id string, passwordHash string) (User, error)
	FindUserByEmail(ctx context.Context, email string, includeDeleted bool) (User, error)
	FindUserByID(ctx context.Context, id string, includeDeleted bool) (User, error)
	ListUsers(ctx context.Context, params UserListParams) ([]User, PaginationMeta, error)
	RecordUserLogin(ctx context.Context, id string, at time.Time) error
	SaveRefreshToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeRefreshTokensByUserID(ctx context.Context, userID string) error
	SavePasswordResetToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error
	FindPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error)
	UsePasswordResetToken(ctx context.Context, tokenHash string) error
	CreateAuditLog(ctx context.Context, input AuditLogInput) error
	ListAuditLogs(ctx context.Context, params AuditLogListParams) ([]AuditLog, PaginationMeta, error)
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

func (r *PostgresRepository) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	const query = `
		insert into users (name, email, password_hash, role, status)
		values ($1, $2, $3, $4, $5)
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(
		r.queryRowContext(ctx, query, params.Name, strings.ToLower(strings.TrimSpace(params.Email)), params.PasswordHash, params.Role, params.Status),
		&user,
	); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, id string, params UpdateUserParams) (User, error) {
	const query = `
		update users
		set name = $2, email = $3, role = $4
		where id = $1 and deleted_at is null
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(
		r.queryRowContext(ctx, query, id, strings.TrimSpace(params.Name), strings.ToLower(strings.TrimSpace(params.Email)), params.Role),
		&user,
	); err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) SetUserStatus(ctx context.Context, id string, status string) (User, error) {
	const query = `
		update users
		set status = $2
		where id = $1 and deleted_at is null
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, id, status), &user); err != nil {
		return User{}, fmt.Errorf("set user status: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) SoftDeleteUser(ctx context.Context, id string, deletedBy *string) (User, error) {
	const query = `
		update users
		set deleted_at = now(), deleted_by = $2, status = 'inactive'
		where id = $1 and deleted_at is null
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, id, deletedBy), &user); err != nil {
		return User{}, fmt.Errorf("soft delete user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) RestoreUser(ctx context.Context, id string) (User, error) {
	const query = `
		update users
		set deleted_at = null, deleted_by = null, status = 'active'
		where id = $1 and deleted_at is not null
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, id), &user); err != nil {
		return User{}, fmt.Errorf("restore user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) HardDeleteUser(ctx context.Context, id string) error {
	result, err := r.execContext(ctx, `delete from users where id = $1`, id)
	if err != nil {
		return fmt.Errorf("hard delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected on hard delete user: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepository) UpdateUserPassword(ctx context.Context, id string, passwordHash string) (User, error) {
	const query = `
		update users
		set password_hash = $2
		where id = $1
		returning id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
	`

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, id, passwordHash), &user); err != nil {
		return User{}, fmt.Errorf("update user password: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string, includeDeleted bool) (User, error) {
	query := `
		select id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
		from users
		where lower(email) = lower($1)
	`
	if !includeDeleted {
		query += " and deleted_at is null"
	}

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, strings.ToLower(strings.TrimSpace(email))), &user); err != nil {
		return User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserByID(ctx context.Context, id string, includeDeleted bool) (User, error) {
	query := `
		select id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
		from users
		where id = $1
	`
	if !includeDeleted {
		query += " and deleted_at is null"
	}

	var user User
	if err := scanUser(r.queryRowContext(ctx, query, id), &user); err != nil {
		return User{}, fmt.Errorf("find user by id: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) ListUsers(ctx context.Context, params UserListParams) ([]User, PaginationMeta, error) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0, 8)
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if query := strings.TrimSpace(params.Query); query != "" {
		pattern := "%" + strings.ToLower(query) + "%"
		nameArg := addArg(pattern)
		emailArg := addArg(pattern)
		conditions = append(conditions, fmt.Sprintf("(lower(name) like %s or lower(email) like %s)", nameArg, emailArg))
	}

	if role := strings.TrimSpace(params.Role); role != "" {
		conditions = append(conditions, fmt.Sprintf("role = %s", addArg(role)))
	}

	if status := strings.TrimSpace(params.Status); status != "" {
		conditions = append(conditions, fmt.Sprintf("status = %s", addArg(status)))
	}

	if !params.IncludeDeleted {
		conditions = append(conditions, "deleted_at is null")
	}

	whereClause := strings.Join(conditions, " and ")
	countQuery := "select count(*) from users where " + whereClause

	var total int
	if err := r.queryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("count users: %w", err)
	}

	page, perPage := normalizePagination(params.Page, params.PerPage)
	meta := buildPaginationMeta(page, perPage, total)

	dataQuery := `
		select id, name, email, password_hash, role, status, deleted_at, deleted_by, last_login_at, created_at, updated_at
		from users
		where ` + whereClause + `
		order by created_at desc
	`

	if !params.IncludeAll {
		dataQuery += fmt.Sprintf(" limit %s offset %s", addArg(perPage), addArg((page-1)*perPage))
	}

	rows, err := r.queryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := scanUser(rows, &user); err != nil {
			return nil, PaginationMeta{}, fmt.Errorf("scan listed user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("iterate listed users: %w", err)
	}

	if params.IncludeAll {
		meta.Page = 1
		meta.PerPage = len(users)
		if len(users) == 0 {
			meta.PerPage = 0
		}
		meta.Total = len(users)
		if len(users) == 0 {
			meta.TotalPages = 0
		} else {
			meta.TotalPages = 1
		}
	}

	return users, meta, nil
}

func (r *PostgresRepository) RecordUserLogin(ctx context.Context, id string, at time.Time) error {
	if _, err := r.execContext(ctx, `update users set last_login_at = $2 where id = $1`, id, at); err != nil {
		return fmt.Errorf("record user login: %w", err)
	}

	return nil
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

func (r *PostgresRepository) RevokeRefreshTokensByUserID(ctx context.Context, userID string) error {
	if _, err := r.execContext(ctx, `
		update refresh_tokens
		set revoked_at = now()
		where user_id = $1 and revoked_at is null
	`, userID); err != nil {
		return fmt.Errorf("revoke refresh tokens by user id: %w", err)
	}

	return nil
}

func (r *PostgresRepository) SavePasswordResetToken(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	if _, err := r.execContext(ctx, `
		insert into password_reset_tokens (user_id, token_hash, expires_at)
		values ($1, $2, $3)
	`, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("save password reset token: %w", err)
	}

	return nil
}

func (r *PostgresRepository) FindPasswordResetToken(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	var token PasswordResetToken
	var usedAt sql.NullTime

	if err := r.queryRowContext(ctx, `
		select id, user_id, token_hash, expires_at, used_at, created_at
		from password_reset_tokens
		where token_hash = $1
	`, tokenHash).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&usedAt,
		&token.CreatedAt,
	); err != nil {
		return PasswordResetToken{}, fmt.Errorf("find password reset token: %w", err)
	}

	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}

	return token, nil
}

func (r *PostgresRepository) UsePasswordResetToken(ctx context.Context, tokenHash string) error {
	result, err := r.execContext(ctx, `
		update password_reset_tokens
		set used_at = now()
		where token_hash = $1 and used_at is null
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("use password reset token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected on use password reset token: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepository) CreateAuditLog(ctx context.Context, input AuditLogInput) error {
	payload, err := json.Marshal(input.Metadata)
	if err != nil {
		return fmt.Errorf("marshal audit metadata: %w", err)
	}

	const query = `
		insert into audit_logs (
			actor_user_id,
			actor_name,
			actor_email,
			action,
			resource,
			resource_id,
			route,
			method,
			ip_address,
			user_agent,
			metadata
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	if _, err := r.execContext(
		ctx,
		query,
		input.ActorUserID,
		nullableString(input.ActorName),
		nullableString(input.ActorEmail),
		input.Action,
		input.Resource,
		nullableString(input.ResourceID),
		coalesceString(input.Route, "unknown"),
		coalesceString(input.Method, "SYSTEM"),
		nullableString(input.IPAddress),
		nullableString(input.UserAgent),
		payload,
	); err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}

	return nil
}

func (r *PostgresRepository) ListAuditLogs(ctx context.Context, params AuditLogListParams) ([]AuditLog, PaginationMeta, error) {
	conditions := []string{"1 = 1"}
	args := make([]any, 0, 8)
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if query := strings.TrimSpace(params.Query); query != "" {
		pattern := "%" + strings.ToLower(query) + "%"
		arg1 := addArg(pattern)
		arg2 := addArg(pattern)
		arg3 := addArg(pattern)
		arg4 := addArg(pattern)
		conditions = append(conditions, fmt.Sprintf("(lower(coalesce(actor_name, '')) like %s or lower(coalesce(actor_email, '')) like %s or lower(coalesce(resource_id, '')) like %s or lower(route) like %s)", arg1, arg2, arg3, arg4))
	}

	if action := strings.TrimSpace(params.Action); action != "" {
		conditions = append(conditions, fmt.Sprintf("action = %s", addArg(action)))
	}

	if resource := strings.TrimSpace(params.Resource); resource != "" {
		conditions = append(conditions, fmt.Sprintf("resource = %s", addArg(resource)))
	}

	whereClause := strings.Join(conditions, " and ")
	countQuery := "select count(*) from audit_logs where " + whereClause

	var total int
	if err := r.queryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("count audit logs: %w", err)
	}

	page, perPage := normalizePagination(params.Page, params.PerPage)
	meta := buildPaginationMeta(page, perPage, total)

	dataQuery := `
		select id, actor_user_id, actor_name, actor_email, action, resource, resource_id, route, method, ip_address, user_agent, metadata, created_at
		from audit_logs
		where ` + whereClause + `
		order by created_at desc
		limit ` + addArg(perPage) + ` offset ` + addArg((page-1)*perPage)

	rows, err := r.queryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	logs := make([]AuditLog, 0)
	for rows.Next() {
		var entry AuditLog
		if err := scanAuditLog(rows, &entry); err != nil {
			return nil, PaginationMeta{}, fmt.Errorf("scan audit log: %w", err)
		}
		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, PaginationMeta{}, fmt.Errorf("iterate audit logs: %w", err)
	}

	return logs, meta, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func scanUser(scanner interface{ Scan(dest ...any) error }, user *User) error {
	var deletedAt sql.NullTime
	var lastLoginAt sql.NullTime
	var deletedBy sql.NullString

	if err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&deletedAt,
		&deletedBy,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return err
	}

	user.DeletedAt = nullTimeToPointer(deletedAt)
	user.LastLoginAt = nullTimeToPointer(lastLoginAt)
	if deletedBy.Valid {
		value := deletedBy.String
		user.DeletedBy = &value
	}

	return nil
}

func scanAuditLog(scanner interface{ Scan(dest ...any) error }, entry *AuditLog) error {
	var actorUserID sql.NullString
	var actorName sql.NullString
	var actorEmail sql.NullString
	var resourceID sql.NullString
	var ipAddress sql.NullString
	var userAgent sql.NullString
	var metadata []byte

	if err := scanner.Scan(
		&entry.ID,
		&actorUserID,
		&actorName,
		&actorEmail,
		&entry.Action,
		&entry.Resource,
		&resourceID,
		&entry.Route,
		&entry.Method,
		&ipAddress,
		&userAgent,
		&metadata,
		&entry.CreatedAt,
	); err != nil {
		return err
	}

	if actorUserID.Valid {
		value := actorUserID.String
		entry.ActorUserID = &value
	}
	if actorName.Valid {
		entry.ActorName = actorName.String
	}
	if actorEmail.Valid {
		entry.ActorEmail = actorEmail.String
	}
	if resourceID.Valid {
		entry.ResourceID = resourceID.String
	}
	if ipAddress.Valid {
		entry.IPAddress = ipAddress.String
	}
	if userAgent.Valid {
		entry.UserAgent = userAgent.String
	}

	entry.Metadata = map[string]any{}
	if len(metadata) > 0 {
		if err := json.Unmarshal(metadata, &entry.Metadata); err != nil {
			return fmt.Errorf("unmarshal audit metadata: %w", err)
		}
	}

	return nil
}

func normalizePagination(page int, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func buildPaginationMeta(page int, perPage int, total int) PaginationMeta {
	totalPages := 0
	if total > 0 && perPage > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	return PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

func nullTimeToPointer(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	parsed := value.Time
	return &parsed
}

func nullableString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func coalesceString(value string, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
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

func (r *PostgresRepository) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if r.tx != nil {
		return r.tx.QueryContext(ctx, query, args...)
	}

	return r.db.QueryContext(ctx, query, args...)
}
