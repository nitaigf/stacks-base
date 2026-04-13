package services

import (
	"context"
	"net/http"
	"testing"

	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
)

func TestAuditServiceListRequiresAdmin(t *testing.T) {
	repo := newMemoryRepository()
	user := repo.seedUser(repositories.User{ID: "user-1", Role: "member", Status: "active"})
	service := NewAuditService(repo)

	_, _, err := service.ListAuditLogs(context.Background(), RequestMetadata{ActorUserID: &user.ID}, schemas.AuditLogListQuery{})
	assertAppErrorCode(t, err, http.StatusForbidden, "forbidden")
}

func TestAuditServiceListReturnsSuccessForAdmin(t *testing.T) {
	repo := newMemoryRepository()
	admin := repo.seedUser(repositories.User{ID: "admin-1", Role: "admin", Status: "active"})
	service := NewAuditService(repo)

	_, _, err := service.ListAuditLogs(context.Background(), adminRequest(admin.ID), schemas.AuditLogListQuery{})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}
