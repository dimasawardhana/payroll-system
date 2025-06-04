package dto

type ReimbursementRequest struct {
	EmployeeID    int     `json:"employee_id"`
	EmployeeEmail string  `json:"employee_email"`
	Amount        float64 `json:"amount" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Date          string  `json:"date" binding:"required"`
}
