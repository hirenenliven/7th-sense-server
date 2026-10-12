package response

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestSendResponse(t *testing.T) {
	tests := []struct {
		name     string
		result   error
		data     any
		wantCode int
		wantBody string
	}{
		{
			name: "success with data", result: Success, data: map[string]int{"id": 1},
			wantCode: 200, wantBody: `{"statusCode":200,"success":true,"message":"success","data":{"id":1}}`,
		},
		{
			name: "created", result: Created, data: []string{"a"},
			wantCode: 201, wantBody: `{"statusCode":201,"success":true,"message":"created","data":["a"]}`,
		},
		{
			name: "string replaces message", result: Success, data: "User deleted",
			wantCode: 200, wantBody: `{"statusCode":200,"success":true,"message":"User deleted"}`,
		},
		{
			name: "predefined failure", result: UserNotFound,
			wantCode: 404, wantBody: `{"statusCode":404,"success":false,"message":"User not found"}`,
		},
		{
			name: "wrapped predefined failure", result: fmt.Errorf("create: %w", EmailExists),
			wantCode: 400, wantBody: `{"statusCode":400,"success":false,"message":"Email already exists"}`,
		},
		{
			name: "unknown error is hidden", result: errors.New("pq: connection refused"),
			wantCode: 500, wantBody: `{"statusCode":500,"success":false,"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			SendResponse(recorder, httptest.NewRequest("GET", "/", nil), tt.result, tt.data)

			if recorder.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantCode)
			}
			if got := recorder.Body.String(); got != tt.wantBody+"\n" {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
		})
	}
}
