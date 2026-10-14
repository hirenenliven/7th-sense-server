package dto

import (
	"net/http"
	"strings"
	"testing"

	"7th-sense-server/internal/response"
)

// validCustomerDetails returns customer details that pass validation.
func validCustomerDetails() CustomerDetails {
	branchID := int64(1)
	return CustomerDetails{
		BranchID:      &branchID,
		CompanyName:   "Enliven Pvt Ltd",
		GSTNumber:     "24AAACE1234F1Z5",
		PANNumber:     "AAACE1234F",
		AddressLine1:  "Plot 12, GIDC",
		Country:       "India",
		State:         "Gujarat",
		City:          "Rajkot",
		ZipCode:       "360001",
		ContactPerson: "Hiren",
		PrimaryPhone:  "+919876543210",
		TimeZone:      "Asia/Kolkata",
		DayStartFrom:  "06:00",
	}
}

func TestCreateUserRequestValidate(t *testing.T) {
	valid := func() CreateUserRequest {
		return CreateUserRequest{
			Username: "hiren_c", Email: "hiren@example.com", Password: "secret123", RoleID: 4,
			CustomerDetails: validCustomerDetails(),
		}
	}

	tests := []struct {
		name    string
		modify  func(*CreateUserRequest)
		wantErr error
	}{
		{"valid", func(*CreateUserRequest) {}, nil},
		{"valid with optional fields", func(r *CreateUserRequest) {
			finixID := int64(42)
			r.FirstName, r.Website, r.CompanyLogo, r.AddressLine2, r.SecondaryPhone = "Hiren", "www.example.com", "uploads/logo.png", "Near bridge", "02812345678"
			r.FinixID, r.DayStartFrom = &finixID, "06:30:00"
		}, nil},
		{"username too short", func(r *CreateUserRequest) { r.Username = "ab" }, response.NewError(http.StatusBadRequest, "Please enter valid username")},
		{"username with space", func(r *CreateUserRequest) { r.Username = "hiren c" }, response.NewError(http.StatusBadRequest, "Please enter valid username")},
		{"email missing", func(r *CreateUserRequest) { r.Email = "" }, response.NewError(http.StatusBadRequest, "Please enter valid email")},
		{"email invalid", func(r *CreateUserRequest) { r.Email = "not-an-email" }, response.NewError(http.StatusBadRequest, "Please enter valid email")},
		{"email with display name", func(r *CreateUserRequest) { r.Email = "Hiren <hiren@example.com>" }, response.NewError(http.StatusBadRequest, "Please enter valid email")},
		{"password too short", func(r *CreateUserRequest) { r.Password = "short" }, response.NewError(http.StatusBadRequest, "Please enter valid password")},
		{"password too long", func(r *CreateUserRequest) { r.Password = strings.Repeat("a", 73) }, response.NewError(http.StatusBadRequest, "Please enter valid password")},
		{"first name too long", func(r *CreateUserRequest) { r.FirstName = strings.Repeat("a", 101) }, response.NewError(http.StatusBadRequest, "Please enter valid name")},
		{"role missing defaults to admin", func(r *CreateUserRequest) { r.RoleID = 0 }, nil},
		{"role negative", func(r *CreateUserRequest) { r.RoleID = -1 }, response.NewError(http.StatusBadRequest, "Please enter valid role")},
		{"branch missing", func(r *CreateUserRequest) { r.BranchID = nil }, response.NewError(http.StatusBadRequest, "Please enter valid branch")},
		{"company name missing", func(r *CreateUserRequest) { r.CompanyName = " " }, response.NewError(http.StatusBadRequest, "Please enter valid company name")},
		{"finix id negative", func(r *CreateUserRequest) { finixID := int64(-1); r.FinixID = &finixID }, response.NewError(http.StatusBadRequest, "Please enter valid Finix ID")},
		{"gst missing", func(r *CreateUserRequest) { r.GSTNumber = "" }, response.NewError(http.StatusBadRequest, "Please enter valid GST number")},
		{"gst invalid", func(r *CreateUserRequest) { r.GSTNumber = "24AAACE1234F1X5" }, response.NewError(http.StatusBadRequest, "Please enter valid GST number")},
		{"pan missing", func(r *CreateUserRequest) { r.PANNumber = "" }, response.NewError(http.StatusBadRequest, "Please enter valid PAN number")},
		{"website invalid", func(r *CreateUserRequest) { r.Website = "not a site" }, response.NewError(http.StatusBadRequest, "Please enter valid website")},
		{"website ftp", func(r *CreateUserRequest) { r.Website = "ftp://example.com" }, response.NewError(http.StatusBadRequest, "Please enter valid website")},
		{"address missing", func(r *CreateUserRequest) { r.AddressLine1 = "" }, response.NewError(http.StatusBadRequest, "Please enter valid address")},
		{"country missing", func(r *CreateUserRequest) { r.Country = " " }, response.NewError(http.StatusBadRequest, "Please enter valid country")},
		{"state missing", func(r *CreateUserRequest) { r.State = "" }, response.NewError(http.StatusBadRequest, "Please enter valid state")},
		{"city missing", func(r *CreateUserRequest) { r.City = "" }, response.NewError(http.StatusBadRequest, "Please enter valid city")},
		{"zip missing", func(r *CreateUserRequest) { r.ZipCode = "" }, response.NewError(http.StatusBadRequest, "Please enter valid zip code")},
		{"zip too long", func(r *CreateUserRequest) { r.ZipCode = strings.Repeat("1", 21) }, response.NewError(http.StatusBadRequest, "Please enter valid zip code")},
		{"contact person missing", func(r *CreateUserRequest) { r.ContactPerson = "" }, response.NewError(http.StatusBadRequest, "Please enter valid contact person")},
		{"primary phone missing", func(r *CreateUserRequest) { r.PrimaryPhone = "" }, response.NewError(http.StatusBadRequest, "Please enter valid mobile")},
		{"secondary phone invalid", func(r *CreateUserRequest) { r.SecondaryPhone = "12ab" }, response.NewError(http.StatusBadRequest, "Please enter valid mobile")},
		{"time zone missing", func(r *CreateUserRequest) { r.TimeZone = "" }, response.NewError(http.StatusBadRequest, "Please enter valid time zone")},
		{"time zone Local", func(r *CreateUserRequest) { r.TimeZone = "Local" }, response.NewError(http.StatusBadRequest, "Please enter valid time zone")},
		{"time zone unknown", func(r *CreateUserRequest) { r.TimeZone = "Mars/Base" }, response.NewError(http.StatusBadRequest, "Please enter valid time zone")},
		{"day start missing", func(r *CreateUserRequest) { r.DayStartFrom = "" }, response.NewError(http.StatusBadRequest, "Please enter valid day start time")},
		{"day start invalid", func(r *CreateUserRequest) { r.DayStartFrom = "25:00" }, response.NewError(http.StatusBadRequest, "Please enter valid day start time")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := valid()
			tt.modify(&request)
			err := request.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, want %v", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestCreateUserRequestValidateNormalizes(t *testing.T) {
	request := CreateUserRequest{
		Username: "  hiren_c ", Email: " Hiren@Example.COM ", Password: "secret123", RoleID: 4,
		CustomerDetails: validCustomerDetails(),
	}
	request.GSTNumber, request.PANNumber, request.City = " 24aaace1234f1z5 ", "aaace1234f", " Rajkot "
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if request.Username != "hiren_c" || request.Email != "hiren@example.com" {
		t.Errorf("got username %q, email %q; want trimmed and lowercase email", request.Username, request.Email)
	}
	if request.GSTNumber != "24AAACE1234F1Z5" || request.PANNumber != "AAACE1234F" || request.City != "Rajkot" {
		t.Errorf("got gst %q, pan %q, city %q; want trimmed and uppercase GST/PAN", request.GSTNumber, request.PANNumber, request.City)
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
		{"deleted", response.NewError(http.StatusBadRequest, "Please enter valid status")},
	}

	for _, tt := range tests {
		t.Run("status "+tt.status, func(t *testing.T) {
			request := UpdateUserRequest{
				Username: "hiren_c", Email: "hiren@example.com", RoleID: 4, Status: tt.status,
				CustomerDetails: validCustomerDetails(),
			}
			err := request.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() error = nil, want %v", tt.wantErr)
				} else if err.Error() != tt.wantErr.Error() {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestParseDayStart(t *testing.T) {
	dayStart, err := ParseDayStart("06:30")
	if err != nil {
		t.Fatalf("ParseDayStart() error = %v", err)
	}
	if dayStart.Hour() != 6 || dayStart.Minute() != 30 {
		t.Errorf("ParseDayStart() = %v; want 06:30", dayStart)
	}
}

func TestCustomerFiltersValidate(t *testing.T) {
	request := CreateUserRequest{
		Username: "hiren_c", Email: "hiren@example.com", Password: "secret123", RoleID: 4,
		CustomerDetails: validCustomerDetails(),
	}
	request.Filters = []string{" Country", "state", "SHOP", "country"}
	if err := request.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if strings.Join(request.Filters, ",") != "country,state,shop" {
		t.Errorf("Filters = %v; want [country state shop] (lowercase, no duplicates)", request.Filters)
	}

	request.Filters = nil
	if err := request.Validate(); err != nil || request.Filters == nil {
		t.Errorf("Validate() with no filters: error = %v, filters = %#v; want nil error and empty list", err, request.Filters)
	}

	request.Filters = []string{"country", "machine"}
	if err := request.Validate(); err == nil || err.Error() != "Please enter valid filter" {
		t.Errorf("Validate() with unknown filter error = %v; want Please enter valid filter", err)
	}
}

func TestUpdateThemeRequestValidate(t *testing.T) {
	for _, testCase := range []struct {
		theme     string
		wantTheme string
		wantErr   bool
	}{
		{ThemeDark, ThemeDark, false},
		{" Light_Version ", ThemeLight, false},
		{"", "", true},
		{"blue_version", "", true},
	} {
		request := UpdateThemeRequest{Theme: testCase.theme}
		err := request.Validate()
		if (err != nil) != testCase.wantErr {
			t.Errorf("Validate(%q) error = %v; want error %v", testCase.theme, err, testCase.wantErr)
			continue
		}
		if err == nil && request.Theme != testCase.wantTheme {
			t.Errorf("Validate(%q) theme = %q; want %q", testCase.theme, request.Theme, testCase.wantTheme)
		}
	}
}
