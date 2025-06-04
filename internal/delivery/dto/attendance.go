package dto

type AttendanceRequest struct {
	Date          string `json:"date" binding:"required"`
	EmployeeID    int    `json:"employee_id"`
	EmployeeEmail string `json:"employee_email"`
}
