package postgres

import (
	"context"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReimbursementRepository struct {
	pool *pgxpool.Pool
}

func NewReimbursementRepository(pool *pgxpool.Pool) *ReimbursementRepository {
	return &ReimbursementRepository{
		pool: pool,
	}
}

func (r *ReimbursementRepository) SubmitReimbursement(ctx context.Context, payload domain.Reimbursement) error {
	if payload.EmployeeID == 0 || payload.Amount <= 0 {
		return error_const.ErrInvalidUser // Return an error if employee ID or amount is invalid
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO reimbursements (employee_id, amount, description, date, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, payload.EmployeeID, payload.Amount, payload.Description, payload.Date, payload.CreatedAt, payload.UpdatedAt, payload.CreatedBy, payload.UpdatedBy)

	if err != nil {
		return err
	}
	return nil
}

func (r *ReimbursementRepository) GetReimbursementsByEmployeeID(ctx context.Context, employeeID int64) ([]domain.Reimbursement, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, amount, description, date, created_at, updated_at, created_by, updated_by
		FROM reimbursements WHERE employee_id = $1
	`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reimbursements []domain.Reimbursement
	for rows.Next() {
		var reimbursement domain.Reimbursement
		err := rows.Scan(&reimbursement.ID, &reimbursement.EmployeeID, &reimbursement.Amount,
			&reimbursement.Description, &reimbursement.Date, &reimbursement.CreatedAt, &reimbursement.UpdatedAt,
			&reimbursement.CreatedBy, &reimbursement.UpdatedBy)
		if err != nil {
			return nil, err
		}
		reimbursements = append(reimbursements, reimbursement)
	}
	return reimbursements, nil
}

func (r *ReimbursementRepository) GetAllReimbursementsFromDateRange(ctx context.Context, startDate, endDate string) ([]domain.Reimbursement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, amount, description, date, created_at, updated_at, created_by, updated_by
		FROM reimbursements WHERE date BETWEEN $1 AND $2
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reimbursements []domain.Reimbursement
	for rows.Next() {
		var reimbursement domain.Reimbursement
		err := rows.Scan(&reimbursement.ID, &reimbursement.EmployeeID, &reimbursement.Amount,
			&reimbursement.Description, &reimbursement.Date, &reimbursement.CreatedAt, &reimbursement.UpdatedAt,
			&reimbursement.CreatedBy, &reimbursement.UpdatedBy)
		if err != nil {
			return nil, err
		}
		reimbursements = append(reimbursements, reimbursement)
	}
	return reimbursements, nil
}

func (r *ReimbursementRepository) GetEmployeeReimbursementByDateRange(ctx context.Context, employeeID int64, startDate, endDate time.Time) ([]domain.Reimbursement, error) {
	if employeeID == 0 {
		return nil, error_const.ErrInvalidUser // Return an error if employee ID is invalid
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, amount, description, date, created_at, updated_at, created_by, updated_by
		FROM reimbursements WHERE employee_id = $1 AND date BETWEEN $2 AND $3
	`, employeeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reimbursements []domain.Reimbursement
	for rows.Next() {
		var reimbursement domain.Reimbursement
		err := rows.Scan(&reimbursement.ID, &reimbursement.EmployeeID, &reimbursement.Amount,
			&reimbursement.Description, &reimbursement.Date, &reimbursement.CreatedAt, &reimbursement.UpdatedAt,
			&reimbursement.CreatedBy, &reimbursement.UpdatedBy)
		if err != nil {
			return nil, err
		}
		reimbursements = append(reimbursements, reimbursement)
	}
	return reimbursements, nil
}

func (r *ReimbursementRepository) GetReimbursementsGroupedByEmployeeID(ctx context.Context, startDate, endDate time.Time) (map[int][]domain.Reimbursement, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, amount, description, date, created_at, updated_at, created_by, updated_by
		FROM reimbursements
		WHERE date BETWEEN $1 AND $2
		ORDER BY employee_id, date
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reimbursementsByEmployee := make(map[int][]domain.Reimbursement)
	for rows.Next() {
		var reimbursement domain.Reimbursement
		err := rows.Scan(&reimbursement.ID, &reimbursement.EmployeeID, &reimbursement.Amount,
			&reimbursement.Description, &reimbursement.Date, &reimbursement.CreatedAt, &reimbursement.UpdatedAt,
			&reimbursement.CreatedBy, &reimbursement.UpdatedBy)
		if err != nil {
			return nil, err
		}
		reimbursementsByEmployee[reimbursement.EmployeeID] = append(reimbursementsByEmployee[reimbursement.EmployeeID], reimbursement)
	}
	return reimbursementsByEmployee, nil
}
