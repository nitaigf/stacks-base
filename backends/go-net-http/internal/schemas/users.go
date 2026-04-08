package schemas

import (
	"strconv"
	"strings"
)

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UserUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserListQuery struct {
	Page           int
	PerPage        int
	Query          string
	Role           string
	Status         string
	IncludeDeleted bool
}

type AuditLogListQuery struct {
	Page     int
	PerPage  int
	Query    string
	Action   string
	Resource string
}

func (r UserCreateRequest) Validate() map[string]string {
	errors := map[string]string{}

	if len(strings.TrimSpace(r.Name)) < 2 {
		errors["name"] = "name must have at least 2 characters"
	}

	if !strings.Contains(strings.TrimSpace(r.Email), "@") {
		errors["email"] = "email must be valid"
	}

	if len(r.Password) < 8 {
		errors["password"] = "password must have at least 8 characters"
	}

	if r.Role != "admin" && r.Role != "member" {
		errors["role"] = "role must be admin or member"
	}

	if r.Status != "active" && r.Status != "inactive" {
		errors["status"] = "status must be active or inactive"
	}

	return errors
}

func (r UserUpdateRequest) Validate() map[string]string {
	errors := map[string]string{}

	if len(strings.TrimSpace(r.Name)) < 2 {
		errors["name"] = "name must have at least 2 characters"
	}

	if !strings.Contains(strings.TrimSpace(r.Email), "@") {
		errors["email"] = "email must be valid"
	}

	if r.Role != "admin" && r.Role != "member" {
		errors["role"] = "role must be admin or member"
	}

	return errors
}

func ParseUserListQuery(raw map[string]string) UserListQuery {
	return UserListQuery{
		Page:           parseIntDefault(raw["page"], 1),
		PerPage:        parseIntDefault(raw["perPage"], 10),
		Query:          strings.TrimSpace(raw["query"]),
		Role:           strings.TrimSpace(raw["role"]),
		Status:         strings.TrimSpace(raw["status"]),
		IncludeDeleted: strings.EqualFold(strings.TrimSpace(raw["includeDeleted"]), "true"),
	}
}

func ParseAuditLogListQuery(raw map[string]string) AuditLogListQuery {
	return AuditLogListQuery{
		Page:     parseIntDefault(raw["page"], 1),
		PerPage:  parseIntDefault(raw["perPage"], 20),
		Query:    strings.TrimSpace(raw["query"]),
		Action:   strings.TrimSpace(raw["action"]),
		Resource: strings.TrimSpace(raw["resource"]),
	}
}

func parseIntDefault(raw string, fallback int) int {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return fallback
	}

	return parsed
}
