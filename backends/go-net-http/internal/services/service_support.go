package services

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"stacks-base/backends/go-net-http/internal/repositories"
)

func resolveActiveActor(ctx context.Context, repo repositories.Repository, request RequestMetadata) (repositories.User, error) {
	if request.ActorUserID == nil || strings.TrimSpace(*request.ActorUserID) == "" {
		return repositories.User{}, NewAppError(http.StatusUnauthorized, "unauthorized", "authentication is required", nil)
	}

	actor, err := repo.FindUserByID(ctx, *request.ActorUserID, true)
	if err != nil {
		if repositories.IsNotFound(err) {
			return repositories.User{}, NewAppError(http.StatusUnauthorized, "unauthorized", "authentication is required", err)
		}
		return repositories.User{}, InternalError("failed to resolve actor", err)
	}

	if err := ensureUserAccessible(actor); err != nil {
		return repositories.User{}, err
	}

	return actor, nil
}

func requireAdminActor(ctx context.Context, repo repositories.Repository, request RequestMetadata) (repositories.User, error) {
	actor, err := resolveActiveActor(ctx, repo, request)
	if err != nil {
		return repositories.User{}, err
	}

	if actor.Role != "admin" {
		return repositories.User{}, NewAppError(http.StatusForbidden, "forbidden", "admin access is required", nil)
	}

	return actor, nil
}

func createAuditLog(repo repositories.Repository, request RequestMetadata, action string, resource string, resourceID string, metadata map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	input := repositories.AuditLogInput{
		ActorUserID: request.ActorUserID,
		ActorName:   request.ActorName,
		ActorEmail:  request.ActorEmail,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Route:       request.Route,
		Method:      request.Method,
		IPAddress:   request.IPAddress,
		UserAgent:   request.UserAgent,
		Metadata:    metadata,
	}

	if input.ActorUserID != nil && (strings.TrimSpace(input.ActorName) == "" || strings.TrimSpace(input.ActorEmail) == "") {
		if actor, err := repo.FindUserByID(ctx, *input.ActorUserID, true); err == nil {
			if strings.TrimSpace(input.ActorName) == "" {
				input.ActorName = actor.Name
			}
			if strings.TrimSpace(input.ActorEmail) == "" {
				input.ActorEmail = actor.Email
			}
		}
	}

	if err := repo.CreateAuditLog(ctx, input); err != nil {
		log.Printf("failed to create audit log %s: %v", action, err)
	}
}
