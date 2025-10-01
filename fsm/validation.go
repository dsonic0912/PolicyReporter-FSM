package fsm

import (
	"fmt"
	"reflect"
	"strings"
)

// Default validation limits
const (
	DefaultMaxStates       = 1000
	DefaultMaxAlphabetSize = 100
	DefaultMaxTransitions  = 10000
	StrictMaxStates        = 100
	StrictMaxAlphabetSize  = 50
	StrictMaxTransitions   = 1000
	MaxStateNameLength     = 50
)

// ValidationRule represents a validation rule for automaton construction.
type ValidationRule[Q State, S Symbol] func(*FiniteAutomaton[Q, S]) error

// ValidatorConfig holds configuration for validation behavior.
type ValidatorConfig struct {
	// StrictMode enables additional validation checks
	StrictMode bool
	// MaxStates limits the number of states (0 = no limit)
	MaxStates int
	// MaxAlphabetSize limits the alphabet size (0 = no limit)
	MaxAlphabetSize int
	// MaxTransitions limits the number of transitions (0 = no limit)
	MaxTransitions int
	// RequireCompleteTransitions ensures all state-symbol combinations have transitions
	RequireCompleteTransitions bool
}

// DefaultValidatorConfig returns a default validation configuration.
func DefaultValidatorConfig() ValidatorConfig {
	return ValidatorConfig{
		StrictMode:                 false,
		MaxStates:                  DefaultMaxStates,
		MaxAlphabetSize:            DefaultMaxAlphabetSize,
		MaxTransitions:             DefaultMaxTransitions,
		RequireCompleteTransitions: false,
	}
}

// StrictValidatorConfig returns a strict validation configuration.
func StrictValidatorConfig() ValidatorConfig {
	return ValidatorConfig{
		StrictMode:                 true,
		MaxStates:                  StrictMaxStates,
		MaxAlphabetSize:            StrictMaxAlphabetSize,
		MaxTransitions:             StrictMaxTransitions,
		RequireCompleteTransitions: true,
	}
}

// InputValidator provides comprehensive input validation.
type InputValidator[Q State, S Symbol] struct {
	config ValidatorConfig
	rules  []ValidationRule[Q, S]
}

// NewInputValidator creates a new input validator with the given configuration.
func NewInputValidator[Q State, S Symbol](config ValidatorConfig) *InputValidator[Q, S] {
	validator := &InputValidator[Q, S]{
		config: config,
		rules:  make([]ValidationRule[Q, S], 0),
	}

	// Add default validation rules
	validator.AddRule(validateNonEmptyStates[Q, S])
	validator.AddRule(validateNonEmptyAlphabet[Q, S])
	validator.AddRule(validateInitialStateInStates[Q, S])
	validator.AddRule(validateAcceptingStatesInStates[Q, S])
	validator.AddRule(validateTransitionStatesInStates[Q, S])
	validator.AddRule(validateTransitionSymbolsInAlphabet[Q, S])

	if config.MaxStates > 0 {
		validator.AddRule(validateMaxStates[Q, S](config.MaxStates))
	}

	if config.MaxAlphabetSize > 0 {
		validator.AddRule(validateMaxAlphabetSize[Q, S](config.MaxAlphabetSize))
	}

	if config.MaxTransitions > 0 {
		validator.AddRule(validateMaxTransitions[Q, S](config.MaxTransitions))
	}

	if config.RequireCompleteTransitions {
		validator.AddRule(validateCompleteTransitions[Q, S])
	}

	if config.StrictMode {
		validator.AddRule(validateNoUnreachableStates[Q, S])
		validator.AddRule(validateNoDuplicateTransitions[Q, S])
		validator.AddRule(validateStateNaming[Q, S])
	}

	return validator
}

// AddRule adds a custom validation rule.
func (v *InputValidator[Q, S]) AddRule(rule ValidationRule[Q, S]) {
	v.rules = append(v.rules, rule)
}

// Validate runs all validation rules against the automaton.
func (v *InputValidator[Q, S]) Validate(automaton *FiniteAutomaton[Q, S]) error {
	collector := NewErrorCollector()

	for _, rule := range v.rules {
		if err := rule(automaton); err != nil {
			collector.Add(err)
		}
	}

	return collector.ToError()
}

// Built-in validation rules

func validateNonEmptyStates[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	if len(automaton.states) == 0 {
		return NewValidationError("automaton must have at least one state")
	}
	return nil
}

func validateNonEmptyAlphabet[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	if len(automaton.alphabet) == 0 {
		return NewValidationError("automaton must have at least one symbol in alphabet")
	}
	return nil
}

func validateInitialStateInStates[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	if !automaton.states[automaton.initialState] {
		return NewValidationError(fmt.Sprintf("initial state %v is not in the set of states", automaton.initialState))
	}
	return nil
}

func validateAcceptingStatesInStates[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	for state := range automaton.acceptingStates {
		if !automaton.states[state] {
			return NewValidationError(fmt.Sprintf("accepting state %v is not in the set of states", state))
		}
	}
	return nil
}

