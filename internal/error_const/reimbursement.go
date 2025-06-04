package error_const

import "errors"

var ErrInvalidReimbursementAmount = errors.New("reimbursement amount must be a positive number")
