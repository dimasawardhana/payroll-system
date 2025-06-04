package tests

import (
	"context"
	"payroll-system/internal/delivery/dto"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"payroll-system/internal/mocks"
	employee_service "payroll-system/internal/service/employee"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLoginAsEmployee_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockEmpRepo.Employee = domain.Employee{}
	mockEmpRepo.Err = error_const.ErrInvalidCredentials

	svc := employee_service.NewEmployeeService(
		mockEmpRepo, nil, nil, nil, nil,
	)
	_, err := svc.LoginAsEmployee(context.Background(), dto.LoginRequest{Email: "", Password: ""})
	if err == nil {
		t.Error("expected error for invalid credentials")
	}
}

func TestGetPayslip_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPayrollRepo := mocks.NewMockPayrollRepository(ctrl)
	mockPayrollRepo.Err = error_const.ErrPayslipNotFound

	svc := employee_service.NewEmployeeService(
		nil, mockPayrollRepo, nil, nil, nil,
	)
	_, err := svc.GetPayslip(context.Background(), dto.PayrollRequest{EmployeeID: 0, PeriodID: 0})
	if err == nil {
		t.Error("expected error for payslip not found")
	}
}

func TestRecordAttendance_InvalidEmployeeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockAttendanceRepo := mocks.NewMockAttendanceRepository(ctrl)

	svc := employee_service.NewEmployeeService(
		mockEmpRepo, nil, mockAttendanceRepo, nil, nil,
	)
	err := svc.RecordAttendance(context.Background(), dto.AttendanceRequest{EmployeeID: 0, Date: "2025-06-04"})
	if err != error_const.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestSubmitOvertime_InvalidHours(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockOvertimeRepo := mocks.NewMockOvertimeRepository(ctrl)

	svc := employee_service.NewEmployeeService(
		mockEmpRepo, nil, nil, mockOvertimeRepo, nil,
	)
	err := svc.SubmitOvertime(context.Background(), dto.OvertimeRequest{EmployeeID: 1, Hours: 0})
	if err != error_const.ErrInvalidOvertimeHours {
		t.Errorf("expected ErrInvalidOvertimeHours, got %v", err)
	}
}

func TestSubmitReimbursement_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockReimbursementRepo := mocks.NewMockReimbursementRepository(ctrl)

	svc := employee_service.NewEmployeeService(
		mockEmpRepo, nil, nil, nil, mockReimbursementRepo,
	)
	err := svc.SubmitReimbursement(context.Background(), dto.ReimbursementRequest{EmployeeID: 1, Amount: 0})
	if err != error_const.ErrInvalidReimbursementAmount {
		t.Errorf("expected ErrInvalidReimbursementAmount, got %v", err)
	}
}