func validateTransitionStatesInStates[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	for fromState, transitions := range automaton.transitions {
		if !automaton.states[fromState] {
			return NewValidationError(fmt.Sprintf("transition from state %v is not in the set of states", fromState))
		}
		for _, toState := range transitions {
			if !automaton.states[toState] {
				return NewValidationError(fmt.Sprintf("transition to state %v is not in the set of states", toState))
			}
		}
	}
	return nil
}

func validateTransitionSymbolsInAlphabet[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	for _, transitions := range automaton.transitions {
		for symbol := range transitions {
			if !automaton.alphabet[symbol] {
				return NewValidationError(fmt.Sprintf("transition symbol %v is not in the alphabet", symbol))
			}
		}
	}
	return nil
}

func validateMaxStates[Q State, S Symbol](maxStates int) ValidationRule[Q, S] {
	return func(automaton *FiniteAutomaton[Q, S]) error {
		if len(automaton.states) > maxStates {
			return NewValidationError(fmt.Sprintf(
				"number of states (%d) exceeds maximum allowed (%d)",
				len(automaton.states), maxStates))
		}
		return nil
	}
}

func validateMaxAlphabetSize[Q State, S Symbol](maxSize int) ValidationRule[Q, S] {
	return func(automaton *FiniteAutomaton[Q, S]) error {
		if len(automaton.alphabet) > maxSize {
			return NewValidationError(fmt.Sprintf(
				"alphabet size (%d) exceeds maximum allowed (%d)",
				len(automaton.alphabet), maxSize))
		}
		return nil
	}
}

func validateMaxTransitions[Q State, S Symbol](maxTransitions int) ValidationRule[Q, S] {
	return func(automaton *FiniteAutomaton[Q, S]) error {
		totalTransitions := 0
		for _, transitions := range automaton.transitions {
			totalTransitions += len(transitions)
		}
		if totalTransitions > maxTransitions {
			return NewValidationError(fmt.Sprintf(
				"number of transitions (%d) exceeds maximum allowed (%d)",
				totalTransitions, maxTransitions))
		}
		return nil
	}
}

func validateCompleteTransitions[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	for state := range automaton.states {
		for symbol := range automaton.alphabet {
			if transitions, exists := automaton.transitions[state]; !exists || transitions[symbol] == *new(Q) {
				return NewValidationError(fmt.Sprintf("missing transition from state %v with symbol %v", state, symbol))
			}
		}
	}
	return nil
}

func validateNoUnreachableStates[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	reachable := make(map[Q]bool)
	reachable[automaton.initialState] = true

	// BFS to find all reachable states
	queue := []Q{automaton.initialState}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if transitions, exists := automaton.transitions[current]; exists {
			for _, nextState := range transitions {
				if !reachable[nextState] {
					reachable[nextState] = true
					queue = append(queue, nextState)
				}
			}
		}
	}

	// Check for unreachable states
	for state := range automaton.states {
		if !reachable[state] {
			return NewValidationError(fmt.Sprintf("state %v is unreachable from initial state", state))
		}
	}

	return nil
}

func validateNoDuplicateTransitions[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	// This is automatically handled by map structure, but we can check for consistency
	for fromState, transitions := range automaton.transitions {
		seen := make(map[S]Q)
		for symbol, toState := range transitions {
			if existing, exists := seen[symbol]; exists && existing != toState {
				return NewValidationError(fmt.Sprintf("duplicate transition from state %v with symbol %v", fromState, symbol))
			}
			seen[symbol] = toState
		}
	}
	return nil
}

func validateStateNaming[Q State, S Symbol](automaton *FiniteAutomaton[Q, S]) error {
	// Check for reasonable state naming conventions (string states only)
	for state := range automaton.states {
		stateStr := fmt.Sprintf("%v", state)

		// Check if state is a string type
		if reflect.TypeOf(state).Kind() == reflect.String {
			// Validate string state naming
			if strings.TrimSpace(stateStr) == "" {
				return NewValidationError("state names cannot be empty or whitespace-only")
			}

			if strings.Contains(stateStr, " ") && len(strings.Fields(stateStr)) > 1 {
				return NewValidationError(fmt.Sprintf(
					"state name '%s' contains multiple words (consider using underscores)",
					stateStr))
			}

			if len(stateStr) > MaxStateNameLength {
				return NewValidationError(fmt.Sprintf(
					"state name '%s' is too long (max %d characters)",
					stateStr, MaxStateNameLength))
			}
		}
	}
	return nil
}

// SanitizeInput provides input sanitization for common cases.
func SanitizeInput[S Symbol](input []S) []S {
	if input == nil {
		return []S{}
	}

	// Remove any zero values (depends on symbol type)
	sanitized := make([]S, 0, len(input))
	var zero S

	for _, symbol := range input {
		if symbol != zero {
			sanitized = append(sanitized, symbol)
		}
	}

	return sanitized
}

// ValidateInputSequence validates an input sequence before processing.
func ValidateInputSequence[S Symbol](input []S, alphabet map[S]bool) error {
	if input == nil {
		return NewInvalidInputError(*new(S), 0, "input sequence cannot be nil")
	}

	for i, symbol := range input {
		if !alphabet[symbol] {
			return NewInvalidInputError(symbol, i, fmt.Sprintf("symbol at position %d is not in alphabet", i))
		}
	}

	return nil
}
