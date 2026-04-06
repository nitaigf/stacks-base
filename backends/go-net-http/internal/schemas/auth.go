package schemas

import "strings"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() map[string]string {
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

	return errors
}

func (r LoginRequest) Validate() map[string]string {
	errors := map[string]string{}

	if !strings.Contains(strings.TrimSpace(r.Email), "@") {
		errors["email"] = "email must be valid"
	}

	if len(r.Password) < 8 {
		errors["password"] = "password must have at least 8 characters"
	}

	return errors
}