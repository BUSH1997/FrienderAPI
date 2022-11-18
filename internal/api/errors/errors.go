package errors

import (
	"fmt"
)

type ApiError interface {
	error
	ErrorType() string
	ErrorExplain() string
	Unwrap() error
}

type ValidationError struct {
	Cause                 error
	Explain               string
	ValidationErrorFields ValidationErrorFields
}

func (e ValidationError) Unwrap() error {
	return e.Cause
}

func (e ValidationError) ErrorType() string {
	return "invalid"
}

func (e ValidationError) ErrorExplain() string {
	return e.Explain
}

func (e ValidationError) Error() string {
	return errorString(e)
}

type ValidationReason string

const (
	FieldForbidden   ValidationReason = "forbidden"
	FieldInvalid     ValidationReason = "invalid"
	FieldNotFound    ValidationReason = "not_found"
	FieldConflict    ValidationReason = "conflict"
	FieldRequired    ValidationReason = "required"
	FieldRequiredAny ValidationReason = "required_any"
)

type ValidationErrorFields map[string]ValidationReason

func (v ValidationErrorFields) AddReason(field string, reason ValidationReason) {
	v[field] = reason
}

type InternalError struct {
	Cause error
}

func (e InternalError) ErrorType() string {
	return "internal"
}

func (e InternalError) Unwrap() error {
	return e.Cause
}

func (e InternalError) ErrorExplain() string {
	return "internal server error"
}

func (e InternalError) Error() string {
	return errorString(e)
}

type ForbiddenReason string

const (
	CSRFTokenInvalid ForbiddenReason = "csrf_invalid"
	EventIsSpam      ForbiddenReason = "event_is_spam"
	NoAccess         ForbiddenReason = "no_access"
)

type ForbiddenError struct {
	Cause  error
	Reason ForbiddenReason
}

func (e ForbiddenError) ErrorType() string {
	return "forbidden"
}

func (e ForbiddenError) Unwrap() error {
	return e.Cause
}

func (e ForbiddenError) ErrorExplain() string {
	return fmt.Sprintf("forbidden due to %s", e.Reason)
}

func (e ForbiddenError) Error() string {
	return errorString(e)
}

type NotFoundError struct {
	Cause error
	Key   string
	Value string
}

func (e NotFoundError) ErrorType() string {
	return "not_found"
}

func (e NotFoundError) Unwrap() error {
	return e.Cause
}

func (e NotFoundError) ErrorExplain() string {
	return fmt.Sprintf("%s was not found by value %s", e.Key, e.Value)
}

func (e NotFoundError) Error() string {
	return errorString(e)
}

type FormatError struct {
	Cause  error
	Value  interface{}
	Format string
}

func (e FormatError) ErrorType() string {
	return "format"
}

func (e FormatError) Unwrap() error {
	return e.Cause
}

func (e FormatError) ErrorExplain() string {
	return fmt.Sprintf("%q doesn't match format %q", e.Value, e.Format)
}

func (e FormatError) Error() string {
	return errorString(e)
}

type UnauthorizedReason string

const (
	UnauthorizedReasonNoCookies       UnauthorizedReason = "no_cookies"
	UnauthorizedReasonTokenInvalid    UnauthorizedReason = "token_invalid"
	UnauthorizedReasonAuthError       UnauthorizedReason = "auth_error"
	UnauthorizedReasonWrongAuthMethod UnauthorizedReason = "wrong_authorization_method"
)

type UnauthorizedError struct {
	Cause  error
	Reason UnauthorizedReason
}

func (e UnauthorizedError) ErrorType() string {
	return "unauthorized"
}

func (e UnauthorizedError) Unwrap() error {
	return e.Cause
}

func (e UnauthorizedError) ErrorExplain() string {
	return fmt.Sprintf("user is not authorized due to %s", e.Reason)
}

func (e UnauthorizedError) Error() string {
	return errorString(e)
}

type UserNotExistsError struct {
	Cause error
}

func (e UserNotExistsError) ErrorType() string {
	return "user_not_exists"
}

func (e UserNotExistsError) OriginalError() error {
	return e.Cause
}

func (e UserNotExistsError) Error() string {
	return errorString(e)
}

func (e UserNotExistsError) Unwrap() error {
	return e.Cause
}

func (e UserNotExistsError) ErrorExplain() string {
	return "user not exists"
}

type UnimplementedMethodError struct {
	Cause error
}

func (e UnimplementedMethodError) ErrorType() string {
	return "unimplemented_method_error"
}

func (e UnimplementedMethodError) OriginalError() error {
	return e.Cause
}

func (e UnimplementedMethodError) Error() string {
	return errorString(e)
}

func (e UnimplementedMethodError) Unwrap() error {
	return e.Cause
}

func (e UnimplementedMethodError) ErrorExplain() string {
	return "unimplemented"
}

type AlreadyExistsError struct {
	Cause  error
	Entity string
}

func (e AlreadyExistsError) ErrorType() string {
	return "already_exists"
}

func (e AlreadyExistsError) Unwrap() error {
	return e.Cause
}

func (e AlreadyExistsError) ErrorExplain() string {
	return fmt.Sprintf("%s already exists", e.Entity)
}

func (e AlreadyExistsError) Error() string {
	return errorString(e)
}

func errorString(err ApiError) string {
	s := fmt.Sprintf("api %s error", err.ErrorType())

	explain := err.ErrorExplain()
	if explain != "" {
		s += ": " + explain
	}

	cause := err.Unwrap()
	if cause != nil {
		s += ": " + cause.Error()
	}

	return s
}
