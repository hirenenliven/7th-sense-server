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

// CompanyProfileRepository runs the SQL queries for the company_profiles table. Soft-deleted
// profiles (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.CompanyProfile, so nullable
// columns are read with COALESCE and given their column name with AS.
type CompanyProfileRepository struct {
	dbPool *pgxpool.Pool
}

// NewCompanyProfileRepository creates a CompanyProfileRepository that uses dbPool.
func NewCompanyProfileRepository(dbPool *pgxpool.Pool) *CompanyProfileRepository {
	return &CompanyProfileRepository{dbPool: dbPool}
}

func (companyProfileRepository *CompanyProfileRepository) FindByUserID(ctx context.Context, userID int64) (model.CompanyProfile, error) {
	query := `SELECT id, user_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at
		FROM company_profiles
		WHERE user_id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, companyProfileRepository.dbPool, query, userID)
	if err != nil {
		return model.CompanyProfile{}, fmt.Errorf("find company profile: %w", err)
	}
	profile, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CompanyProfile])
	return profile, mapCompanyProfileError("find company profile", err)
}

func (companyProfileRepository *CompanyProfileRepository) Create(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	query := `INSERT INTO company_profiles (user_id, gst_number, pan_number, website, finix_id)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), $5)
		RETURNING id, user_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, companyProfileRepository.dbPool, query,
		profile.UserID, profile.GSTNumber, profile.PANNumber, profile.Website, profile.FinixID,
	)
	if err != nil {
		return model.CompanyProfile{}, fmt.Errorf("create company profile: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CompanyProfile])
	return created, mapCompanyProfileError("create company profile", err)
}

func (companyProfileRepository *CompanyProfileRepository) Update(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	query := `UPDATE company_profiles
		SET gst_number = NULLIF($2, ''), pan_number = NULLIF($3, ''), website = NULLIF($4, ''), finix_id = $5
		WHERE user_id = $1 AND is_deleted = false
		RETURNING id, user_id,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, companyProfileRepository.dbPool, query,
		profile.UserID, profile.GSTNumber, profile.PANNumber, profile.Website, profile.FinixID,
	)
	if err != nil {
		return model.CompanyProfile{}, fmt.Errorf("update company profile: %w", err)
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CompanyProfile])
	return updated, mapCompanyProfileError("update company profile", err)
}

// Delete soft-deletes the company profile: the row is kept with is_deleted = true.
func (companyProfileRepository *CompanyProfileRepository) Delete(ctx context.Context, userID int64) error {
	query := `UPDATE company_profiles
		SET is_deleted = true, deleted_at = NOW()
		WHERE user_id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, companyProfileRepository.dbPool, query, userID)
	if err != nil {
		return fmt.Errorf("delete company profile: %w", err)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("delete company profile: %w", err)
	}
	if rows.CommandTag().RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// mapCompanyProfileError turns driver errors into repository errors. A nil err stays nil.
func mapCompanyProfileError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
