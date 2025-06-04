package tests

import (
	"context"
	"payroll-system/internal/delivery/dto"
	"payroll-system/internal/domain"
	"payroll-system/internal/error_const"
	"payroll-system/internal/mocks"
	admin_service "payroll-system/internal/service/admin"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestLoginAsAdmin_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAdminRepo := mocks.NewMockAdminRepository(ctrl)
	mockAdminRepo.Admin = domain.Admin{}
	mockAdminRepo.Err = error_const.ErrInvalidCredentials

	svc := admin_service.NewAdminService(
		mockAdminRepo,
		nil, nil, nil, nil, nil,
	)
	_, err := svc.LoginAsAdmin(context.Background(), dto.LoginRequest{Email: "", Password: ""})
	if err == nil {
		t.Error("expected error for invalid credentials")
	}
}

func TestViewPayrollSummary_NoPayrolls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPayrollRepo := mocks.NewMockPayrollRepository(ctrl)
	mockPayrollRepo.Payrolls = []domain.Payroll{}
	mockPayrollRepo.Err = nil
	mockEmpRepo := mocks.NewMockEmployeeRepository(ctrl)
	mockAttendanceRepo := mocks.NewMockAttendanceRepository(ctrl)
	mockOvertimeRepo := mocks.NewMockOvertimeRepository(ctrl)
	mockReimbursementRepo := mocks.NewMockReimbursementRepository(ctrl)

	svc := admin_service.NewAdminService(
		nil,
		mockEmpRepo,
		mockPayrollRepo,
		mockAttendanceRepo,
		mockOvertimeRepo,
		mockReimbursementRepo,
	)
	_, err := svc.ViewPayrollSummary(context.Background(), 9999)
	if err != error_const.ErrNoPayrollsFound {
		t.Errorf("expected ErrNoPayrollsFound, got %v", err)
	}
}
