package error_const

import "errors"

var ErrInvalidOvertimeHours = errors.New("overtime hours must be a positive number")
var ErrOvertimeHoursExceeded = errors.New("overtime hours cannot exceed 3 hours per day")
var ErrOvertimeAlreadyExists = errors.New("overtime record already exists for the given date")
