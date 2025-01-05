package util

import (
	"errors"
	"fmt"
)

type RomanError interface {
	// Error is to enforce the implementation of the Error interface of STD lib. Should return what FullError is in string format
	Error() string

	// FullError Should return all the data stored, even the error cases if the context wasn't specified
	FullError() error
	// AddErrorCase Stores cases without modifying the full error until a context is applied
	AddErrorCase(err error)

	// DisplayError should be a user-readable error without exposing inner workings
	DisplayError() string
	SetDisplayError(err string)

	// WriteCurrentContext should be stored and applied as a prefix to the constructed error + apply the stored error cases
	WriteCurrentContext(context string)
}

type romanError struct {
	fullError         error
	contextErrorCases string
	displayError      string
}

func (r *romanError) Error() string {
	temp := r.fullError
	if r.contextErrorCases != "" {
		temp = fmt.Errorf("%s: %w", r.contextErrorCases, r.fullError)
	}

	return temp.Error()
}

func (r *romanError) FullError() error {
	temp := r.fullError
	if r.contextErrorCases != "" {
		temp = fmt.Errorf("%s %w", r.contextErrorCases, r.fullError)
	}

	return temp
}

func (r *romanError) AddErrorCase(err error) {
	if r.contextErrorCases == "" {
		r.contextErrorCases = err.Error()
		return
	}

	r.contextErrorCases = fmt.Sprintf("%s, %s", err.Error(), r.contextErrorCases)
}

func (r *romanError) DisplayError() string {
	return r.displayError
}

func (r *romanError) SetDisplayError(err string) {
	r.displayError = err
}

func (r *romanError) WriteCurrentContext(context string) {
	contextError := fmt.Sprintf("%s %s", context, r.contextErrorCases)
	r.fullError = fmt.Errorf("%s %w", contextError, r.fullError)

	r.contextErrorCases = ""
}

// NewError returns nil if initial error is nil
func NewError(context string, err error) RomanError {
	if err == nil {
		return nil
	}

	outErr := romanError{}
	outErr.fullError = fmt.Errorf("%s %w", context, err)

	var romErr RomanError
	if errors.As(err, &romErr) {
		outErr.SetDisplayError(romErr.DisplayError())
	}

	return &outErr
}

// NewErrorWithDisplay returns nil if initial error is nil
func NewErrorWithDisplay(context string, fullErr error, display string) RomanError {
	err := NewError(context, fullErr)
	if err == nil {
		return nil
	}

	err.SetDisplayError(display)
	return err
}
