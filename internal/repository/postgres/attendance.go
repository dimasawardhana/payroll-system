package postgres

import (
	"context"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepository struct {
	pool *pgxpool.Pool
}

func NewAttendanceRepository(pool *pgxpool.Pool) *AttendanceRepository {
	return &AttendanceRepository{
		pool: pool,
	}
}

func (r *AttendanceRepository) RecordAttendance(ctx context.Context, payload domain.Attendance) error {
	// Implement logic to record employee attendance
	// This is a placeholder implementation
	if payload.EmployeeID == 0 {
		return error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO attendance (employee_id, date, status, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, payload.EmployeeID, payload.Date, payload.Status, payload.CreatedAt, payload.UpdatedAt, payload.CreatedBy, payload.UpdatedBy)

	if err != nil {
		return err
	}
	return nil
}

func (r *AttendanceRepository) GetAttendanceByEmployeeID(ctx context.Context, employeeID int64) ([]domain.Attendance, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, date, status, created_at, updated_at, created_by, updated_by
		FROM attendance WHERE employee_id = $1
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(&attendance.ID, &attendance.EmployeeID, &attendance.Date, &attendance.Status,
			&attendance.CreatedAt, &attendance.UpdatedAt, &attendance.CreatedBy, &attendance.UpdatedBy)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}
func (r *AttendanceRepository) GetEmployeeAttendanceByDateRange(ctx context.Context, employeeID int64, startDate, endDate time.Time) ([]domain.Attendance, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, date, status, created_at, updated_at, created_by, updated_by
		FROM attendance WHERE employee_id = $1 AND date BETWEEN $2 AND $3
	`, employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(&attendance.ID, &attendance.EmployeeID, &attendance.Date, &attendance.Status,
			&attendance.CreatedAt, &attendance.UpdatedAt, &attendance.CreatedBy, &attendance.UpdatedBy)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}
func (r *AttendanceRepository) GetAllAttendanceByDateRange(ctx context.Context, startDate, endDate time.Time) ([]domain.Attendance, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, date, status, created_at, updated_at, created_by, updated_by
		FROM attendance WHERE date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []domain.Attendance
	for rows.Next() {
		var attendance domain.Attendance
		err := rows.Scan(&attendance.ID, &attendance.EmployeeID, &attendance.Date, &attendance.Status,
			&attendance.CreatedAt, &attendance.UpdatedAt, &attendance.CreatedBy, &attendance.UpdatedBy)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}

func (r *AttendanceRepository) GetTotalAttendanceByDateRangeGroupedByEmployee(ctx context.Context, startDate, endDate time.Time) (map[int]int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT employee_id, COUNT(*) as total_attendance
		FROM attendance
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
		var total int
		if err := rows.Scan(&employeeID, &total); err != nil {
			return nil, err
		}
		result[employeeID] = total
	}
	return result, nil
}
