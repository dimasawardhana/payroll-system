package dto

type OvertimeRequest struct {
	EmployeeID    int    `json:"employee_id"`
	EmployeeEmail string `json:"employee_email"`
	Date          string `json:"date" binding:"required"`
	Hours         int    `json:"hours" binding:"required"`
}
