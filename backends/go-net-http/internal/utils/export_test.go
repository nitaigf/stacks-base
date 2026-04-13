package utils

import (
	"strings"
	"testing"
	"time"

	"stacks-base/backends/go-net-http/internal/repositories"
)

func TestBuildUsersCSVUsesSemicolonDelimiter(t *testing.T) {
	createdAt := time.Date(2026, 4, 8, 10, 0, 0, 0, time.UTC)
	users := []repositories.User{
		{
			ID:        "user-1",
			Name:      "Admin",
			Email:     "admin@example.com",
			Role:      "admin",
			Status:    "active",
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
	}

	payload, err := BuildUsersCSV(users)
	if err != nil {
		t.Fatalf("expected csv generation to succeed: %v", err)
	}

	output := string(payload)
	if !strings.Contains(output, "ID;Nome;Email;Papel;Status;Excluido em;Ultimo login;Criado em;Atualizado em") {
		t.Fatalf("expected semicolon-separated header, got %q", output)
	}
	if strings.Contains(output, "ID,Nome,Email") {
		t.Fatalf("expected csv output not to use comma as delimiter, got %q", output)
	}
}

func TestBuildUsersPDFRendersTablePreview(t *testing.T) {
	loginAt := time.Date(2026, 4, 8, 12, 30, 0, 0, time.UTC)
	users := []repositories.User{
		{
			ID:          "user-1",
			Name:        "Admin",
			Email:       "admin@example.com",
			Role:        "admin",
			Status:      "active",
			LastLoginAt: &loginAt,
			CreatedAt:   loginAt,
			UpdatedAt:   loginAt,
		},
	}

	payload, err := BuildUsersPDF(users)
	if err != nil {
		t.Fatalf("expected pdf generation to succeed: %v", err)
	}

	output := string(payload)
	if !strings.Contains(output, "%PDF-1.4") {
		t.Fatalf("expected valid pdf header, got %q", output[:min(len(output), 32)])
	}
	if !strings.Contains(output, "Relatorio de usuarios") {
		t.Fatalf("expected title in pdf output")
	}
	if !strings.Contains(output, "(Nome) Tj") || !strings.Contains(output, "(Email) Tj") {
		t.Fatalf("expected tabular header cells in pdf output")
	}
	if !strings.Contains(output, " re S") {
		t.Fatalf("expected stroked table cells in pdf output")
	}
}
