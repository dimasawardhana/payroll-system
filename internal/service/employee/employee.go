package employee_service

import (
	"context"
	"errors"
	"payroll-system/internal/delivery/dto"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"payroll-system/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EmployeeRepository interface {
	GetEmployee(ctx context.Context, credential domain.Employee) (domain.Employee, error)
}
type PayrollRepository interface {
	GetPayrollPeriodFromDate(ctx context.Context, date time.Time) (domain.PayrollPeriod, error)
	GetEmployeePayslipByPeriod(ctx context.Context, payroll domain.Payroll) (domain.Payroll, error)
}
type AttendanceRepository interface {
	RecordAttendance(ctx context.Context, attendance domain.Attendance) error
}
type OvertimeRepository interface {
	SubmitOvertime(ctx context.Context, overtime domain.Overtime) error
}
type ReimbursementRepository interface {
	SubmitReimbursement(ctx context.Context, reimbursement domain.Reimbursement) error
}

type EmployeeService struct {
	empRepo           EmployeeRepository
	payrollRepo       PayrollRepository
	attendanceRepo    AttendanceRepository
	overtimeRepo      OvertimeRepository
	reimbursementRepo ReimbursementRepository
}

func NewEmployeeService(empRepo EmployeeRepository, payrollRepo PayrollRepository,
	attendanceRepo AttendanceRepository, overtimeRepo OvertimeRepository,
	reimbursementRepo ReimbursementRepository) *EmployeeService {
	return &EmployeeService{
		empRepo:           empRepo,
		payrollRepo:       payrollRepo,
		attendanceRepo:    attendanceRepo,
		overtimeRepo:      overtimeRepo,
		reimbursementRepo: reimbursementRepo,
	}
}

func (s *EmployeeService) LoginAsEmployee(ctx context.Context, credentials dto.LoginRequest) (string, error) {
	admin, err := s.empRepo.GetEmployee(ctx, domain.Employee{
		Email: credentials.Email,
	})
	if err != nil {
		return "", err
	}

	validPassword := utils.CheckPassword(admin.Password_hash, credentials.Password)
	if !validPassword {
		return "", error_const.ErrInvalidCredentials
	}
	token, err := utils.GenerateJWT(admin.ID, admin.Email, admin.Role)
	if err != nil {
		return "", err

	}
	return token, nil
}
func (s *EmployeeService) RecordAttendance(ctx context.Context, payload dto.AttendanceRequest) error {
	var attendance domain.Attendance

	if payload.EmployeeID == 0 {
		return error_const.ErrInvalidCredentials
	}

	attendanceDate, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return error_const.ErrInvalidDateFormat
	}
	if attendanceDate.Weekday() == time.Saturday || attendanceDate.Weekday() == time.Sunday {
		return error_const.ErrAttendanceOnWeekend
	}

	attendance.Date = attendanceDate
	attendance.EmployeeID = payload.EmployeeID
	currentTime := time.Now()
	attendance.CreatedAt = currentTime
	attendance.UpdatedAt = currentTime
	attendance.Status = "present"
	attendance.CreatedBy = payload.EmployeeEmail
	attendance.UpdatedBy = payload.EmployeeEmail

	// making sure the attendance date is within the payroll period
	payrollPeriod, err := s.payrollRepo.GetPayrollPeriodFromDate(ctx, attendance.Date)
	if err != nil {
		return error_const.ErrPayrollPeriodNotFound
	}
	if payrollPeriod.Locked {
		return error_const.ErrPayrollPeriodLocked
	}

	if err := s.attendanceRepo.RecordAttendance(ctx, attendance); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return error_const.ErrAttendanceAlreadyExists
		}
		return err
	}
	return nil
}

func (s *EmployeeService) SubmitOvertime(ctx context.Context, payload dto.OvertimeRequest) error {

	var overtime domain.Overtime
	if payload.EmployeeID == 0 {
		return error_const.ErrInvalidCredentials
	}
	if payload.Hours <= 0 {
		return error_const.ErrInvalidOvertimeHours
	}
	if payload.Hours > 3 {
		return error_const.ErrOvertimeHoursExceeded
	}
	overtime.EmployeeID = payload.EmployeeID
	overtime.Hours = payload.Hours
	currentTime := time.Now()
	overtime.Date = currentTime
	overtime.CreatedAt = currentTime
	overtime.UpdatedAt = currentTime
	overtime.CreatedBy = payload.EmployeeEmail
	overtime.UpdatedBy = payload.EmployeeEmail
	if err := s.overtimeRepo.SubmitOvertime(ctx, overtime); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return error_const.ErrOvertimeAlreadyExists
		}
	}
	return nil
}

func (s *EmployeeService) SubmitReimbursement(ctx context.Context, payload dto.ReimbursementRequest) error {
	var reimbursement domain.Reimbursement
	if payload.EmployeeID == 0 {
		return error_const.ErrInvalidCredentials
	}
	if payload.Amount <= 0 {
		return error_const.ErrInvalidReimbursementAmount
	}
	reimbursement.EmployeeID = payload.EmployeeID
	reimbursement.Amount = payload.Amount
	reimbursement.Description = payload.Description
	currentTime := time.Now()
	reimbursement.CreatedAt = currentTime
	reimbursement.UpdatedAt = currentTime
	reimbursement.CreatedBy = payload.EmployeeEmail
	reimbursement.UpdatedBy = payload.EmployeeEmail
	if err := s.reimbursementRepo.SubmitReimbursement(ctx, reimbursement); err != nil {
		return err
	}
	return nil
}

func (s *EmployeeService) GetPayslip(ctx context.Context, payload dto.PayrollRequest) (interface{}, error) {
	payroll, err := s.payrollRepo.GetEmployeePayslipByPeriod(ctx, domain.Payroll{
		EmployeeID: payload.EmployeeID,
		PeriodID:   payload.PeriodID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, error_const.ErrPayslipNotFound
		}
		return nil, err
	}
	return payroll, nil
}
