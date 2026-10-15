package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
)

// holidaySelect reads model.Holiday: the holidays row (h) plus its customer's display
// name. DATE columns are read as "YYYY-MM-DD".
const holidaySelect = `SELECT h.id, h.user_id, h.holiday_type, COALESCE(h.name, '') AS name,
		h.start_date::text AS start_date, h.end_date::text AS end_date, h.week_day,
		h.status, h.is_deleted, h.deleted_at, h.created_at, h.updated_at,
		COALESCE(NULLIF(TRIM(CONCAT_WS(' ', u.first_name, u.last_name)), ''), u.username, '') AS customer_name,
		COALESCE(cp.company_name, '') AS company_name
	FROM holidays h
	LEFT JOIN users u ON u.id = h.user_id
	LEFT JOIN company_profiles cp ON cp.user_id = h.user_id AND cp.is_deleted = false`

// HolidayRepository runs the SQL queries for the holidays table (holidays and week
// offs). Soft-deleted rows (is_deleted = true) are invisible to every method.
//
// A nil ownerID does not narrow a query; otherwise only the rows of that customer
// (holidays.user_id) are found.
type HolidayRepository struct {
	database db.Querier
}

// NewHolidayRepository creates a HolidayRepository that runs its queries on database
// (a connection pool or a transaction).
func NewHolidayRepository(database db.Querier) *HolidayRepository {
	return &HolidayRepository{database: database}
}

// FindAll returns the holidays and week offs of ownerID (every customer's when nil):
// week offs first by weekday, then holidays by date.
func (holidayRepository *HolidayRepository) FindAll(ctx context.Context, ownerID *int64) ([]model.Holiday, error) {
	query := holidaySelect + `
		WHERE h.is_deleted = false AND ($1::integer IS NULL OR h.user_id = $1)
		ORDER BY h.user_id, h.holiday_type DESC, h.week_day, h.start_date, h.id`

	rows, err := db.Execute(ctx, holidayRepository.database, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("find all holidays: %w", err)
	}
	holidays, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Holiday])
	if err != nil {
		return nil, fmt.Errorf("find all holidays: %w", err)
	}
	return holidays, nil
}

// FindByID returns the holiday, or ErrNotFound when it is not one of ownerID's.
func (holidayRepository *HolidayRepository) FindByID(ctx context.Context, ownerID *int64, id int64) (model.Holiday, error) {
	query := holidaySelect + `
		WHERE h.id = $1 AND h.is_deleted = false AND ($2::integer IS NULL OR h.user_id = $2)`

	rows, err := db.Execute(ctx, holidayRepository.database, query, id, ownerID)
	if err != nil {
		return model.Holiday{}, fmt.Errorf("find holiday: %w", err)
	}
	holiday, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Holiday])
	return holiday, mapHolidayError("find holiday", err)
}

// Create inserts holiday and returns the new id. Blank name and dates are stored as NULL.
func (holidayRepository *HolidayRepository) Create(ctx context.Context, holiday model.Holiday) (int64, error) {
	query := `INSERT INTO holidays (user_id, holiday_type, name, start_date, end_date, week_day, status)
		VALUES ($1, $2, NULLIF($3, ''), $4::date, $5::date, $6, $7)
		RETURNING id`

	rows, err := db.Execute(ctx, holidayRepository.database, query,
		holiday.UserID, holiday.HolidayType, holiday.Name, holiday.StartDate, holiday.EndDate, holiday.WeekDay, holiday.Status)
	if err != nil {
		return 0, fmt.Errorf("create holiday: %w", err)
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int64])
	return id, mapHolidayError("create holiday", err)
}

// Update replaces the fields of holiday.ID. A blank Status keeps the saved status.
func (holidayRepository *HolidayRepository) Update(ctx context.Context, holiday model.Holiday) error {
	query := `UPDATE holidays
		SET user_id = $2, holiday_type = $3, name = NULLIF($4, ''), start_date = $5::date, end_date = $6::date,
			week_day = $7, status = COALESCE(NULLIF($8, ''), status)
		WHERE id = $1 AND is_deleted = false`

	affected, err := executeAffected(ctx, holidayRepository.database, query,
		holiday.ID, holiday.UserID, holiday.HolidayType, holiday.Name, holiday.StartDate, holiday.EndDate, holiday.WeekDay, holiday.Status)
	if err != nil {
		return mapHolidayError("update holiday", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete soft-deletes the holiday of ownerID (any customer's when nil).
func (holidayRepository *HolidayRepository) Delete(ctx context.Context, ownerID *int64, id int64) error {
	query := `UPDATE holidays
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1 AND is_deleted = false AND ($2::integer IS NULL OR user_id = $2)`

	affected, err := executeAffected(ctx, holidayRepository.database, query, id, ownerID)
	if err != nil {
		return fmt.Errorf("delete holiday: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

// mapHolidayError turns driver errors into repository errors. A nil err stays nil.
func mapHolidayError(operation string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch {
		case pgError.Code == "23505" && pgError.ConstraintName == "holidays_user_week_day_key": // unique_violation
			return ErrDuplicateWeekOff
		case pgError.Code == "23505" && pgError.ConstraintName == "holidays_user_name_date_key":
			return ErrDuplicateName
		case pgError.Code == "23503" && pgError.ConstraintName == "holidays_user_id_fkey": // foreign_key_violation
			return ErrInvalidUser
		}
	}
	return fmt.Errorf("%s: %w", operation, err)
}
