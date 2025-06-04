package main

import (
	"log"
	"payroll-system/internal/config"
	httpRoutes "payroll-system/internal/delivery/http"
	"payroll-system/internal/delivery/http/handler"
	"payroll-system/internal/repository/postgres"
	admin_service "payroll-system/internal/service/admin"
	employee_service "payroll-system/internal/service/employee"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Payslip Generation System server starting...")

	_config := config.Load()

	if _config.Env != "production" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	pool := config.InitDB(_config.DBUrl)
	defer pool.Close()

	router := gin.Default()

	adminRepo := postgres.NewAdminRepository(pool)
	employeeRepo := postgres.NewEmployeeRepository(pool)
	payrollRepo := postgres.NewPayrollRepository(pool)
	attendanceRepo := postgres.NewAttendanceRepository(pool)
	overtimeRepo := postgres.NewOvertimeRepository(pool)
	reimbursementRepo := postgres.NewReimbursementRepository(pool)

	adminService := admin_service.NewAdminService(adminRepo, employeeRepo, payrollRepo, attendanceRepo, overtimeRepo, reimbursementRepo)
	empService := employee_service.NewEmployeeService(employeeRepo, payrollRepo, attendanceRepo, overtimeRepo, reimbursementRepo)

	adminHandler := handler.NewAdminHandler(adminService, empService)
	employeeHandler := handler.NewEmployeeHandler(empService)

	_http := httpRoutes.NewRoutes(router, adminHandler, employeeHandler)
	_http.InitRoutes()
	port := _config.ServerPort
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
