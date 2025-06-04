package postgres

import (
	"context"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OvertimeRepository struct {
	pool *pgxpool.Pool
}

func NewOvertimeRepository(pool *pgxpool.Pool) *OvertimeRepository {
	return &OvertimeRepository{
		pool: pool,
	}
}

func (r *OvertimeRepository) SubmitOvertime(ctx context.Context, payload domain.Overtime) error {

	if payload.EmployeeID == 0 || payload.Hours <= 0 {
		return error_const.ErrInvalidUser // Return an error if employee ID or hours are invalid
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO overtime (employee_id, hours, date, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, payload.EmployeeID, payload.Hours, payload.Date, payload.CreatedAt, payload.UpdatedAt, payload.CreatedBy, payload.UpdatedBy)
	if err != nil {
		return err
	}
	return nil
}

func (r *OvertimeRepository) GetOvertimesByEmployeeID(ctx context.Context, employeeID int64) ([]domain.Overtime, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, hours, date, created_at, updated_at, created_by, updated_by
		FROM overtime WHERE employee_id = $1
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overtimes []domain.Overtime
	for rows.Next() {
		var overtime domain.Overtime
		err := rows.Scan(&overtime.ID, &overtime.EmployeeID, &overtime.Hours,
			&overtime.Date, &overtime.CreatedAt, &overtime.UpdatedAt,
			&overtime.CreatedBy, &overtime.UpdatedBy)
		if err != nil {
			return nil, err
		}
		overtimes = append(overtimes, overtime)
	}
	return overtimes, nil
}
func (r *OvertimeRepository) GetEmployeeOvertimeByDateRange(ctx context.Context, employeeID int64, startDate, endDate time.Time) ([]domain.Overtime, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, hours, date, created_at, updated_at, created_by, updated_by
		FROM overtime WHERE employee_id = $1 AND date BETWEEN $2 AND $3
	`, employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overtimes []domain.Overtime
	for rows.Next() {
		var overtime domain.Overtime
		err := rows.Scan(&overtime.ID, &overtime.EmployeeID, &overtime.Hours,
			&overtime.Date, &overtime.CreatedAt, &overtime.UpdatedAt,
			&overtime.CreatedBy, &overtime.UpdatedBy)
		if err != nil {
			return nil, err
		}
		overtimes = append(overtimes, overtime)
	}
	return overtimes, nil
}
func (r *OvertimeRepository) GetOvertimesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]domain.Overtime, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, hours, date, created_at, updated_at, created_by, updated_by
		FROM overtime WHERE date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overtimes []domain.Overtime
	for rows.Next() {
		var overtime domain.Overtime
		err := rows.Scan(&overtime.ID, &overtime.EmployeeID, &overtime.Hours,
			&overtime.Date, &overtime.CreatedAt, &overtime.UpdatedAt,
			&overtime.CreatedBy, &overtime.UpdatedBy)
		if err != nil {
			return nil, err
		}
		overtimes = append(overtimes, overtime)
	}
	return overtimes, nil
}

func (r *OvertimeRepository) GetTotalOvertimeHoursByDateRange(ctx context.Context, startDate, endDate time.Time) (map[int]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT employee_id, SUM(hours) as total_hours
		FROM overtime
		WHERE date BETWEEN $1 AND $2
		GROUP BY employee_id
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var employeeID int
		var totalHours int
		if err := rows.Scan(&employeeID, &totalHours); err != nil {
			return nil, err
		}
		result[employeeID] = totalHours
	}
	return result, nil
}

func (r *OvertimeRepository) GetOvertimesGroupedByEmployeeID(ctx context.Context, startDate, endDate time.Time) (map[int][]domain.Overtime, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, hours, date, created_at, updated_at, created_by, updated_by
		FROM overtime
		WHERE date BETWEEN $1 AND $2
		ORDER BY employee_id, date
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]domain.Overtime)
	for rows.Next() {
		var overtime domain.Overtime
		err := rows.Scan(&overtime.ID, &overtime.EmployeeID, &overtime.Hours,
			&overtime.Date, &overtime.CreatedAt, &overtime.UpdatedAt,
			&overtime.CreatedBy, &overtime.UpdatedBy)
		if err != nil {
			return nil, err
		}
		result[overtime.EmployeeID] = append(result[overtime.EmployeeID], overtime)
	}
	return result, nil
}
