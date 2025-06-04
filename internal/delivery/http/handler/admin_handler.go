package handler

import (
	"payroll-system/internal/delivery/dto"
	admin_service "payroll-system/internal/service/admin"
	employee_service "payroll-system/internal/service/employee"
	"payroll-system/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminService    *admin_service.AdminService
	EmployeeService *employee_service.EmployeeService
}

func NewAdminHandler(adminSvc *admin_service.AdminService, empSvc *employee_service.EmployeeService) *AdminHandler {
	return &AdminHandler{
		AdminService:    adminSvc,
		EmployeeService: empSvc,
	}
}

func (h *AdminHandler) AdminLoginHandler(c *gin.Context) {
	var credentials dto.LoginRequest

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	token, err := h.AdminService.LoginAsAdmin(c.Request.Context(), credentials)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Login failed", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Admin login successful", dto.LoginResponse{
		Token: token,
	}))
}
func (h *AdminHandler) AdminCreatePayrollPeriodHandler(c *gin.Context) {
	var payrollPeriodPayload dto.PayrollPeriodRequest
	if err := c.ShouldBindJSON(&payrollPeriodPayload); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	payrollPeriodPayload.ActorEmail = claims.Email
	id, err := h.AdminService.CreatePayrollPeriod(c.Request.Context(), payrollPeriodPayload)
	if err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to create payroll period", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Payroll period created successfully", dto.PayrollPeriodResponse{
		ID:        id,
		StartDate: payrollPeriodPayload.StartDate,
		EndDate:   payrollPeriodPayload.EndDate,
	}))
}
func (h *AdminHandler) AdminRunPayrollPeriodHandler(c *gin.Context) {

	var payrollPeriodRunPayload dto.PayrollRequest
	if err := c.ShouldBindJSON(&payrollPeriodRunPayload); err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid request", err))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	payrollPeriodRunPayload.ActorEmail = claims.Email
	err = h.AdminService.RunPayrollPeriod(c.Request.Context(), payrollPeriodRunPayload)
	if err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to run payroll period", err))
		return
	}

	c.JSON(200, dto.NewSuccessResponse("Payroll period run initiated successfully", nil))
}
func (h *AdminHandler) AdminViewPayrollSummaryHandler(c *gin.Context) {
	var payrollSummaryPayload dto.PayrollRequest
	periodIdStr := c.Param("period_id")

	periodID, err := strconv.Atoi(periodIdStr)
	if err != nil {
		c.JSON(400, dto.NewErrorResponse("Invalid period ID", err))
		return
	}
	payrollSummaryPayload.PeriodID = periodID

	if payrollSummaryPayload.PeriodID == 0 {
		c.JSON(400, dto.NewErrorResponse("Period ID is required", nil))
		return
	}
	claims, err := utils.GetClaimsFromJWTUsingContext(c)
	if err != nil {
		c.JSON(401, dto.NewErrorResponse("Unauthorized", err))
		return
	}
	payrollSummaryPayload.ActorEmail = claims.Email
	summary, err := h.AdminService.ViewPayrollSummary(c.Request.Context(), payrollSummaryPayload.PeriodID)
	if err != nil {
		c.JSON(500, dto.NewErrorResponse("Failed to retrieve payroll summary", err))
		return
	}
	c.JSON(200, dto.NewSuccessResponse("Payroll summary retrieved successfully", summary))
}
