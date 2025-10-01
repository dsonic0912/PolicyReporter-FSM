package fsm

import (
	"fmt"
	"strings"
)

// ErrorType represents different categories of errors that can occur.
type ErrorType int

const (
	// ErrorTypeValidation indicates a validation error
	ErrorTypeValidation ErrorType = iota
	// ErrorTypeTransition indicates a transition error
	ErrorTypeTransition
	// ErrorTypeInvalidInput indicates invalid input
	ErrorTypeInvalidInput
	// ErrorTypeInvalidConfiguration indicates invalid configuration
	ErrorTypeInvalidConfiguration
	// ErrorTypeInternal indicates an internal error
	ErrorTypeInternal
)

// String returns a string representation of the error type.
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeValidation:
		return "ValidationError"
	case ErrorTypeTransition:
		return "TransitionError"
	case ErrorTypeInvalidInput:
		return "InvalidInputError"
	case ErrorTypeInvalidConfiguration:
		return "InvalidConfigurationError"
	case ErrorTypeInternal:
		return "InternalError"
	default:
		return "UnknownError"
	}
}

// AutomatonError represents a structured error with context.
type AutomatonError struct {
	Type    ErrorType
	Message string
	Context map[string]interface{}
	Cause   error
}

// Error implements the error interface.
func (e *AutomatonError) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[%s]", e.Type.String()))
	parts = append(parts, e.Message)

	if len(e.Context) > 0 {
		var contextParts []string
		for key, value := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
		}
		parts = append(parts, fmt.Sprintf("context: {%s}", strings.Join(contextParts, ", ")))
	}

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}

	return strings.Join(parts, " ")
}

// Unwrap returns the underlying cause of the error.
func (e *AutomatonError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a target error type.
func (e *AutomatonError) Is(target error) bool {
	if targetErr, ok := target.(*AutomatonError); ok {
		return e.Type == targetErr.Type
	}
	return false
}

// NewError creates a new AutomatonError.
func NewError(errorType ErrorType, message string) *AutomatonError {
	return &AutomatonError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// NewErrorWithContext creates a new AutomatonError with context.
func NewErrorWithContext(errorType ErrorType, message string, context map[string]interface{}) *AutomatonError {
	return &AutomatonError{
		Type:    errorType,
		Message: message,
		Context: context,
	}
}

// NewErrorWithCause creates a new AutomatonError with a cause.
func NewErrorWithCause(errorType ErrorType, message string, cause error) *AutomatonError {
	return &AutomatonError{
		Type:    errorType,
		Message: message,
		Context: make(map[string]interface{}),
		Cause:   cause,
	}
}

// WithContext adds context to an existing error.
func (e *AutomatonError) WithContext(key string, value interface{}) *AutomatonError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithCause adds a cause to an existing error.
func (e *AutomatonError) WithCause(cause error) *AutomatonError {
	e.Cause = cause
	return e
}

// Predefined error constructors for common scenarios

// NewValidationError creates a validation error.
func NewValidationError(message string) *AutomatonError {
	return NewError(ErrorTypeValidation, message)
}

// NewTransitionError creates a transition error with context.
func NewTransitionError[Q State, S Symbol](state Q, symbol S, message string) *AutomatonError {
	return NewErrorWithContext(ErrorTypeTransition, message, map[string]interface{}{
		"state":  state,
		"symbol": symbol,
	})
}

// NewInvalidInputError creates an invalid input error with context.
func NewInvalidInputError[S Symbol](symbol S, position int, message string) *AutomatonError {
	return NewErrorWithContext(ErrorTypeInvalidInput, message, map[string]interface{}{
		"symbol":   symbol,
		"position": position,
	})
}

// NewInvalidConfigurationError creates an invalid configuration error.
func NewInvalidConfigurationError(component string, message string) *AutomatonError {
	return NewErrorWithContext(ErrorTypeInvalidConfiguration, message, map[string]interface{}{
		"component": component,
	})
}

// ErrorCollector collects multiple errors and presents them as a single error.
type ErrorCollector struct {
	errors []error
}

// NewErrorCollector creates a new error collector.
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{
		errors: make([]error, 0),
	}
}

// Add adds an error to the collector.
func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

// HasErrors returns true if there are any collected errors.
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// Error returns a combined error message.
func (ec *ErrorCollector) Error() string {
	if len(ec.errors) == 0 {
		return ""
	}

	if len(ec.errors) == 1 {
		return ec.errors[0].Error()
	}

	var messages []string
	for i, err := range ec.errors {
		messages = append(messages, fmt.Sprintf("%d. %s", i+1, err.Error()))
	}

	return fmt.Sprintf("Multiple errors occurred:\n%s", strings.Join(messages, "\n"))
}

// Errors returns all collected errors.
func (ec *ErrorCollector) Errors() []error {
	return ec.errors
}

// ToError returns the collector as an error if there are any errors, nil otherwise.
func (ec *ErrorCollector) ToError() error {
	if ec.HasErrors() {
		return ec
	}
	return nil
}

// ErrorHandler defines an interface for handling errors.
type ErrorHandler interface {
	HandleError(err error) error
}

// LoggingErrorHandler logs errors and optionally transforms them.
type LoggingErrorHandler struct {
	logger func(string)
}

// NewLoggingErrorHandler creates a new logging error handler.
func NewLoggingErrorHandler(logger func(string)) *LoggingErrorHandler {
	return &LoggingErrorHandler{
		logger: logger,
	}
}

// HandleError logs the error and returns it unchanged.
func (h *LoggingErrorHandler) HandleError(err error) error {
	if h.logger != nil {
		h.logger(fmt.Sprintf("Automaton error: %v", err))
	}
	return err
}

// RetryableErrorHandler wraps errors with retry information.
type RetryableErrorHandler struct {
	maxRetries int
}

// NewRetryableErrorHandler creates a new retryable error handler.
func NewRetryableErrorHandler(maxRetries int) *RetryableErrorHandler {
	return &RetryableErrorHandler{
		maxRetries: maxRetries,
	}
}

// HandleError wraps the error with retry information.
func (h *RetryableErrorHandler) HandleError(err error) error {
	if automatonErr, ok := err.(*AutomatonError); ok {
		return automatonErr.WithContext("max_retries", h.maxRetries)
	}
	return err
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) bool {
	if automatonErr, ok := err.(*AutomatonError); ok {
		return automatonErr.Type == ErrorTypeValidation
	}
	return false
}

// IsTransitionError checks if an error is a transition error.
func IsTransitionError(err error) bool {
	if automatonErr, ok := err.(*AutomatonError); ok {
		return automatonErr.Type == ErrorTypeTransition
	}
	return false
}

// IsInvalidInputError checks if an error is an invalid input error.
func IsInvalidInputError(err error) bool {
	if automatonErr, ok := err.(*AutomatonError); ok {
		return automatonErr.Type == ErrorTypeInvalidInput
	}
	return false
}
