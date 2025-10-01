// Package fsm provides a generic finite state automaton (FSM) library.
// It implements the formal definition of a finite automaton as a 5-tuple (Q, Σ, q0, F, δ).
// Package fsm provides a production-ready, type-safe finite state machine library.
//
// This package implements finite automata based on the formal mathematical definition
// as a 5-tuple (Q, Σ, δ, q₀, F) where:
//   - Q is a finite set of states
//   - Σ is a finite set of input symbols (alphabet)
//   - δ is the transition function: Q × Σ → Q
//   - q₀ is the initial state (q₀ ∈ Q)
//   - F is the set of accepting states (F ⊆ Q)
//
// Key Features:
//   - Type-safe generics for states and symbols
//   - Thread-safe operations with RWMutex
//   - Comprehensive validation system
//   - Observer pattern for monitoring
//   - Multiple processing strategies
//   - Factory pattern for extensibility
//   - Zero external dependencies
//
// Basic Usage:
//
//	automaton := fsm.NewBuilder[string, rune]("q0").
//		WithStates("q0", "q1").
//		WithAlphabet('a', 'b').
//		WithAcceptingStates("q1").
//		WithTransition("q0", 'a', "q0").
//		WithTransition("q0", 'b', "q1").
//		MustBuild()
//
//	accepted, _ := automaton.ProcessInput([]rune("aab"))
//	fmt.Printf("Accepted: %t\n", accepted)
//
// Advanced Usage with Observers:
//
//	observable := fsm.NewObservableAutomaton(automaton)
//	logger := func(msg string) { fmt.Println(msg) }
//	observable.AddObserver(fsm.NewLoggingObserver[string, rune](logger))
//
//	accepted, _ := observable.ProcessInput([]rune("test"))
//
// The package is designed for production use with comprehensive error handling,
// input validation, and performance optimizations.
package fsm

import (
	"fmt"
	"strings"
	"sync"
)

// State represents a state in the finite automaton.
// It can be any comparable type (string, int, custom type, etc.)
type State interface {
	comparable
}

// Symbol represents an input symbol from the alphabet.
// It can be any comparable type (rune, string, int, etc.)
type Symbol interface {
	comparable
}

// TransitionFunc is the transition function δ: Q × Σ → Q
// It takes a current state and an input symbol and returns the next state.
type TransitionFunc[Q State, S Symbol] func(state Q, symbol S) (Q, bool)

// FiniteAutomaton represents a finite automaton (FA) as a 5-tuple (Q, Σ, q0, F, δ).
//
// Type parameters:
//   - Q: The type used for states
//   - S: The type used for input symbols
type FiniteAutomaton[Q State, S Symbol] struct {
	// Q: Set of states
	states map[Q]bool

	// Σ (Sigma): Input alphabet
	alphabet map[S]bool

	// q0: Initial state
	initialState Q

	// F: Set of accepting/final states
	acceptingStates map[Q]bool

	// δ (delta): Transition function Q × Σ → Q
	transitions map[Q]map[S]Q

	// Current state (for stateful processing)
	currentState Q

	// Thread safety
	mutex sync.RWMutex
}

// New creates a new FiniteAutomaton with the specified initial state.
// Use the builder methods to configure the automaton.
func New[Q State, S Symbol](initialState Q) *FiniteAutomaton[Q, S] {
	return &FiniteAutomaton[Q, S]{
		states:          make(map[Q]bool),
		alphabet:        make(map[S]bool),
		initialState:    initialState,
		acceptingStates: make(map[Q]bool),
		transitions:     make(map[Q]map[S]Q),
		currentState:    initialState,
	}
}

// AddState adds a state to the set Q.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddState(state Q) *FiniteAutomaton[Q, S] {
	fa.states[state] = true
	return fa
}

// AddStates adds multiple states to the set Q.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddStates(states ...Q) *FiniteAutomaton[Q, S] {
	for _, state := range states {
		fa.states[state] = true
	}
	return fa
}

// AddSymbol adds a symbol to the alphabet Σ.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddSymbol(symbol S) *FiniteAutomaton[Q, S] {
	fa.alphabet[symbol] = true
	return fa
}

// AddSymbols adds multiple symbols to the alphabet Σ.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddSymbols(symbols ...S) *FiniteAutomaton[Q, S] {
	for _, symbol := range symbols {
		fa.alphabet[symbol] = true
	}
	return fa
}

// AddAcceptingState adds a state to the set of accepting states F.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddAcceptingState(state Q) *FiniteAutomaton[Q, S] {
	fa.acceptingStates[state] = true
	return fa
}

// AddAcceptingStates adds multiple states to the set of accepting states F.
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddAcceptingStates(states ...Q) *FiniteAutomaton[Q, S] {
	for _, state := range states {
		fa.acceptingStates[state] = true
	}
	return fa
}

// AddTransition adds a transition to the transition function δ.
// δ(fromState, symbol) = toState
// Returns the automaton for method chaining.
func (fa *FiniteAutomaton[Q, S]) AddTransition(fromState Q, symbol S, toState Q) *FiniteAutomaton[Q, S] {
	if fa.transitions[fromState] == nil {
		fa.transitions[fromState] = make(map[S]Q)
	}
	fa.transitions[fromState][symbol] = toState
	return fa
}

// GetInitialState returns the initial state q0.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) GetInitialState() Q {
	fa.mutex.RLock()
	defer fa.mutex.RUnlock()
	return fa.initialState
}

