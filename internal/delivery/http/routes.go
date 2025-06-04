package http

import (
	"payroll-system/internal/delivery/http/handler"
	middleware "payroll-system/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	router          *gin.Engine
	adminHandler    *handler.AdminHandler
	employeeHandler *handler.EmployeeHandler
}

func NewRoutes(router *gin.Engine, adminHandler *handler.AdminHandler, employeeHandler *handler.EmployeeHandler) *Routes {
	return &Routes{
		router:          router,
		adminHandler:    adminHandler,
		employeeHandler: employeeHandler,
	}
}

func (r *Routes) InitRoutes() {
	r.router.GET("/healthz", handler.HealthCheckHandler)

	httpV1 := r.router.Group("/api/v1")
	login := httpV1.Group("/login")
	{
		login.POST("/admin", r.adminHandler.AdminLoginHandler)
		login.POST("/employee", r.employeeHandler.EmployeeLoginHandler)
	}
	httpV1.Use(middleware.CheckJWT())
	{
		RegisterAdminRoutes(httpV1, r.adminHandler)
		RegisterEmployeeRoutes(httpV1, r.employeeHandler)
	}
}

func RegisterAdminRoutes(router *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.CheckRole("admin"))
	{
		adminGroup.POST("/payroll-period", adminHandler.AdminCreatePayrollPeriodHandler)
		adminGroup.POST("/payroll-period/run", adminHandler.AdminRunPayrollPeriodHandler)
		adminGroup.GET("/payroll-summary/:period_id", adminHandler.AdminViewPayrollSummaryHandler)
	}
}

func RegisterEmployeeRoutes(router *gin.RouterGroup, employeeHandler *handler.EmployeeHandler) {
	employeeGroup := router.Group("/employee")
	employeeGroup.Use(middleware.CheckRole("employee"))
	{
		employeeGroup.POST("/attendance", employeeHandler.EmployeeAttendanceHandler)
		employeeGroup.POST("/overtime", employeeHandler.EmployeeOvertimeSubmissionHandler)
		employeeGroup.POST("/reimbursement", employeeHandler.EmployeeReimbursementHandler)
		employeeGroup.GET("/payslip/:period_id", employeeHandler.EmployeePayslipHandler)
	}
}
