package fsm

import (
	"strings"
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

// TestFiniteAutomaton_EdgeCases tests various edge cases and boundary conditions
func TestFiniteAutomaton_EdgeCases(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddAcceptingState("q0")

		accepted, err := fa.ProcessInput([]rune{})
		if err != nil {
			t.Fatalf("ProcessInput([]) returned error: %v", err)
		}
		if !accepted {
			t.Error("Empty input should be accepted when initial state is accepting")
		}
	})

	t.Run("empty input non-accepting initial state", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddAcceptingState("q1")

		accepted, err := fa.ProcessInput([]rune{})
		if err != nil {
			t.Fatalf("ProcessInput([]) returned error: %v", err)
		}
		if accepted {
			t.Error("Empty input should not be accepted when initial state is not accepting")
		}
	})

	t.Run("single symbol alphabet", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddAcceptingState("q1").
			AddTransition("q0", 'a', "q1")

		accepted, err := fa.ProcessInput([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInput('a') returned error: %v", err)
		}
		if !accepted {
			t.Error("Single symbol should be accepted")
		}
	})

	t.Run("self-loop transitions", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddSymbol('a').
			AddAcceptingState("q0").
			AddTransition("q0", 'a', "q0")

		accepted, err := fa.ProcessInput([]rune("aaaa"))
		if err != nil {
			t.Fatalf("ProcessInput('aaaa') returned error: %v", err)
		}
		if !accepted {
			t.Error("Self-loop should accept repeated symbols")
		}
	})

	t.Run("no accepting states", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddTransition("q0", 'a', "q1")

		accepted, err := fa.ProcessInput([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInput('a') returned error: %v", err)
		}
		if accepted {
			t.Error("Should not accept when no accepting states exist")
		}
	})
}

// TestFiniteAutomaton_StringRepresentation tests the String() method
func TestFiniteAutomaton_StringRepresentation(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbols('a', 'b').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q1")

	str := fa.String()

	// Check that all components are present in the string representation
	if !strings.Contains(str, "Finite Automaton") {
		t.Error("String representation should contain 'Finite Automaton'")
	}
	if !strings.Contains(str, "Q (States)") {
		t.Error("String representation should contain states section")
	}
	if !strings.Contains(str, "Σ (Alphabet)") {
		t.Error("String representation should contain alphabet section")
	}
	if !strings.Contains(str, "q0 (Initial)") {
		t.Error("String representation should contain initial state")
	}
	if !strings.Contains(str, "F (Accepting)") {
		t.Error("String representation should contain accepting states")
	}
	if !strings.Contains(str, "δ (Transitions)") {
		t.Error("String representation should contain transitions")
	}
}

// TestFiniteAutomaton_ConcurrentAccess tests thread safety (basic check)
func TestFiniteAutomaton_ConcurrentAccess(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbol('a').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q1")

	// Test concurrent reads (should be safe)
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_ = fa.GetInitialState()
			_ = fa.IsAcceptingState("q1")
			_ = fa.String()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestFiniteAutomaton_LargeAlphabet tests performance with large alphabets
func TestFiniteAutomaton_LargeAlphabet(t *testing.T) {
	fa := New[string, rune]("q0").AddState("q0").AddAcceptingState("q0")

	// Add large alphabet (all printable ASCII)
	for r := rune(32); r <= rune(126); r++ {
		fa.AddSymbol(r)
		fa.AddTransition("q0", r, "q0")
	}

	// Test with various inputs
	testInputs := []string{"hello", "world", "123", "!@#$%"}
	for _, input := range testInputs {
		accepted, err := fa.ProcessInput([]rune(input))
		if err != nil {
			t.Fatalf("ProcessInput(%q) returned error: %v", input, err)
		}
		if !accepted {
			t.Errorf("Input %q should be accepted", input)
		}
	}
}

// TestFiniteAutomaton_DeepStateChain tests performance with many states
func TestFiniteAutomaton_DeepStateChain(t *testing.T) {
	const numStates = 100
	fa := New[int, rune](0)

	// Create a chain of states 0 -> 1 -> 2 -> ... -> numStates-1
	for i := 0; i < numStates; i++ {
		fa.AddState(i)
		if i < numStates-1 {
			fa.AddTransition(i, 'a', i+1)
		}
	}
	fa.AddSymbol('a').AddAcceptingState(numStates - 1)

	// Create input that traverses the entire chain
	input := make([]rune, numStates-1)
	for i := range input {
		input[i] = 'a'
	}

	accepted, err := fa.ProcessInput(input)
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}
	if !accepted {
		t.Error("Long chain should be accepted")
	}
}

