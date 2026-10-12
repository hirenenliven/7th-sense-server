package handler

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
)

func TestPathID(t *testing.T) {
	tests := []struct {
		pathValue string
		wantID    int64
		wantErr   error
	}{
		{"1", 1, nil},
		{"9223372036854775807", 9223372036854775807, nil},
		{"0", 0, response.InvalidID},
		{"-5", 0, response.InvalidID},
		{"abc", 0, response.InvalidID},
		{"", 0, response.InvalidID},
		{"9223372036854775808", 0, response.InvalidID}, // overflows int64
	}

	for _, tt := range tests {
		t.Run(tt.pathValue, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			request.SetPathValue("id", tt.pathValue)

			id, err := pathID(request)
			if id != tt.wantID || !errors.Is(err, tt.wantErr) {
				t.Errorf("pathID() = %d, %v; want %d, %v", id, err, tt.wantID, tt.wantErr)
			}
		})
	}
}

func TestDecodeAndValidate(t *testing.T) {
	validBody := `{"username":"hiren_c","email":"hiren@example.com","password":"secret123"}`

	tests := []struct {
		name    string
		body    string
		wantErr error
	}{
		{"valid", validBody, nil},
		{"malformed json", `{"username":`, response.InvalidJSON},
		{"empty body", ``, response.InvalidJSON},
		{"fails validation", `{"username":"ab"}`, response.ValidUserName},
		{"too large", `{"username":"` + strings.Repeat("a", maxBodyBytes) + `"}`, response.RequestEntryLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			var requestBody dto.CreateUserRequest

			err := decodeAndValidate(httptest.NewRecorder(), request, &requestBody)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("decodeAndValidate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
