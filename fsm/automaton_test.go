package fsm

import (
	"testing"
)

// TestFiniteAutomaton_BasicConstruction tests basic automaton construction
func TestFiniteAutomaton_BasicConstruction(t *testing.T) {
	fa := New[string, rune]("q0")

	if fa.GetInitialState() != "q0" {
		t.Errorf("Initial state = %v, want q0", fa.GetInitialState())
	}

	if fa.GetCurrentState() != "q0" {
		t.Errorf("Current state = %v, want q0", fa.GetCurrentState())
	}
}

// TestFiniteAutomaton_AddStates tests adding states
func TestFiniteAutomaton_AddStates(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1", "q2")

	if !fa.states["q0"] {
		t.Error("State q0 not added")
	}
	if !fa.states["q1"] {
		t.Error("State q1 not added")
	}
	if !fa.states["q2"] {
		t.Error("State q2 not added")
	}
}

// TestFiniteAutomaton_AddSymbols tests adding alphabet symbols
func TestFiniteAutomaton_AddSymbols(t *testing.T) {
	fa := New[string, rune]("q0").
		AddSymbols('0', '1')

	if !fa.alphabet['0'] {
		t.Error("Symbol '0' not added")
	}
	if !fa.alphabet['1'] {
		t.Error("Symbol '1' not added")
	}
}

// TestFiniteAutomaton_AddTransitions tests adding transitions
func TestFiniteAutomaton_AddTransitions(t *testing.T) {
	fa := New[string, rune]("q0").
		AddTransition("q0", '0', "q1").
		AddTransition("q0", '1', "q2")

	if fa.transitions["q0"]['0'] != "q1" {
		t.Errorf("Transition δ(q0, '0') = %v, want q1", fa.transitions["q0"]['0'])
	}
	if fa.transitions["q0"]['1'] != "q2" {
		t.Errorf("Transition δ(q0, '1') = %v, want q2", fa.transitions["q0"]['1'])
	}
}

// TestFiniteAutomaton_AcceptingStates tests accepting state management
func TestFiniteAutomaton_AcceptingStates(t *testing.T) {
	fa := New[string, rune]("q0").
		AddAcceptingStates("q1", "q2")

	if fa.IsAcceptingState("q0") {
		t.Error("q0 should not be accepting")
	}
	if !fa.IsAcceptingState("q1") {
		t.Error("q1 should be accepting")
	}
	if !fa.IsAcceptingState("q2") {
		t.Error("q2 should be accepting")
	}
}

// TestFiniteAutomaton_Reset tests state reset
func TestFiniteAutomaton_Reset(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbol('a').
		AddTransition("q0", 'a', "q1")

	// Move to q1
	fa.currentState = "q1"

	// Reset should go back to q0
	fa.Reset()

	if fa.GetCurrentState() != "q0" {
		t.Errorf("After reset, current state = %v, want q0", fa.GetCurrentState())
	}
}

// TestFiniteAutomaton_Step tests single step transitions
func TestFiniteAutomaton_Step(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbol('a').
		AddTransition("q0", 'a', "q1")

	state, err := fa.Step('a')
	if err != nil {
		t.Fatalf("Step returned error: %v", err)
	}

	if state != "q1" {
		t.Errorf("After step, state = %v, want q1", state)
	}

	if fa.GetCurrentState() != "q1" {
		t.Errorf("Current state = %v, want q1", fa.GetCurrentState())
	}
}

// TestFiniteAutomaton_StepInvalidSymbol tests error handling for invalid symbols
func TestFiniteAutomaton_StepInvalidSymbol(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0").
		AddSymbol('a')

	_, err := fa.Step('b')
	if err == nil {
		t.Error("Expected error for invalid symbol, got nil")
	}
}

// TestFiniteAutomaton_StepUndefinedTransition tests error for undefined transitions
func TestFiniteAutomaton_StepUndefinedTransition(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0").
		AddSymbol('a')

	_, err := fa.Step('a')
	if err == nil {
		t.Error("Expected error for undefined transition, got nil")
	}
}

// TestFiniteAutomaton_ProcessInput tests processing a sequence of inputs
func TestFiniteAutomaton_ProcessInput(t *testing.T) {
	// Simple automaton that accepts strings ending in 'b'
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbols('a', 'b').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q0").
		AddTransition("q0", 'b', "q1").
		AddTransition("q1", 'a', "q0").
		AddTransition("q1", 'b', "q1")

	tests := []struct {
		input    string
		accepted bool
	}{
		{"b", true},
		{"ab", true},
		{"aab", true},
		{"abb", true},
		{"a", false},
		{"aa", false},
		{"aba", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			accepted, err := fa.ProcessInput([]rune(tt.input))
			if err != nil {
				t.Fatalf("ProcessInput(%q) returned error: %v", tt.input, err)
			}
			if accepted != tt.accepted {
				t.Errorf("ProcessInput(%q) = %v, want %v", tt.input, accepted, tt.accepted)
			}
		})
	}
}

// TestFiniteAutomaton_ProcessInputWithTrace tests trace generation
func TestFiniteAutomaton_ProcessInputWithTrace(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbols('a', 'b').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q0").
		AddTransition("q0", 'b', "q1")

	trace, accepted, err := fa.ProcessInputWithTrace([]rune("aab"))
	if err != nil {
		t.Fatalf("ProcessInputWithTrace returned error: %v", err)
	}

	expectedTrace := []string{"q0", "q0", "q0", "q1"}
	if len(trace) != len(expectedTrace) {
		t.Fatalf("Trace length = %d, want %d", len(trace), len(expectedTrace))
	}

	for i, state := range trace {
		if state != expectedTrace[i] {
			t.Errorf("Trace[%d] = %v, want %v", i, state, expectedTrace[i])
		}
	}

	if !accepted {
		t.Error("Input should be accepted")
	}
}

// TestFiniteAutomaton_Validate tests automaton validation
func TestFiniteAutomaton_Validate(t *testing.T) {
	t.Run("valid automaton", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddAcceptingState("q1").
			AddTransition("q0", 'a', "q1")

		if err := fa.Validate(); err != nil {
			t.Errorf("Validate() returned error for valid automaton: %v", err)
		}
	})

	t.Run("initial state not in Q", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q1")

		if err := fa.Validate(); err == nil {
			t.Error("Expected error for initial state not in Q")
		}
	})

	t.Run("accepting state not in Q", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddAcceptingState("q1")

		if err := fa.Validate(); err == nil {
			t.Error("Expected error for accepting state not in Q")
		}
	})

	t.Run("transition from state not in Q", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddTransition("q2", 'a', "q1")

		if err := fa.Validate(); err == nil {
			t.Error("Expected error for transition from state not in Q")
		}
	})

	t.Run("transition with symbol not in Σ", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddTransition("q0", 'b', "q1")

		if err := fa.Validate(); err == nil {
			t.Error("Expected error for transition with symbol not in Σ")
		}
	})

	t.Run("transition to state not in Q", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddSymbol('a').
			AddTransition("q0", 'a', "q1")

		if err := fa.Validate(); err == nil {
			t.Error("Expected error for transition to state not in Q")
		}
	})
}

