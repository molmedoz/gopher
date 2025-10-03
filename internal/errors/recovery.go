package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Recoverer provides panic recovery functionality
type Recoverer struct {
	logger *ErrorLogger
}

// NewRecoverer creates a new recoverer
func NewRecoverer(logger *ErrorLogger) *Recoverer {
	return &Recoverer{
		logger: logger,
	}
}

// Recover recovers from a panic and converts it to an error
func (r *Recoverer) Recover() error {
	if recovered := recover(); recovered != nil {
		// Get stack trace
		stack := debug.Stack()

		// Create a GopherError for the panic
		err := Newf(ErrCodeUnknown, "panic occurred: %v", recovered)
		err.WithDetails(fmt.Sprintf("Stack trace:\n%s", string(stack)))

		// Log the panic
		if r.logger != nil {
			r.logger.LogGopherError(err, map[string]interface{}{
				"panic_value": recovered,
				"goroutine":   runtime.NumGoroutine(),
			})
		}

		return err
	}
	return nil
}

// RecoverWithHandler recovers from a panic and calls a handler function
func (r *Recoverer) RecoverWithHandler(handler func(error)) {
	if err := r.Recover(); err != nil {
		handler(err)
	}
}

// SafeExecute executes a function with panic recovery
func (r *Recoverer) SafeExecute(fn func() error) (err error) {
	defer func() {
		if recovered := r.Recover(); recovered != nil {
			err = recovered
		}
	}()

	return fn()
}

// SafeExecuteWithResult executes a function with panic recovery and returns a result
func SafeExecuteWithResult[T any](fn func() (T, error)) (result T, err error) {
	recoverer := NewRecoverer(nil)
	defer func() {
		if recovered := recoverer.Recover(); recovered != nil {
			err = recovered
		}
	}()

	return fn()
}

// RecoverFunc is a function type that can be used with defer
type RecoverFunc func()

// Recover creates a recovery function that can be used with defer
func Recover(logger *ErrorLogger) RecoverFunc {
	recoverer := NewRecoverer(logger)
	return func() {
		recoverer.Recover()
	}
}

// RecoverWithHandler creates a recovery function with a custom handler
func RecoverWithHandler(logger *ErrorLogger, handler func(error)) RecoverFunc {
	recoverer := NewRecoverer(logger)
	return func() {
		recoverer.RecoverWithHandler(handler)
	}
}

// Must panics if the error is not nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustValue panics if the error is not nil, otherwise returns the value
func MustValue[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Must2 panics if the second return value (error) is not nil
func Must2[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Must3 panics if the third return value (error) is not nil
func Must3[T1, T2 any](value1 T1, value2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return value1, value2
}
