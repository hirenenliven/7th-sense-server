package tests

import (
	"context"
	"errors"
	"testing"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

func TestCompanyProfileCRUD(t *testing.T) {
	prepareProfileTable(t, "company_profiles")

	ctx := context.Background()
	companyProfileService := service.NewCompanyProfileService(repository.NewCompanyProfileRepository(testDBPool))

	t.Run("get before create returns not found", func(t *testing.T) {
		_, err := companyProfileService.GetByUserID(ctx, testUserID)
		if !errors.Is(err, response.CompanyProfileNotFound) {
			t.Fatalf("GetByUserID() error = %v; want %v", err, response.CompanyProfileNotFound)
		}
	})

	finixID := int64(1001)
	t.Run("create", func(t *testing.T) {
		created, err := companyProfileService.Create(ctx, model.CompanyProfile{
			UserID:    testUserID,
			GSTNumber: "24AAACC1206D1ZM",
			PANNumber: "AAACC1206D",
			Website:   "https://example.com",
			FinixID:   &finixID,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if created.ID == 0 || created.UserID != testUserID {
			t.Errorf("Create() id = %d, userId = %d; want non-zero id, userId %d", created.ID, created.UserID, testUserID)
		}
		if created.GSTNumber != "24AAACC1206D1ZM" || created.PANNumber != "AAACC1206D" || created.Website != "https://example.com" {
			t.Errorf("Create() = %+v; fields not saved", created)
		}
		if created.FinixID == nil || *created.FinixID != finixID {
			t.Errorf("Create() finixId = %v; want %d", created.FinixID, finixID)
		}
		if created.IsDeleted {
			t.Error("Create() isDeleted = true; want false")
		}
	})

	t.Run("read", func(t *testing.T) {
		profile, err := companyProfileService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if profile.GSTNumber != "24AAACC1206D1ZM" || profile.Website != "https://example.com" {
			t.Errorf("GetByUserID() = %+v; want the created profile", profile)
		}
	})

	t.Run("update", func(t *testing.T) {
		updated, err := companyProfileService.Update(ctx, model.CompanyProfile{
			UserID:    testUserID,
			GSTNumber: "27AAACC1206D1ZN",
			PANNumber: "", // stored as NULL, read back as ""
			Website:   "https://example.org",
			FinixID:   nil,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if updated.GSTNumber != "27AAACC1206D1ZN" || updated.PANNumber != "" || updated.Website != "https://example.org" || updated.FinixID != nil {
			t.Errorf("Update() = %+v; fields not updated", updated)
		}

		profile, err := companyProfileService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() after update error = %v", err)
		}
		if profile.GSTNumber != "27AAACC1206D1ZN" {
			t.Errorf("GetByUserID() gstNumber = %q; want updated value", profile.GSTNumber)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := companyProfileService.Delete(ctx, testUserID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if _, err := companyProfileService.GetByUserID(ctx, testUserID); !errors.Is(err, response.CompanyProfileNotFound) {
			t.Errorf("GetByUserID() after delete error = %v; want %v", err, response.CompanyProfileNotFound)
		}
		if _, err := companyProfileService.Update(ctx, model.CompanyProfile{UserID: testUserID}); !errors.Is(err, response.CompanyProfileNotFound) {
			t.Errorf("Update() after delete error = %v; want %v", err, response.CompanyProfileNotFound)
		}
		if err := companyProfileService.Delete(ctx, testUserID); !errors.Is(err, response.CompanyProfileNotFound) {
			t.Errorf("second Delete() error = %v; want %v", err, response.CompanyProfileNotFound)
		}
	})
}
