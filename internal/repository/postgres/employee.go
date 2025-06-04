package postgres

import (
	"context"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		pool: pool,
	}
}
func (r *EmployeeRepository) GetEmployee(ctx context.Context, credential domain.Employee) (domain.Employee, error) {
	var employee domain.Employee
	// This is a placeholder implementation
	if credential.Email == "" {
		return domain.Employee{}, nil // Return an error or false if credentials are invalid
	}

	err := r.pool.
		QueryRow(ctx, "SELECT id, name, email, password_hash, role, salary FROM employees WHERE email = $1", credential.Email).
		Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Password_hash, &employee.Role, &employee.Salary)

	if err != nil {
		return domain.Employee{}, err
	}
	if employee.ID == 0 {
		return domain.Employee{}, error_const.ErrUserNotFound
	}
	return employee, nil
}
func (r *EmployeeRepository) GetAllEmployees(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, email, password_hash, role, salary FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.Password_hash, &employee.Role, &employee.Salary)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (r *EmployeeRepository) GetPayslip(employeeID string, period string) (interface{}, error) {
	// Implement logic to retrieve employee payslip
	// This is a placeholder implementation
	if employeeID == "" || period == "" {
		return nil, nil // Return an error or empty payslip if parameters are invalid
	}
	return nil, nil
}

func (r *EmployeeRepository) GetEmployeeByID(ctx context.Context, employeeID int) (*domain.Employee, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, salary
		FROM employees
		WHERE id = $1
	`, employeeID)

	var e domain.Employee
	if err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Password_hash, &e.Role, &e.Salary); err != nil {
		return nil, err
	}
	return &e, nil
}
