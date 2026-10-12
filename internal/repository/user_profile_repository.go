package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// UserProfileRepository runs the SQL queries for the user_profiles table. Soft-deleted
// profiles (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.UserProfile, so nullable
// columns are read with COALESCE and given their column name with AS.
type UserProfileRepository struct {
	dbPool *pgxpool.Pool
}

// NewUserProfileRepository creates a UserProfileRepository that uses dbPool.
func NewUserProfileRepository(dbPool *pgxpool.Pool) *UserProfileRepository {
	return &UserProfileRepository{dbPool: dbPool}
}

func (userProfileRepository *UserProfileRepository) FindByUserID(ctx context.Context, userID int64) (model.UserProfile, error) {
	query := `SELECT id, user_id,
			COALESCE(designation, '') AS designation,
			COALESCE(grade, '') AS grade,
			birth_date, joining_date,
			COALESCE(blood_group, '') AS blood_group,
			COALESCE(profile_picture, '') AS profile_picture,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			COALESCE(theme, '') AS theme,
			COALESCE(timezone, '') AS timezone,
			COALESCE(language, '') AS language,
			start_day_time,
			is_deleted, deleted_at, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userProfileRepository.dbPool, query, userID)
	if err != nil {
		return model.UserProfile{}, fmt.Errorf("find user profile: %w", err)
	}
	profile, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserProfile])
	return profile, mapProfileError("find user profile", err)
}

func (userProfileRepository *UserProfileRepository) Create(ctx context.Context, profile model.UserProfile) (model.UserProfile, error) {
	query := `INSERT INTO user_profiles (user_id, designation, grade, birth_date, joining_date, blood_group, profile_picture, contact_person, contact_phone, theme, timezone, language, start_day_time)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), $4, $5, NULLIF($6, ''), NULLIF($7, ''), NULLIF($8, ''), NULLIF($9, ''), NULLIF($10, ''), NULLIF($11, ''), NULLIF($12, ''), $13)
		RETURNING id, user_id,
			COALESCE(designation, '') AS designation,
			COALESCE(grade, '') AS grade,
			birth_date, joining_date,
			COALESCE(blood_group, '') AS blood_group,
			COALESCE(profile_picture, '') AS profile_picture,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			COALESCE(theme, '') AS theme,
			COALESCE(timezone, '') AS timezone,
			COALESCE(language, '') AS language,
			start_day_time,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, userProfileRepository.dbPool, query,
		profile.UserID, profile.Designation, profile.Grade, profile.BirthDate, profile.JoiningDate, profile.BloodGroup, profile.ProfilePicture, profile.ContactPerson, profile.ContactPhone, profile.Theme, profile.Timezone, profile.Language, profile.StartDayTime,
	)
	if err != nil {
		return model.UserProfile{}, fmt.Errorf("create user profile: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserProfile])
	return created, mapProfileError("create user profile", err)
}

func (userProfileRepository *UserProfileRepository) Update(ctx context.Context, profile model.UserProfile) (model.UserProfile, error) {
	query := `UPDATE user_profiles
		SET designation = NULLIF($2, ''), grade = NULLIF($3, ''), birth_date = $4, joining_date = $5,
			blood_group = NULLIF($6, ''), profile_picture = NULLIF($7, ''), contact_person = NULLIF($8, ''),
			contact_phone = NULLIF($9, ''), theme = NULLIF($10, ''), timezone = NULLIF($11, ''),
			language = NULLIF($12, ''), start_day_time = $13
		WHERE user_id = $1 AND is_deleted = false
		RETURNING id, user_id,
			COALESCE(designation, '') AS designation,
			COALESCE(grade, '') AS grade,
			birth_date, joining_date,
			COALESCE(blood_group, '') AS blood_group,
			COALESCE(profile_picture, '') AS profile_picture,
			COALESCE(contact_person, '') AS contact_person,
			COALESCE(contact_phone, '') AS contact_phone,
			COALESCE(theme, '') AS theme,
			COALESCE(timezone, '') AS timezone,
			COALESCE(language, '') AS language,
			start_day_time,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, userProfileRepository.dbPool, query,
		profile.UserID, profile.Designation, profile.Grade, profile.BirthDate, profile.JoiningDate, profile.BloodGroup, profile.ProfilePicture, profile.ContactPerson, profile.ContactPhone, profile.Theme, profile.Timezone, profile.Language, profile.StartDayTime,
	)
	if err != nil {
		return model.UserProfile{}, fmt.Errorf("update user profile: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.UserProfile])
	return updated, mapProfileError("update user profile", err)
}

// Delete soft-deletes the user profile: the row is kept with is_deleted = true.
func (userProfileRepository *UserProfileRepository) Delete(ctx context.Context, userID int64) error {
	query := `UPDATE user_profiles
		SET is_deleted = true, deleted_at = NOW()
		WHERE user_id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, userProfileRepository.dbPool, query, userID)
	if err != nil {
		return fmt.Errorf("delete user profile: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete user profile: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapProfileError turns driver errors into repository errors. A nil err stays nil.
func mapProfileError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
