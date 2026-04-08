package services

import (
	"context"

	"stacks-base/backends/go-net-http/internal/repositories"
	"stacks-base/backends/go-net-http/internal/schemas"
)

type AuditService struct {
	repo repositories.Repository
}

func NewAuditService(repo repositories.Repository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) ListAuditLogs(ctx context.Context, request RequestMetadata, query schemas.AuditLogListQuery) ([]repositories.AuditLog, repositories.PaginationMeta, error) {
	actor, err := requireAdminActor(ctx, s.repo, request)
	if err != nil {
		return nil, repositories.PaginationMeta{}, err
	}

	request.ActorUserID = &actor.ID
	request.ActorName = actor.Name
	request.ActorEmail = actor.Email
	request.ActorRole = actor.Role

	logs, meta, err := s.repo.ListAuditLogs(ctx, repositories.AuditLogListParams{
		Page:     query.Page,
		PerPage:  query.PerPage,
		Query:    query.Query,
		Action:   query.Action,
		Resource: query.Resource,
	})
	if err != nil {
		return nil, repositories.PaginationMeta{}, InternalError("failed to list audit logs", err)
	}

	go createAuditLog(s.repo, request, "audit-logs.list", "audit_logs", "", map[string]any{
		"query":    query.Query,
		"action":   query.Action,
		"resource": query.Resource,
		"returned": len(logs),
	})

	return logs, meta, nil
}
