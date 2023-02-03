package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/rpc"
)

// StackErrors wraps multiple errors into a single error.
func StackErrors(errs ...error) error {
	return NewStakedError(errs...)
}

// StakedError is a struct that holds multiple errors and returns them as a single error.
type StakedError struct {
	errors []error
}

// NewStakedError creates a new StakedError instance.
func NewStakedError(errs ...error) *StakedError {
	return &StakedError{errors: errs}
}

// Error returns the error message.
// Implements the error interface.
func (e *StakedError) Error() string {
	if len(e.errors) == 0 {
		return ""
	}

	if len(e.errors) == 1 {
		return errToString(e.errors[0])
	}

	var result string

	for _, e := range e.errors {
		if e == nil {
			continue
		}

		if result == "" {
			result = errToString(e)
			continue
		}

		result = fmt.Sprintf("%s: %s", result, errToString(e))
	}

	return result
}

// Unwrap returns the underlying error.
// Implements the errors.Unwrap interface.
func (e *StakedError) Unwrap() error {
	if len(e.errors) == 0 {
		return nil
	}

	if len(e.errors) == 1 {
		return e.errors[0]
	}

	return e
}

// Is returns true if the error is equal to target.
func (e *StakedError) Is(target error) bool {
	if len(e.errors) == 0 {
		return false
	}

	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// As returns the first error that matches target.
func (e *StakedError) As(target interface{}) bool {
	if len(e.errors) == 0 {
		return false
	}

	for _, err := range e.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

// errToString converts the error to a string.
func errToString(err error) string {
	if rpcErr, ok := err.(*rpc.JsonRpcError); ok {
		log.Printf("JsonRpcError: %+v", rpcErr)
		return rpcErr.Message
	}

	return err.Error()
}

// WrapError wraps the rpc.JsonRpcError to a standard error.
func WrapError(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	if len(errs) == 1 {
		return errs[0]
	}

	var err error

	for _, e := range errs {
		if e == nil {
			continue
		}

		if err == nil {
			err = e
			continue
		}

		if rpcErr, ok := e.(*rpc.JsonRpcError); ok {
			err = fmt.Errorf("%w: %s", err, rpcErr.Message)
			continue
		}

		err = fmt.Errorf("%w: %s", err, e.Error())
	}

	return err
}
