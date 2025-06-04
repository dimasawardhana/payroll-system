package error_const

import "errors"

var ErrJWTTokenRequired = errors.New("JWT token is required")
var ErrJWTTokenInvalid = errors.New("JWT token is invalid")
var ErrDBConnFailed = errors.New("failed to connect to the database")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrNotAllowedAccess = errors.New("not allowed to access this resource")
var ErrInvalidDateFormat = errors.New("invalid date format, expected YYYY-MM-DD")
var ErrInvalidUser = errors.New("invalid user")
var ErrInvalidID = errors.New("invalid ID provided")
var ErrInvalidInput = errors.New("invalid input provided")
