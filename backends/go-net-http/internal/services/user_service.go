package services

import (
	"context"
	"net/http"
	"strings"

	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
)

type UserService struct {
	repo repositories.Repository
}

func NewUserService(repo repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) ListUsers(ctx context.Context, request RequestMetadata, query schemas.UserListQuery) ([]repositories.User, repositories.PaginationMeta, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return nil, repositories.PaginationMeta{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	users, meta, err := s.repo.ListUsers(ctx, repositories.UserListParams{
		Page:           query.Page,
		PerPage:        query.PerPage,
		Query:          query.Query,
		Role:           query.Role,
		Status:         query.Status,
		IncludeDeleted: query.IncludeDeleted,
	})
	if err != nil {
		return nil, repositories.PaginationMeta{}, InternalError("failed to list users", err)
	}

	go createAuditLog(s.repo, request, "users.list", "users", "", map[string]any{
		"query":          query.Query,
		"role":           query.Role,
		"status":         query.Status,
		"includeDeleted": query.IncludeDeleted,
		"returned":       len(users),
	})

	return users, meta, nil
}

func (s *UserService) GetUser(ctx context.Context, request RequestMetadata, id string) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	user, err := s.repo.FindUserByID(ctx, id, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
		}
		return repositories.User{}, InternalError("failed to load user", err)
	}

	go createAuditLog(s.repo, request, "users.view", "users", user.ID, map[string]any{
		"email": user.Email,
	})

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, request RequestMetadata, input schemas.UserCreateRequest) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	if _, err := s.repo.FindUserByEmail(ctx, input.Email, false); err == nil {
		return repositories.User{}, NewAppError(http.StatusConflict, "email_taken", "email already exists", nil)
	} else if err != nil && !repositories.IsNotFound(err) {
		return repositories.User{}, InternalError("failed to inspect existing user", err)
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return repositories.User{}, InternalError("failed to hash password", err)
	}

	user, err := s.repo.CreateUser(ctx, repositories.CreateUserParams{
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.TrimSpace(input.Email),
		PasswordHash: passwordHash,
		Role:         input.Role,
		Status:       input.Status,
	})
	if err != nil {
		return repositories.User{}, InternalError("failed to create user", err)
	}

	go createAuditLog(s.repo, request, "users.create", "users", user.ID, map[string]any{
		"email":  user.Email,
		"role":   user.Role,
		"status": user.Status,
	})

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, request RequestMetadata, id string, input schemas.UserUpdateRequest) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	existing, err := s.repo.FindUserByID(ctx, id, false)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
		}
		return repositories.User{}, InternalError("failed to inspect existing user", err)
	}

	if actor.ID == existing.ID && input.Role != existing.Role {
		return repositories.User{}, NewAppError(http.StatusBadRequest, "self_role_change_forbidden", "you cannot change your own role", nil)
	}

	other, err := s.repo.FindUserByEmail(ctx, input.Email, false)
	if err == nil && other.ID != existing.ID {
		return repositories.User{}, NewAppError(http.StatusConflict, "email_taken", "email already exists", nil)
	} else if err != nil && !repositories.IsNotFound(err) {
		return repositories.User{}, InternalError("failed to inspect duplicated email", err)
	}

	user, err := s.repo.UpdateUser(ctx, id, repositories.UpdateUserParams{
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	})
	if err != nil {
		return repositories.User{}, InternalError("failed to update user", err)
	}

	go createAuditLog(s.repo, request, "users.update", "users", user.ID, map[string]any{
		"email": user.Email,
		"role":  user.Role,
	})

	return user, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, request RequestMetadata, id string) (repositories.User, error) {
	return s.setUserStatus(ctx, request, id, "inactive", "users.deactivate")
}

func (s *UserService) ReactivateUser(ctx context.Context, request RequestMetadata, id string) (repositories.User, error) {
	return s.setUserStatus(ctx, request, id, "active", "users.reactivate")
}