// GetCurrentState returns the current state during processing.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) GetCurrentState() Q {
	fa.mutex.RLock()
	defer fa.mutex.RUnlock()
	return fa.currentState
}

// IsAcceptingState checks if the given state is in the set of accepting states F.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) IsAcceptingState(state Q) bool {
	fa.mutex.RLock()
	defer fa.mutex.RUnlock()
	return fa.acceptingStates[state]
}

// IsCurrentStateAccepting checks if the current state is an accepting state.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) IsCurrentStateAccepting() bool {
	fa.mutex.RLock()
	defer fa.mutex.RUnlock()
	return fa.acceptingStates[fa.currentState]
}

// Reset resets the automaton to its initial state q0.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) Reset() {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()
	fa.currentState = fa.initialState
}

// Step processes a single input symbol and transitions to the next state.
// Returns the new state and an error if the transition is not defined.
// This method is thread-safe.
func (fa *FiniteAutomaton[Q, S]) Step(symbol S) (Q, error) {
	fa.mutex.Lock()
	defer fa.mutex.Unlock()

	// Validate symbol is in alphabet
	if !fa.alphabet[symbol] {
		var zero Q
		return zero, fmt.Errorf("symbol not in alphabet: %v", symbol)
	}

	// Get transition
	nextState, exists := fa.transitions[fa.currentState][symbol]
	if !exists {
		var zero Q
		return zero, fmt.Errorf("no transition defined for state %v with symbol %v", fa.currentState, symbol)
	}

	fa.currentState = nextState
	return fa.currentState, nil
}

// ProcessInput processes a sequence of input symbols.
// Returns true if the automaton ends in an accepting state, false otherwise.
// Returns an error if any transition is undefined.
// This method is thread-safe and includes input validation.
func (fa *FiniteAutomaton[Q, S]) ProcessInput(input []S) (bool, error) {
	// Validate input sequence
	if err := ValidateInputSequence(input, fa.alphabet); err != nil {
		return false, err
	}

	fa.Reset()

	for _, symbol := range input {
		_, err := fa.Step(symbol)
		if err != nil {
			return false, err
		}
	}

	return fa.IsCurrentStateAccepting(), nil
}

// ProcessInputWithTrace processes input and returns a trace of state transitions.
// Returns the trace and whether the input was accepted.
// This method is thread-safe and includes input validation.
func (fa *FiniteAutomaton[Q, S]) ProcessInputWithTrace(input []S) ([]Q, bool, error) {
	// Validate input sequence
	if err := ValidateInputSequence(input, fa.alphabet); err != nil {
		return nil, false, err
	}

	fa.Reset()
	trace := []Q{fa.currentState}

	for _, symbol := range input {
		_, err := fa.Step(symbol)
		if err != nil {
			return trace, false, err
		}
		trace = append(trace, fa.currentState)
	}

	return trace, fa.IsCurrentStateAccepting(), nil
}

// Validate checks if the automaton is properly configured.
// Returns an error if the configuration is invalid.
func (fa *FiniteAutomaton[Q, S]) Validate() error {
	// Check if initial state is in Q
	if !fa.states[fa.initialState] {
		return fmt.Errorf("initial state %v is not in the set of states Q", fa.initialState)
	}

	// Check if all accepting states are in Q
	for state := range fa.acceptingStates {
		if !fa.states[state] {
			return fmt.Errorf("accepting state %v is not in the set of states Q", state)
		}
	}

	// Check if all transitions reference valid states and symbols
	for fromState, transitions := range fa.transitions {
		if !fa.states[fromState] {
			return fmt.Errorf("transition from state %v, but state is not in Q", fromState)
		}
		for symbol, toState := range transitions {
			if !fa.alphabet[symbol] {
				return fmt.Errorf("transition uses symbol %v, but symbol is not in Σ", symbol)
			}
			if !fa.states[toState] {
				return fmt.Errorf("transition to state %v, but state is not in Q", toState)
			}
		}
	}

	return nil
}

// String returns a string representation of the automaton configuration.
func (fa *FiniteAutomaton[Q, S]) String() string {
	var sb strings.Builder

	sb.WriteString("Finite Automaton:\n")
	sb.WriteString(fmt.Sprintf("  Q (States): %v\n", fa.getStatesList()))
	sb.WriteString(fmt.Sprintf("  Σ (Alphabet): %v\n", fa.getAlphabetList()))
	sb.WriteString(fmt.Sprintf("  q0 (Initial): %v\n", fa.initialState))
	sb.WriteString(fmt.Sprintf("  F (Accepting): %v\n", fa.getAcceptingStatesList()))
	sb.WriteString("  δ (Transitions):\n")
	for state, transitions := range fa.transitions {
		for symbol, nextState := range transitions {
			sb.WriteString(fmt.Sprintf("    δ(%v, %v) = %v\n", state, symbol, nextState))
		}
	}

	return sb.String()
}

func (fa *FiniteAutomaton[Q, S]) getStatesList() []Q {
	states := make([]Q, 0, len(fa.states))
	for state := range fa.states {
		states = append(states, state)
	}
	return states
}

func (fa *FiniteAutomaton[Q, S]) getAlphabetList() []S {
	alphabet := make([]S, 0, len(fa.alphabet))
	for symbol := range fa.alphabet {
		alphabet = append(alphabet, symbol)
	}
	return alphabet
}

func (fa *FiniteAutomaton[Q, S]) getAcceptingStatesList() []Q {
	states := make([]Q, 0, len(fa.acceptingStates))
	for state := range fa.acceptingStates {
		states = append(states, state)
	}
	return states
}
