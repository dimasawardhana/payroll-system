package dto

type PayrollPeriodRequest struct {
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	ActorEmail string `json:"actor_id"`
}

type PayrollPeriodResponse struct {
	ID        string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type PayrollRequest struct {
	PeriodID   int    `json:"period_id" binding:"required"`
	ActorEmail string `json:"actor_email"`
	EmployeeID int    `json:"employee_id"` // Optional, if running for a specific employee
}
