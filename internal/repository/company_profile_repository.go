package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// CompanyProfileRepository runs the SQL queries for the company_profiles table. Soft-deleted
// profiles (is_deleted = true) are invisible to every method.
//
// Selected column names must match the db tags on model.CompanyProfile, so nullable
// columns are read with COALESCE and given their column name with AS.
type CompanyProfileRepository struct {
	database db.Querier
}

// NewCompanyProfileRepository creates a CompanyProfileRepository that runs its queries on database
// (a connection pool or a transaction).
func NewCompanyProfileRepository(database db.Querier) *CompanyProfileRepository {
	return &CompanyProfileRepository{database: database}
}

func (companyProfileRepository *CompanyProfileRepository) FindByUserID(ctx context.Context, userID int64) (model.CompanyProfile, error) {
	query := `SELECT id, user_id,
			COALESCE(company_name, '') AS company_name,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at
		FROM company_profiles
		WHERE user_id = $1 AND is_deleted = false`

	rows, err := db.Execute(ctx, companyProfileRepository.database, query, userID)
	if err != nil {
		return model.CompanyProfile{}, fmt.Errorf("find company profile: %w", err)
	}
	profile, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CompanyProfile])
	return profile, mapCompanyProfileError("find company profile", err)
}

func (companyProfileRepository *CompanyProfileRepository) Create(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	query := `INSERT INTO company_profiles (user_id, company_name, gst_number, pan_number, website, finix_id)
		VALUES ($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), $6)
		RETURNING id, user_id,
			COALESCE(company_name, '') AS company_name,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, companyProfileRepository.database, query,
		profile.UserID, profile.CompanyName, profile.GSTNumber, profile.PANNumber, profile.Website, profile.FinixID,
	)
	if err != nil {
		return model.CompanyProfile{}, fmt.Errorf("create company profile: %w", err)
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CompanyProfile])
	return created, mapCompanyProfileError("create company profile", err)
}

func (companyProfileRepository *CompanyProfileRepository) Update(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	query := `UPDATE company_profiles
		SET company_name = NULLIF($2, ''), gst_number = NULLIF($3, ''), pan_number = NULLIF($4, ''),
			website = NULLIF($5, ''), finix_id = $6
		WHERE user_id = $1 AND is_deleted = false
		RETURNING id, user_id,
			COALESCE(company_name, '') AS company_name,
			COALESCE(gst_number, '') AS gst_number,
			COALESCE(pan_number, '') AS pan_number,
			COALESCE(website, '') AS website,
			finix_id,
			is_deleted, deleted_at, created_at, updated_at`

	rows, err := db.Execute(ctx, companyProfileRepository.database, query,
		profile.UserID, profile.CompanyName, profile.GSTNumber, profile.PANNumber, profile.Website, profile.FinixID,
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

	rows, err := db.Execute(ctx, companyProfileRepository.database, query, userID)
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
