package dto

import (
	"errors"
	"strings"
	"testing"

	"7th-sense-server/internal/response"
)

func TestCreateUserRequestValidate(t *testing.T) {
	valid := func() CreateUserRequest {
		return CreateUserRequest{Username: "hiren_c", Email: "hiren@example.com", Password: "secret123"}
	}

	tests := []struct {
		name    string
		modify  func(*CreateUserRequest)
		wantErr error
	}{
		{"valid", func(*CreateUserRequest) {}, nil},
		{"valid with optional fields", func(r *CreateUserRequest) { r.FirstName, r.Mobile = "Hiren", "+919876543210" }, nil},
		{"username too short", func(r *CreateUserRequest) { r.Username = "ab" }, response.ValidUserName},
		{"username with space", func(r *CreateUserRequest) { r.Username = "hiren c" }, response.ValidUserName},
		{"email missing", func(r *CreateUserRequest) { r.Email = "" }, response.ValidEmail},
		{"email invalid", func(r *CreateUserRequest) { r.Email = "not-an-email" }, response.ValidEmail},
		{"email with display name", func(r *CreateUserRequest) { r.Email = "Hiren <hiren@example.com>" }, response.ValidEmail},
		{"password too short", func(r *CreateUserRequest) { r.Password = "short" }, response.ValidPassword},
		{"password too long", func(r *CreateUserRequest) { r.Password = strings.Repeat("a", 73) }, response.ValidPassword},
		{"first name too long", func(r *CreateUserRequest) { r.FirstName = strings.Repeat("a", 101) }, response.ValidName},
		{"mobile with letters", func(r *CreateUserRequest) { r.Mobile = "98765abc" }, response.ValidMobile},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := valid()
			tt.modify(&request)
			if err := request.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUserRequestValidateNormalizes(t *testing.T) {
	request := CreateUserRequest{Username: "  hiren_c ", Email: " Hiren@Example.COM ", Password: "secret123"}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if request.Username != "hiren_c" || request.Email != "hiren@example.com" {
		t.Errorf("got username %q, email %q; want trimmed and lowercase email", request.Username, request.Email)
	}
}

func TestUpdateUserRequestValidateStatus(t *testing.T) {
	tests := []struct {
		status  string
		wantErr error
	}{
		{"", nil},
		{"active", nil},
		{"INACTIVE", nil},
		{"deleted", response.ValidStatus},
	}

	for _, tt := range tests {
		t.Run("status "+tt.status, func(t *testing.T) {
			request := UpdateUserRequest{Username: "hiren_c", Email: "hiren@example.com", Status: tt.status}
			if err := request.Validate(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
