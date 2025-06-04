package error_const

import "errors"

var ErrAttendanceAlreadyExists = errors.New("attendance already exists for the given date and employee")
var ErrAttendanceOnWeekend = errors.New("attendance cannot be recorded on weekends")