// TestFiniteAutomaton_ErrorMessages tests that error messages are descriptive
func TestFiniteAutomaton_ErrorMessages(t *testing.T) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbol('a')

	t.Run("invalid symbol error message", func(t *testing.T) {
		_, err := fa.Step('b')
		if err == nil {
			t.Fatal("Expected error for invalid symbol")
		}
		if !strings.Contains(err.Error(), "symbol not in alphabet") {
			t.Errorf("Error message should mention 'symbol not in alphabet', got: %v", err)
		}
		if !strings.Contains(err.Error(), "b") {
			t.Errorf("Error message should include the invalid symbol 'b', got: %v", err)
		}
	})

	t.Run("undefined transition error message", func(t *testing.T) {
		_, err := fa.Step('a')
		if err == nil {
			t.Fatal("Expected error for undefined transition")
		}
		if !strings.Contains(err.Error(), "no transition defined") {
			t.Errorf("Error message should mention 'no transition defined', got: %v", err)
		}
		if !strings.Contains(err.Error(), "q0") && !strings.Contains(err.Error(), "a") {
			t.Errorf("Error message should include state and symbol, got: %v", err)
		}
	})
}

// TestFiniteAutomaton_ValidationEdgeCases tests additional validation scenarios
func TestFiniteAutomaton_ValidationEdgeCases(t *testing.T) {
	t.Run("empty states set", func(t *testing.T) {
		fa := New[string, rune]("q0")
		// Don't add any states
		err := fa.Validate()
		if err == nil {
			t.Error("Should fail validation when initial state not in Q")
		}
	})

	t.Run("empty alphabet", func(t *testing.T) {
		fa := New[string, rune]("q0").AddState("q0")
		// Don't add any symbols
		err := fa.Validate()
		if err != nil {
			t.Errorf("Empty alphabet should be valid: %v", err)
		}
	})

	t.Run("transitions without alphabet symbols", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddTransition("q0", 'a', "q1")
		// Symbol 'a' not added to alphabet
		err := fa.Validate()
		if err == nil {
			t.Error("Should fail validation when transition uses symbol not in alphabet")
		}
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q1").                // q0 not in states
			AddAcceptingState("q2").       // q2 not in states
			AddTransition("q3", 'a', "q4") // q3, q4 not in states, 'a' not in alphabet
		err := fa.Validate()
		if err == nil {
			t.Error("Should fail validation with multiple errors")
		}
	})
}

// TestFiniteAutomaton_StateManagement tests state-related operations
func TestFiniteAutomaton_StateManagement(t *testing.T) {
	t.Run("duplicate state addition", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddState("q0").             // Add same state twice
			AddStates("q0", "q1", "q0") // Add duplicates in batch

		if len(fa.states) != 2 {
			t.Errorf("Expected 2 unique states, got %d", len(fa.states))
		}
	})

	t.Run("duplicate accepting state addition", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddAcceptingState("q0").
			AddAcceptingState("q0").             // Add same accepting state twice
			AddAcceptingStates("q0", "q1", "q0") // Add duplicates in batch

		if len(fa.acceptingStates) != 2 {
			t.Errorf("Expected 2 unique accepting states, got %d", len(fa.acceptingStates))
		}
	})

	t.Run("duplicate symbol addition", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddSymbol('a').
			AddSymbol('a').           // Add same symbol twice
			AddSymbols('a', 'b', 'a') // Add duplicates in batch

		if len(fa.alphabet) != 2 {
			t.Errorf("Expected 2 unique symbols, got %d", len(fa.alphabet))
		}
	})

	t.Run("transition overwriting", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1", "q2").
			AddSymbol('a').
			AddTransition("q0", 'a', "q1").
			AddTransition("q0", 'a', "q2") // Overwrite previous transition

		if fa.transitions["q0"]['a'] != "q2" {
			t.Errorf("Transition should be overwritten to q2, got %v", fa.transitions["q0"]['a'])
		}
	})
}

