package error_const

import "errors"

var ErrEmptyPayrollPeriod = errors.New("payroll period cannot be empty")
var ErrStartDateAfterEndDate = errors.New("start date cannot be after end date")
var ErrPayrollPeriodNotFound = errors.New("payroll period not found")
var ErrPayrollPeriodLocked = errors.New("payroll period is locked, cannot modify")
var ErrPayrollPeriodAlreadyExists = errors.New("payroll period already exists for the given date range")
var ErrPayrollPeriodInvalid = errors.New("payroll period is invalid, check start and end dates")
var ErrNoEmployeesFound = errors.New("no employees found for payroll period")
var ErrPayslipNotFound = errors.New("payslip not found for the given employee and period")
var ErrNoPayrollsFound = errors.New("no payrolls found for this period")