func (s *UserService) SoftDeleteUser(ctx context.Context, request RequestMetadata, id string) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	if actor.ID == id {
		return repositories.User{}, NewAppError(http.StatusBadRequest, "self_delete_forbidden", "you cannot soft-delete your own account", nil)
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	var user repositories.User
	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		current, err := store.FindUserByID(ctx, id, true)
		if err != nil {
			if repositories.IsNotFound(err) {
				return NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
			}
			return InternalError("failed to resolve user for soft delete", err)
		}

		if current.DeletedAt != nil {
			return NewAppError(http.StatusBadRequest, "user_already_deleted", "user is already deleted", nil)
		}

		user, err = store.SoftDeleteUser(ctx, id, &actor.ID)
		if err != nil {
			return InternalError("failed to soft-delete user", err)
		}

		if err := store.RevokeRefreshTokensByUserID(ctx, id); err != nil {
			return InternalError("failed to revoke refresh tokens for deleted user", err)
		}

		return nil
	}); err != nil {
		return repositories.User{}, err
	}

	go createAuditLog(s.repo, request, "users.soft-delete", "users", user.ID, map[string]any{
		"email": user.Email,
	})

	return user, nil
}

func (s *UserService) RestoreUser(ctx context.Context, request RequestMetadata, id string) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	user, err := s.repo.RestoreUser(ctx, id)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
		}
		return repositories.User{}, InternalError("failed to restore user", err)
	}

	go createAuditLog(s.repo, request, "users.restore", "users", user.ID, map[string]any{
		"email": user.Email,
	})

	return user, nil
}

func (s *UserService) HardDeleteUser(ctx context.Context, request RequestMetadata, id string) error {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return err
	}

	if actor.ID == id {
		return NewAppError(http.StatusBadRequest, "self_delete_forbidden", "you cannot hard-delete your own account", nil)
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	target, err := s.repo.FindUserByID(ctx, id, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
		}
		return InternalError("failed to resolve user for hard delete", err)
	}

	if err := s.repo.HardDeleteUser(ctx, id); err != nil {
		if repositories.IsNotFound(err) {
			return NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
		}
		return InternalError("failed to hard-delete user", err)
	}

	go createAuditLog(s.repo, request, "users.hard-delete", "users", id, map[string]any{
		"email": target.Email,
		"name":  target.Name,
	})

	return nil
}

func (s *UserService) ExportUsers(ctx context.Context, request RequestMetadata, query schemas.UserListQuery, action string) ([]repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return nil, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	users, _, err := s.repo.ListUsers(ctx, repositories.UserListParams{
		Page:           1,
		PerPage:        1000,
		Query:          query.Query,
		Role:           query.Role,
		Status:         query.Status,
		IncludeDeleted: query.IncludeDeleted,
		IncludeAll:     true,
	})
	if err != nil {
		return nil, InternalError("failed to load users for export", err)
	}

	go createAuditLog(s.repo, request, action, "users", "", map[string]any{
		"query":          query.Query,
		"role":           query.Role,
		"status":         query.Status,
		"includeDeleted": query.IncludeDeleted,
		"returned":       len(users),
	})

	return users, nil
}

func (s *UserService) setUserStatus(ctx context.Context, request RequestMetadata, id string, status string, auditAction string) (repositories.User, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	if actor.ID == id {
		return repositories.User{}, NewAppError(http.StatusBadRequest, "self_status_change_forbidden", "you cannot change your own status", nil)
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	var user repositories.User
	if err := s.repo.RunInTx(ctx, func(store repositories.Repository) error {
		current, err := store.FindUserByID(ctx, id, true)
		if err != nil {
			if repositories.IsNotFound(err) {
				return NewAppError(http.StatusNotFound, "user_not_found", "user was not found", err)
			}
			return InternalError("failed to resolve user for status change", err)
		}

		if current.DeletedAt != nil {
			return NewAppError(http.StatusBadRequest, "user_deleted", "deleted users must be restored first", nil)
		}

		user, err = store.SetUserStatus(ctx, id, status)
		if err != nil {
			return InternalError("failed to update user status", err)
		}

		if status != "active" {
			if err := store.RevokeRefreshTokensByUserID(ctx, id); err != nil {
				return InternalError("failed to revoke refresh tokens after status change", err)
			}
		}

		return nil
	}); err != nil {
		return repositories.User{}, err
	}

	go createAuditLog(s.repo, request, auditAction, "users", user.ID, map[string]any{
		"email":  user.Email,
		"status": user.Status,
	})

	return user, nil
}
