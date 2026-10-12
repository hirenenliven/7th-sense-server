package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

func TestUserProfileCRUD(t *testing.T) {
	prepareProfileTable(t, "user_profiles")

	ctx := context.Background()
	userProfileService := service.NewUserProfileService(repository.NewUserProfileRepository(testDBPool))

	birthDate := time.Date(1995, time.March, 14, 0, 0, 0, 0, time.UTC)
	joiningDate := time.Date(2022, time.July, 1, 0, 0, 0, 0, time.UTC)
	startDayTime := time.Date(2000, time.January, 1, 9, 30, 0, 0, time.UTC)

	t.Run("get before create returns not found", func(t *testing.T) {
		_, err := userProfileService.GetByUserID(ctx, testUserID)
		if !errors.Is(err, response.UserProfileNotFound) {
			t.Fatalf("GetByUserID() error = %v; want %v", err, response.UserProfileNotFound)
		}
	})

	t.Run("create", func(t *testing.T) {
		created, err := userProfileService.Create(ctx, model.UserProfile{
			UserID:         testUserID,
			Designation:    "Software Engineer",
			Grade:          "G2",
			BirthDate:      &birthDate,
			JoiningDate:    &joiningDate,
			BloodGroup:     "O+",
			ProfilePicture: "/uploads/profile/6.png",
			ContactPerson:  "Test Contact",
			ContactPhone:   "9876543210",
			Theme:          "dark_version",
			Timezone:       "Asia/Kolkata",
			Language:       "en",
			StartDayTime:   &startDayTime,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if created.ID == 0 || created.UserID != testUserID {
			t.Errorf("Create() id = %d, userId = %d; want non-zero id, userId %d", created.ID, created.UserID, testUserID)
		}
		if created.Designation != "Software Engineer" || created.Grade != "G2" || created.BloodGroup != "O+" ||
			created.ContactPhone != "9876543210" || created.Theme != "dark_version" || created.Language != "en" {
			t.Errorf("Create() = %+v; fields not saved", created)
		}
		if created.BirthDate == nil || !sameDate(*created.BirthDate, birthDate) {
			t.Errorf("Create() birthDate = %v; want %v", created.BirthDate, birthDate)
		}
		if created.JoiningDate == nil || !sameDate(*created.JoiningDate, joiningDate) {
			t.Errorf("Create() joiningDate = %v; want %v", created.JoiningDate, joiningDate)
		}
		if created.StartDayTime == nil || !sameClock(*created.StartDayTime, startDayTime) {
			t.Errorf("Create() startDayTime = %v; want %v", created.StartDayTime, startDayTime)
		}
		if created.IsDeleted {
			t.Error("Create() isDeleted = true; want false")
		}
	})

	t.Run("read", func(t *testing.T) {
		profile, err := userProfileService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() error = %v", err)
		}
		if profile.Designation != "Software Engineer" || profile.Timezone != "Asia/Kolkata" {
			t.Errorf("GetByUserID() = %+v; want the created profile", profile)
		}
	})

	t.Run("update", func(t *testing.T) {
		updated, err := userProfileService.Update(ctx, model.UserProfile{
			UserID:      testUserID,
			Designation: "Senior Software Engineer",
			Grade:       "G3",
			BirthDate:   &birthDate,
			JoiningDate: nil, // cleared
			BloodGroup:  "O+",
			Theme:       "light_version",
			Timezone:    "Asia/Kolkata",
			Language:    "hi",
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		if updated.Designation != "Senior Software Engineer" || updated.Grade != "G3" || updated.Language != "hi" ||
			updated.ContactPhone != "" || updated.JoiningDate != nil || updated.StartDayTime != nil {
			t.Errorf("Update() = %+v; fields not updated", updated)
		}

		profile, err := userProfileService.GetByUserID(ctx, testUserID)
		if err != nil {
			t.Fatalf("GetByUserID() after update error = %v", err)
		}
		if profile.Designation != "Senior Software Engineer" {
			t.Errorf("GetByUserID() designation = %q; want updated value", profile.Designation)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := userProfileService.Delete(ctx, testUserID); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if _, err := userProfileService.GetByUserID(ctx, testUserID); !errors.Is(err, response.UserProfileNotFound) {
			t.Errorf("GetByUserID() after delete error = %v; want %v", err, response.UserProfileNotFound)
		}
		if _, err := userProfileService.Update(ctx, model.UserProfile{UserID: testUserID}); !errors.Is(err, response.UserProfileNotFound) {
			t.Errorf("Update() after delete error = %v; want %v", err, response.UserProfileNotFound)
		}
		if err := userProfileService.Delete(ctx, testUserID); !errors.Is(err, response.UserProfileNotFound) {
			t.Errorf("second Delete() error = %v; want %v", err, response.UserProfileNotFound)
		}
	})
}

func sameDate(got, want time.Time) bool {
	return got.Year() == want.Year() && got.Month() == want.Month() && got.Day() == want.Day()
}

func sameClock(got, want time.Time) bool {
	return got.Hour() == want.Hour() && got.Minute() == want.Minute() && got.Second() == want.Second()
}
