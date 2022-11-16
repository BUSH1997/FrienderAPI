package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	New    = errors.New
	Errorf = errors.Errorf
	As     = errors.As

	Unwrap = errors.Unwrap
)

type ErrorWithStackTrace interface {
	error
	StackTrace() errors.StackTrace
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	var stackTracedErr ErrorWithStackTrace
	if hasStackTrace := As(err, &stackTracedErr); hasStackTrace {
		return errors.WithMessage(err, message)
	}

	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var stackTracedErr ErrorWithStackTrace
	if hasStackTrace := As(err, &stackTracedErr); hasStackTrace {
		return errors.WithMessagef(err, format, args...)
	}

	return errors.Wrapf(err, format, args...)
}

func Is(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

var _ error = transformedError{}

type transformedError struct {
	TargetError   error
	OriginalError error
}

func (err transformedError) Error() string {
	return fmt.Sprintf("%s: %s", err.TargetError, err.OriginalError)
}

func (err transformedError) Is(another error) bool {
	return Is(err.TargetError, another)
}

func (err transformedError) As(dst interface{}) bool {
	return As(err.TargetError, dst)
}

func (err transformedError) Unwrap() error {
	return err.OriginalError
}

func Transform(originalError error, targetError error) error {
	if originalError == nil {
		return nil
	}

	return transformedError{
		TargetError:   targetError,
		OriginalError: originalError,
	}
}
