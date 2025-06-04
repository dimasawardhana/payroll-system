package dto

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Message: message,
		Data:    data,
	}
}
func NewErrorResponse(message string, err error) *Response {
	return &Response{
		Message: message,
		Error:   err.Error(),
	}
}

type EmployeePayrollSummary struct {
	EmployeeID   int     `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	TotalSalary  float64 `json:"total_salary"`
}

type PayrollSummaryResponse struct {
	PeriodID          int                      `json:"period_id"`
	EmployeeSummaries []EmployeePayrollSummary `json:"employee_summaries"`
	TotalSalary       float64                  `json:"total_salary"`
}
