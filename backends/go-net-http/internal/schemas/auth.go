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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
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

func (r ForgotPasswordRequest) Validate() map[string]string {
	errors := map[string]string{}

	if !strings.Contains(strings.TrimSpace(r.Email), "@") {
		errors["email"] = "email must be valid"
	}

	return errors
}

func (r ResetPasswordRequest) Validate() map[string]string {
	errors := map[string]string{}

	if strings.TrimSpace(r.Token) == "" {
		errors["token"] = "token is required"
	}

	if len(r.NewPassword) < 8 {
		errors["newPassword"] = "newPassword must have at least 8 characters"
	}

	return errors
}

func (r ChangePasswordRequest) Validate() map[string]string {
	errors := map[string]string{}

	if len(r.CurrentPassword) < 8 {
		errors["currentPassword"] = "currentPassword must have at least 8 characters"
	}

	if len(r.NewPassword) < 8 {
		errors["newPassword"] = "newPassword must have at least 8 characters"
	}

	return errors
}
