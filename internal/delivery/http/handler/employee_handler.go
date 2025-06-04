package handler

import (
	"payroll-system/internal/delivery/dto"
	employee_service "payroll-system/internal/service/employee"
	"payroll-system/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	empService *employee_service.EmployeeService
}

func NewEmployeeHandler(empSvc *employee_service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		empService: empSvc,
	}
}
func (h *EmployeeHandler) EmployeeLoginHandler(c *gin.Context) {
	var credential dto.LoginRequest
	if err := c.ShouldBindJSON(&credential); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	token, err := h.empService.LoginAsEmployee(c.Request.Context(), credential)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Login failed", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Employee login successful", dto.LoginResponse{
		Token: token,
	}))
}

func (h *EmployeeHandler) EmployeeAttendanceHandler(c *gin.Context) {
	var attendancePayload dto.AttendanceRequest
	if err := c.ShouldBindJSON(&attendancePayload); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	attendancePayload.EmployeeID = claims.UserID
	attendancePayload.EmployeeEmail = claims.Email
	if err := h.empService.RecordAttendance(c.Request.Context(), attendancePayload); err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to record attendance", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Attendance recorded successfully", nil))
}
func (h *EmployeeHandler) EmployeeOvertimeSubmissionHandler(c *gin.Context) {
	var overtimePayload dto.OvertimeRequest
	if err := c.ShouldBindJSON(&overtimePayload); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	overtimePayload.EmployeeID = claims.UserID
	overtimePayload.EmployeeEmail = claims.Email
	if err := h.empService.SubmitOvertime(c.Request.Context(), overtimePayload); err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to submit overtime", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Overtime submitted successfully", nil))
}
func (h *EmployeeHandler) EmployeeReimbursementHandler(c *gin.Context) {
	var reimbursementPayload dto.ReimbursementRequest
	if err := c.ShouldBindJSON(&reimbursementPayload); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	reimbursementPayload.EmployeeID = claims.UserID
	reimbursementPayload.EmployeeEmail = claims.Email
	if err := h.empService.SubmitReimbursement(c.Request.Context(), reimbursementPayload); err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to submit reimbursement", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Reimbursement submitted successfully", nil))
}
func (h *EmployeeHandler) EmployeePayslipHandler(c *gin.Context) {
	var payslipPayload dto.PayrollRequest
	periodIdStr := c.Param("period_id")

	periodID, err := strconv.Atoi(periodIdStr)
	if err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid period ID", err))
		return
	}
	payslipPayload.PeriodID = periodID

	if payslipPayload.PeriodID == 0 {
		c.JSON(400, dto.NewErrorResponse("Period ID is required", nil))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	payslipPayload.EmployeeID = claims.UserID
	payslipPayload.ActorEmail = claims.Email

	payslip, err := h.empService.GetPayslip(c.Request.Context(), payslipPayload)
	if err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to retrieve payslip", err))
		return
	}

	c.JSON(200, dto.NewSuccessResponse("Payslip retrieved successfully", payslip))
}
