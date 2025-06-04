package postgres

import (
	"context"
	"encoding/json"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PayrollRepository struct {
	pool *pgxpool.Pool
}

func NewPayrollRepository(pool *pgxpool.Pool) *PayrollRepository {
	return &PayrollRepository{
		pool: pool,
	}
}

func (r *PayrollRepository) CreatePayroll(ctx context.Context, payroll domain.Payroll) error {
	if payroll.EmployeeID == 0 || payroll.PeriodID == 0 {
		return error_const.ErrInvalidUser
	}
	if payroll.CreatedBy == "" || payroll.UpdatedBy == "" {
		return error_const.ErrInvalidUser
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO payrolls (employee_id, period_id, payslip, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, NOW(), NOW(), $4, $5)
		RETURNING id`, payroll.EmployeeID, payroll.PeriodID, payroll.Payslip, payroll.CreatedBy, payroll.UpdatedBy)

	if err != nil {
		return err
	}
	return nil
}

// bulkInsert payrolls
func (r *PayrollRepository) BulkInsertPayrolls(ctx context.Context, payrolls []domain.Payroll) error {
	if len(payrolls) == 0 {
		return error_const.ErrInvalidInput
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare data for COPY FROM
	rows := make([][]interface{}, 0, len(payrolls))
	for _, payroll := range payrolls {
		if payroll.EmployeeID == 0 || payroll.PeriodID == 0 {
			return error_const.ErrInvalidUser
		}
		if payroll.CreatedBy == "" || payroll.UpdatedBy == "" {
			return error_const.ErrInvalidUser
		}
		rows = append(rows, []interface{}{
			payroll.EmployeeID,
			payroll.PeriodID,
			payroll.Payslip,
			payroll.CreatedBy,
			payroll.UpdatedBy,
		})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"payrolls"},
		[]string{"employee_id", "period_id", "payslip", "created_by", "updated_by"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PayrollRepository) GetPayrollPeriod(ctx context.Context, period domain.PayrollPeriod) (domain.PayrollPeriod, error) {
	var result domain.PayrollPeriod

	if period.ID == 0 {
		return domain.PayrollPeriod{}, error_const.ErrInvalidID
	}
	query := `
		SELECT id, start_date, end_date, locked, created_at, updated_at, created_by, updated_by
		FROM payroll_periods
		WHERE id = $1
		LIMIT 1
	`
	row := r.pool.QueryRow(ctx, query, period.ID)
	err := row.Scan(
		&result.ID,
		&result.StartDate,
		&result.EndDate,
		&result.Locked,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.CreatedBy,
		&result.UpdatedBy,
	)
	if err != nil {
		return domain.PayrollPeriod{}, err
	}
	return result, nil
}
func (r *PayrollRepository) GetPayrollPeriodFromDateRange(ctx context.Context, period domain.PayrollPeriod) (domain.PayrollPeriod, error) {
	var result domain.PayrollPeriod

	if period.StartDate.IsZero() || period.EndDate.IsZero() {
		return domain.PayrollPeriod{}, error_const.ErrInvalidDateFormat
	}
	query := `
		SELECT id, start_date, end_date, locked, created_at, updated_at, created_by, updated_by
		FROM payroll_periods
		WHERE start_date <= $1 AND end_date >= $2
		LIMIT 1
	`
	row := r.pool.QueryRow(ctx, query, period.StartDate, period.EndDate)
	err := row.Scan(
		&result.ID,
		&result.StartDate,
		&result.EndDate,
		&result.Locked,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.CreatedBy,
		&result.UpdatedBy,
	)
	if err != nil {
		return domain.PayrollPeriod{}, err
	}
	return result, nil
}
func (r *PayrollRepository) GetPayrollPeriodFromDate(ctx context.Context, date domain.PayrollPeriodDate) (domain.PayrollPeriod, error) {
	var result domain.PayrollPeriod

	if date.IsZero() {
		return domain.PayrollPeriod{}, error_const.ErrInvalidDateFormat
	}
	query := `
		SELECT id, start_date, end_date, locked, created_at, updated_at, created_by, updated_by
		FROM payroll_periods
		WHERE start_date <= $1 AND end_date >= $1
		LIMIT 1
	`
	row := r.pool.QueryRow(ctx, query, date)
	err := row.Scan(
		&result.ID,
		&result.StartDate,
		&result.EndDate,
		&result.Locked,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.CreatedBy,
		&result.UpdatedBy,
	)
	if err != nil {
		return domain.PayrollPeriod{}, err
	}
	return result, nil
}

func (r *PayrollRepository) CreatePayrollPeriod(ctx context.Context, payroll domain.PayrollPeriod) (string, error) {
	var payrollID string
	if payroll.StartDate.IsZero() || payroll.EndDate.IsZero() {
		return "", error_const.ErrInvalidDateFormat
	}
	if payroll.CreatedBy == "" || payroll.UpdatedBy == "" {
		return "", error_const.ErrInvalidUser
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO payroll_periods (start_date, end_date, locked, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, false, NOW(), NOW(), $3, $4)
		RETURNING id`, payroll.StartDate, payroll.EndDate, payroll.CreatedBy, payroll.UpdatedBy).Scan(&payrollID)

	if err != nil {
		return "", err
	}
	return payrollID, nil
}

func (r *PayrollRepository) GetEmployeePayslipByPeriod(ctx context.Context, _payroll domain.Payroll) (domain.Payroll, error) {
	if _payroll.EmployeeID == 0 || _payroll.PeriodID == 0 {
		return domain.Payroll{}, error_const.ErrInvalidUser
	}

	var payroll domain.Payroll
	query := `
		SELECT id, employee_id, period_id, payslip, created_at, updated_at, created_by, updated_by
		FROM payrolls
		WHERE employee_id = $1 AND period_id = $2
		LIMIT 1
	`
	row := r.pool.QueryRow(ctx, query, _payroll.EmployeeID, _payroll.EmployeeID)
	err := row.Scan(
		&payroll.ID,
		&payroll.EmployeeID,
		&payroll.PeriodID,
		&payroll.Payslip,
		&payroll.CreatedAt,
		&payroll.UpdatedAt,
		&payroll.CreatedBy,
		&payroll.UpdatedBy,
	)
	if err != nil {
		return domain.Payroll{}, err
	}
	return payroll, nil
}

func (r *PayrollRepository) GetPayrollsByPeriodID(ctx context.Context, periodID int) ([]domain.Payroll, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, employee_id, period_id, payslip, created_by, updated_by
		FROM payrolls
		WHERE period_id = $1
	`, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payrolls []domain.Payroll
	for rows.Next() {
		var p domain.Payroll
		var payslipData []byte
		if err := rows.Scan(&p.ID, &p.EmployeeID, &p.PeriodID, &payslipData, &p.CreatedBy, &p.UpdatedBy); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payslipData, &p.Payslip); err != nil {
			return nil, err
		}
		payrolls = append(payrolls, p)
	}
	return payrolls, nil
}

func (r *PayrollRepository) LockPayrollPeriod(ctx context.Context, periodID int) error {
	if periodID == 0 {
		return error_const.ErrInvalidID
	}

	_, err := r.pool.Exec(ctx, `
		UPDATE payroll_periods
		SET locked = true, updated_at = NOW()
		WHERE id = $1
	`, periodID)
	if err != nil {
		return err
	}
	return nil
}