// TestFiniteAutomaton_ProcessInputEdgeCases tests edge cases in input processing
func TestFiniteAutomaton_ProcessInputEdgeCases(t *testing.T) {
	t.Run("process input resets state", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddAcceptingState("q1").
			AddTransition("q0", 'a', "q1")

		// Move to q1 manually
		fa.currentState = "q1"

		// ProcessInput should reset to q0 first
		accepted, err := fa.ProcessInput([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInput returned error: %v", err)
		}
		if !accepted {
			t.Error("ProcessInput should reset and then process from initial state")
		}
	})

	t.Run("trace includes initial state", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddStates("q0", "q1").
			AddSymbol('a').
			AddAcceptingState("q1").
			AddTransition("q0", 'a', "q1")

		trace, accepted, err := fa.ProcessInputWithTrace([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInputWithTrace returned error: %v", err)
		}
		if !accepted {
			t.Error("Input should be accepted")
		}
		if len(trace) != 2 {
			t.Fatalf("Trace should have 2 states, got %d", len(trace))
		}
		if trace[0] != "q0" {
			t.Errorf("First state in trace should be q0, got %v", trace[0])
		}
		if trace[1] != "q1" {
			t.Errorf("Second state in trace should be q1, got %v", trace[1])
		}
	})

	t.Run("trace with empty input", func(t *testing.T) {
		fa := New[string, rune]("q0").
			AddState("q0").
			AddAcceptingState("q0")

		trace, accepted, err := fa.ProcessInputWithTrace([]rune{})
		if err != nil {
			t.Fatalf("ProcessInputWithTrace returned error: %v", err)
		}
		if !accepted {
			t.Error("Empty input should be accepted")
		}
		if len(trace) != 1 {
			t.Fatalf("Trace should have 1 state for empty input, got %d", len(trace))
		}
		if trace[0] != "q0" {
			t.Errorf("Trace should contain initial state q0, got %v", trace[0])
		}
	})
}

// BenchmarkFiniteAutomaton_Step benchmarks single step operations
func BenchmarkFiniteAutomaton_Step(b *testing.B) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbol('a').
		AddTransition("q0", 'a', "q1").
		AddTransition("q1", 'a', "q0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fa.Reset()
		_, _ = fa.Step('a')
	}
}

// BenchmarkFiniteAutomaton_ProcessInput benchmarks input processing
func BenchmarkFiniteAutomaton_ProcessInput(b *testing.B) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbols('a', 'b').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q0").
		AddTransition("q0", 'b', "q1").
		AddTransition("q1", 'a', "q0").
		AddTransition("q1", 'b', "q1")

	input := []rune("aaabaaabaaab")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = fa.ProcessInput(input)
	}
}

// BenchmarkFiniteAutomaton_ProcessInputWithTrace benchmarks trace generation
func BenchmarkFiniteAutomaton_ProcessInputWithTrace(b *testing.B) {
	fa := New[string, rune]("q0").
		AddStates("q0", "q1").
		AddSymbols('a', 'b').
		AddAcceptingState("q1").
		AddTransition("q0", 'a', "q0").
		AddTransition("q0", 'b', "q1").
		AddTransition("q1", 'a', "q0").
		AddTransition("q1", 'b', "q1")

	input := []rune("aaabaaabaaab")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = fa.ProcessInputWithTrace(input)
	}
}

// BenchmarkFiniteAutomaton_LargeInput benchmarks with large inputs
func BenchmarkFiniteAutomaton_LargeInput(b *testing.B) {
	fa := New[string, rune]("q0").
		AddState("q0").
		AddSymbols('0', '1').
		AddAcceptingState("q0").
		AddTransition("q0", '0', "q0").
		AddTransition("q0", '1', "q0")

	// Create large input
	input := make([]rune, 10000)
	for i := range input {
		if i%2 == 0 {
			input[i] = '0'
		} else {
			input[i] = '1'
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = fa.ProcessInput(input)
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
