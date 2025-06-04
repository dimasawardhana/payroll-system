package employee_service

import (
	"context"
	"payroll-system/internal/delivery/dto"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"payroll-system/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLoginAsEmployee_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockEmpRepo.Employee = domain.Employee{}
	mockEmpRepo.Err = error_const.ErrInvalidCredentials

	svc := NewEmployeeService(
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

	svc := NewEmployeeService(
		nil, mockPayrollRepo, nil, nil, nil,
	)
	_, err := svc.GetPayslip(context.Background(), dto.PayrollRequest{EmployeeID: 0, PeriodID: 0})
	if err == nil {
		t.Error("expected error for payslip not found")
	}
}
