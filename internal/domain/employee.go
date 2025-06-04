package domain

import "time"

type Employee struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password_hash string  `json:"password_hash"`
	Role          string  `json:"role"`
	Salary        float64 `json:"salary"`
	Created_at    string  `json:"created_at"`
	Updated_at    string  `json:"updated_at"`
	Created_by    string  `json:"created_by"`
	Updated_by    string  `json:"updated_by"`
}

type Attendance struct {
	ID         int       `json:"id"`
	Date       time.Time `json:"date"`
	EmployeeID int       `json:"employee_id"`
	Status     string    `json:"status"` // e.g., "present", "absent", "leave"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  string    `json:"updated_by"`
}

type Overtime struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	Hours      int       `json:"hours"`
	Date       time.Time `json:"date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  string    `json:"updated_by"`
}

type OvertimeRecap struct {
	Date   time.Time `json:"date"`
	Hours  int       `json:"hours"`
	Amount float64   `json:"amount"` // Total salary for the overtime hours
}
type Reimbursement struct {
	ID          int       `json:"id"`
	EmployeeID  int       `json:"employee_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
}
type Payslip struct {
	ID                        int             `json:"id"`
	EmployeeID                int             `json:"employee_id"`
	PeriodID                  int             `json:"period_id"`
	NumberAttendances         int             `json:"num_attendances"`
	TotalWorkDays             int             `json:"total_work_days"`
	SalaryByAttendance        float64         `json:"salary_by_attendance"`
	OvertimesRecap            []OvertimeRecap `json:"overtimes_recap"`
	OvetimeTotalSalary        float64         `json:"overtime_total_salary"`
	Reimbursements            []Reimbursement `json:"reimbursements"`
	ReimbursementsTotalSalary float64         `json:"reimbursements_total_salary"`
	TotalSalary               float64         `json:"total_salary"`
	Description               string          `json:"description"`
}
type Payroll struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	PeriodID   int       `json:"period_id"`
	Payslip    Payslip   `json:"payslip"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  string    `json:"updated_by"`
}
