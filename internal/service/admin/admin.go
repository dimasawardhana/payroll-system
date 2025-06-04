package admin_service

import (
	"context"
	"errors"
	"fmt"
	"payroll-system/internal/delivery/dto"
	domain "payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"payroll-system/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

type AdminRepository interface {
	GetAdmin(ctx context.Context, credential domain.Admin) (domain.Admin, error)
}

type EmployeeRepository interface {
	GetEmployeeByID(ctx context.Context, employeeID int) (*domain.Employee, error)
	GetAllEmployees(ctx context.Context) ([]domain.Employee, error)
}

type PayrollRepository interface {
	BulkInsertPayrolls(ctx context.Context, payrolls []domain.Payroll) error
	GetPayrollsByPeriodID(ctx context.Context, periodID int) ([]domain.Payroll, error)
	GetPayrollPeriod(ctx context.Context, period domain.PayrollPeriod) (domain.PayrollPeriod, error)
	GetPayrollPeriodFromDateRange(ctx context.Context, period domain.PayrollPeriod) (domain.PayrollPeriod, error)
	CreatePayrollPeriod(ctx context.Context, payroll domain.PayrollPeriod) (string, error)
	LockPayrollPeriod(ctx context.Context, periodID int) error
}

type AttendanceRepository interface {
	GetTotalAttendanceByDateRangeGroupedByEmployee(ctx context.Context, startDate, endDate time.Time) (map[int]int, error)
}
type OvertimeRepository interface {
	GetOvertimesGroupedByEmployeeID(ctx context.Context, startDate, endDate time.Time) (map[int][]domain.Overtime, error)
}
type ReimbursementRepository interface {
	GetReimbursementsGroupedByEmployeeID(ctx context.Context, startDate, endDate time.Time) (map[int][]domain.Reimbursement, error)
}

type AdminService struct {
	adminRepository         AdminRepository
	employeeRepository      EmployeeRepository
	payrollRepository       PayrollRepository
	attendanceRepository    AttendanceRepository
	overtimeRepository      OvertimeRepository
	reimbursementRepository ReimbursementRepository
}

func NewAdminService(adminRepo AdminRepository, empRepo EmployeeRepository,
	payrollRepo PayrollRepository, attendanceRepo AttendanceRepository,
	overtimeRepo OvertimeRepository, reimbursementRepo ReimbursementRepository) *AdminService {
	return &AdminService{
		adminRepository:         adminRepo,
		employeeRepository:      empRepo,
		payrollRepository:       payrollRepo,
		attendanceRepository:    attendanceRepo,
		overtimeRepository:      overtimeRepo,
		reimbursementRepository: reimbursementRepo,
	}
}

func (s *AdminService) LoginAsAdmin(ctx context.Context, credentials dto.LoginRequest) (string, error) {

	admin, err := s.adminRepository.GetAdmin(ctx, domain.Admin{
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
func (s *AdminService) CreatePayrollPeriod(ctx context.Context, payrollPeriodPayload dto.PayrollPeriodRequest) (string, error) {
	start_date, err := time.Parse("2006-01-02", payrollPeriodPayload.StartDate)
	if err != nil {
		return "", error_const.ErrInvalidDateFormat
	}
	end_date, err := time.Parse("2006-01-02", payrollPeriodPayload.EndDate)
	if err != nil {
		return "", error_const.ErrInvalidDateFormat
	}
	if start_date.After(end_date) {
		return "", error_const.ErrStartDateAfterEndDate
	}

	existingPeriod, err := s.payrollRepository.GetPayrollPeriodFromDateRange(ctx, domain.PayrollPeriod{
		StartDate: start_date,
		EndDate:   end_date,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
	}
	if existingPeriod.ID != 0 {
		return "", error_const.ErrPayrollPeriodAlreadyExists
	}

	Id, err := s.payrollRepository.CreatePayrollPeriod(ctx, domain.PayrollPeriod{
		StartDate: start_date,
		EndDate:   end_date,
		CreatedBy: payrollPeriodPayload.ActorEmail,
		UpdatedBy: payrollPeriodPayload.ActorEmail,
	})
	if err != nil {
		return "", err
	}
	return Id, nil
}
func (s *AdminService) RunPayrollPeriod(ctx context.Context, payrollPayload dto.PayrollRequest) error {
	payrollPeriod, err := s.payrollRepository.GetPayrollPeriod(context.Background(), domain.PayrollPeriod{
		ID: payrollPayload.PeriodID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return error_const.ErrPayrollPeriodNotFound
		}
	}
	if payrollPeriod.Locked {
		return error_const.ErrPayrollPeriodLocked
	}
	employees, err := s.employeeRepository.GetAllEmployees(ctx)
	if err != nil {
		return err
	}
	attendance, err := s.attendanceRepository.GetTotalAttendanceByDateRangeGroupedByEmployee(ctx, payrollPeriod.StartDate, payrollPeriod.EndDate)
	if err != nil {
		return err
	}
	overtime, err := s.overtimeRepository.GetOvertimesGroupedByEmployeeID(ctx, payrollPeriod.StartDate, payrollPeriod.EndDate)

	if err != nil {
		return err
	}
	reimbursement, err := s.reimbursementRepository.GetReimbursementsGroupedByEmployeeID(ctx, payrollPeriod.StartDate, payrollPeriod.EndDate)
	if err != nil {
		return err
	}
	allPayrolls := make([]domain.Payroll, 0, len(employees))
	if len(employees) == 0 {
		return error_const.ErrNoEmployeesFound
	}
	for _, employee := range employees {
		var payroll domain.Payroll
		payroll.EmployeeID = employee.ID
		payroll.PeriodID = payrollPeriod.ID
		var payslip domain.Payslip
		payslip.EmployeeID = employee.ID
		payslip.PeriodID = payrollPeriod.ID

		totalWorkDay := utils.GetTotalWorkdays(payrollPeriod.StartDate, payrollPeriod.EndDate)
		payslip.NumberAttendances = attendance[employee.ID]
		payslip.TotalWorkDays = totalWorkDay
		attendanceSalary := employee.Salary * (float64(payslip.NumberAttendances) / float64(totalWorkDay))
		salaryPerHours := employee.Salary / 20 // Assuming 20 working days in a month
		payslip.SalaryByAttendance = attendanceSalary

		var overtimeRecaps []domain.OvertimeRecap
		var overtimeSalary float64
		var overtimeTotalHours int

		overtimeTotalHours = 0
		if len(overtime[employee.ID]) > 0 {
			recapOvertime := overtime[employee.ID]
			for _, o := range recapOvertime {
				amount := float64(o.Hours) * salaryPerHours * 2 // Assuming overtime is paid at double rate
				overtimeRecap := domain.OvertimeRecap{
					Date:   o.Date,
					Hours:  o.Hours,
					Amount: amount,
				}
				overtimeSalary += amount
				overtimeTotalHours += o.Hours
				overtimeRecaps = append(overtimeRecaps, overtimeRecap)
			}
		}
		payslip.OvertimesRecap = overtimeRecaps
		payslip.OvetimeTotalSalary = overtimeSalary

		payslip.Reimbursements = reimbursement[employee.ID]
		totalReimbursement := 0.0
		for _, r := range payslip.Reimbursements {
			totalReimbursement += r.Amount
		}
		payslip.ReimbursementsTotalSalary = totalReimbursement

		payslip.TotalSalary = attendanceSalary + overtimeSalary + totalReimbursement
		payslip.Description = fmt.Sprintf(
			"Total Salary: %.2f\n"+
				"Attendance Salary: %.2f (Base Salary: %.2f x Attendance: %d / Workdays: %d)\n"+
				"Overtime Salary: %.2f (Overtime Hours: %d x Salary/Day: %.2f x 2)\n"+
				"Total Reimbursement: %.2f",
			payslip.TotalSalary,
			attendanceSalary, employee.Salary, payslip.NumberAttendances, totalWorkDay,
			overtimeSalary, overtimeTotalHours, salaryPerHours,
			totalReimbursement,
		)
		payroll.Payslip = payslip
		payroll.CreatedBy = payrollPayload.ActorEmail
		payroll.UpdatedBy = payrollPayload.ActorEmail

		allPayrolls = append(allPayrolls, payroll)
	}

	err = s.payrollRepository.BulkInsertPayrolls(ctx, allPayrolls)

	if err != nil {
		return err
	}
	err = s.payrollRepository.LockPayrollPeriod(ctx, payrollPayload.PeriodID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return error_const.ErrPayrollPeriodNotFound
		}
	}

	return nil
}

func (s *AdminService) ViewPayrollSummary(ctx context.Context, periodID int) (*dto.PayrollSummaryResponse, error) {
	payrolls, err := s.payrollRepository.GetPayrollsByPeriodID(ctx, periodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, error_const.ErrPayrollPeriodNotFound
		}
		return nil, err
	}

	if len(payrolls) == 0 {
		return nil, error_const.ErrNoPayrollsFound
	}

	var totalSalary float64
	employeeSummaries := make([]dto.EmployeePayrollSummary, 0, len(payrolls))

	for _, payroll := range payrolls {
		employee, err := s.employeeRepository.GetEmployeeByID(ctx, payroll.EmployeeID)
		if err != nil {
			return nil, err
		}
		summary := dto.EmployeePayrollSummary{
			EmployeeID:   employee.ID,
			EmployeeName: employee.Name,
			TotalSalary:  payroll.Payslip.TotalSalary,
		}
		employeeSummaries = append(employeeSummaries, summary)
		totalSalary += payroll.Payslip.TotalSalary
	}

	response := &dto.PayrollSummaryResponse{
		PeriodID:          periodID,
		EmployeeSummaries: employeeSummaries,
		TotalSalary:       totalSalary,
	}

	return response, nil
}
