package fsm

import (
	"strings"
	"testing"
)

// TestBuilder_BasicConstruction tests basic builder usage
func TestBuilder_BasicConstruction(t *testing.T) {
	fa, err := NewBuilder[string, rune]("q0").
		WithStates("q0", "q1").
		WithAlphabet('a', 'b').
		WithAcceptingStates("q1").
		WithTransition("q0", 'a', "q1").
		Build()

	if err != nil {
		t.Fatalf("Build() returned error: %v", err)
	}

	if fa.GetInitialState() != "q0" {
		t.Errorf("Initial state = %v, want q0", fa.GetInitialState())
	}
}

// TestBuilder_WithTransitions tests bulk transition addition
func TestBuilder_WithTransitions(t *testing.T) {
	fa := NewBuilder[string, rune]("q0").
		WithStates("q0", "q1", "q2").
		WithAlphabet('a', 'b').
		WithAcceptingStates("q2").
		WithTransitions(
			T("q0", 'a', "q1"),
			T("q1", 'b', "q2"),
		).
		MustBuild()

	// Test transitions by processing input
	accepted, err := fa.ProcessInput([]rune("ab"))
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}
	if !accepted {
		t.Error("Transitions not added correctly - input 'ab' should be accepted")
	}
}

// TestBuilder_MustBuild tests MustBuild panic behavior
func TestBuilder_MustBuild(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBuild should panic on invalid automaton")
		}
	}()

	// This should panic because q1 is not in Q
	NewBuilder[string, rune]("q0").
		WithStates("q0").
		WithAcceptingStates("q1").
		MustBuild()
}

// TestBuilder_ChainedCalls tests method chaining
func TestBuilder_ChainedCalls(t *testing.T) {
	fa := NewBuilder[string, rune]("start").
		WithStates("start", "middle", "end").
		WithAlphabet('x', 'y', 'z').
		WithAcceptingStates("end").
		WithTransition("start", 'x', "middle").
		WithTransition("middle", 'y', "end").
		MustBuild()

	if fa.GetInitialState() != "start" {
		t.Errorf("Initial state = %v, want start", fa.GetInitialState())
	}

	if !fa.IsAcceptingState("end") {
		t.Error("end should be an accepting state")
	}
}

// TestBuilder_IntegerStates tests using integers as states
func TestBuilder_IntegerStates(t *testing.T) {
	fa := NewBuilder[int, rune](0).
		WithStates(0, 1, 2).
		WithAlphabet('a', 'b').
		WithAcceptingStates(2).
		WithTransitions(
			T(0, 'a', 1),
			T(1, 'b', 2),
		).
		MustBuild()

	accepted, err := fa.ProcessInput([]rune("ab"))
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}

	if !accepted {
		t.Error("Input 'ab' should be accepted")
	}
}

// TestBuilder_CustomTypes tests using custom types
func TestBuilder_CustomTypes(t *testing.T) {
	type State string
	type Symbol int

	const (
		StateA State = "A"
		StateB State = "B"
	)

	const (
		Symbol0 Symbol = 0
		Symbol1 Symbol = 1
	)

	fa := NewBuilder[State, Symbol](StateA).
		WithStates(StateA, StateB).
		WithAlphabet(Symbol0, Symbol1).
		WithAcceptingStates(StateB).
		WithTransition(StateA, Symbol0, StateB).
		MustBuild()

	accepted, err := fa.ProcessInput([]Symbol{Symbol0})
	if err != nil {
		t.Fatalf("ProcessInput returned error: %v", err)
	}

	if !accepted {
		t.Error("Input should be accepted")
	}
}

// TestBuilder_EdgeCases tests edge cases in builder usage
func TestBuilder_EdgeCases(t *testing.T) {
	t.Run("builder with initial state auto-added", func(t *testing.T) {
		// Initial state should be automatically added by WithStates
		fa := NewBuilder[string, rune]("q0").
			WithStates("q0"). // This ensures initial state is in Q
			WithAlphabet('a').
			WithAcceptingStates("q0").
			WithTransition("q0", 'a', "q0").
			MustBuild()

		// Test that initial state is properly configured by checking it's the initial state
		if fa.GetInitialState() != "q0" {
			t.Error("Initial state should be q0")
		}
	})

	t.Run("builder with empty alphabet", func(t *testing.T) {
		// Empty alphabet should now cause validation error
		_, err := NewBuilder[string, rune]("q0").
			WithStates("q0").
			WithAcceptingStates("q0").
			Build()

		if err == nil {
			t.Error("Expected validation error for empty alphabet")
		}

		if !strings.Contains(err.Error(), "alphabet") {
			t.Errorf("Expected alphabet validation error, got: %v", err)
		}
	})

	t.Run("builder with no accepting states", func(t *testing.T) {
		fa := NewBuilder[string, rune]("q0").
			WithStates("q0", "q1").
			WithAlphabet('a').
			WithTransition("q0", 'a', "q1").
			MustBuild()

		accepted, err := fa.ProcessInput([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInput('a') returned error: %v", err)
		}
		if accepted {
			t.Error("Should not accept when no accepting states")
		}
	})

	t.Run("builder with duplicate transitions", func(t *testing.T) {
		fa := NewBuilder[string, rune]("q0").
			WithStates("q0", "q1", "q2").
			WithAlphabet('a').
			WithTransition("q0", 'a', "q1").
			WithTransition("q0", 'a', "q2"). // Overwrite
			WithAcceptingStates("q2").
			MustBuild()

		accepted, err := fa.ProcessInput([]rune("a"))
		if err != nil {
			t.Fatalf("ProcessInput('a') returned error: %v", err)
		}
		if !accepted {
			t.Error("Should use the last transition definition")
		}
	})
}

// TestBuilder_ValidationErrors tests various validation error scenarios
func TestBuilder_ValidationErrors(t *testing.T) {
	t.Run("accepting state not in states", func(t *testing.T) {
		_, err := NewBuilder[string, rune]("q0").
			WithStates("q0").
			WithAcceptingStates("q1"). // q1 not in states
			Build()

		if err == nil {
			t.Error("Expected validation error for accepting state not in Q")
		}
	})

	t.Run("transition from state not in states", func(t *testing.T) {
		_, err := NewBuilder[string, rune]("q0").
			WithStates("q0").
			WithAlphabet('a').
			WithTransition("q1", 'a', "q0"). // q1 not in states
			Build()

		if err == nil {
			t.Error("Expected validation error for transition from state not in Q")
		}
	})

	t.Run("transition to state not in states", func(t *testing.T) {
		_, err := NewBuilder[string, rune]("q0").
			WithStates("q0").
			WithAlphabet('a').
			WithTransition("q0", 'a', "q1"). // q1 not in states
			Build()

		if err == nil {
			t.Error("Expected validation error for transition to state not in Q")
		}
	})

	t.Run("transition with symbol not in alphabet", func(t *testing.T) {
		_, err := NewBuilder[string, rune]("q0").
			WithStates("q0", "q1").
			WithAlphabet('a').
			WithTransition("q0", 'b', "q1"). // 'b' not in alphabet
			Build()

		if err == nil {
			t.Error("Expected validation error for transition with symbol not in Σ")
		}
	})
}

// TestBuilder_MethodChaining tests that all methods return the builder for chaining
func TestBuilder_MethodChaining(t *testing.T) {
	builder := NewBuilder[string, rune]("q0")

	// Test that all methods return the same builder instance
	b1 := builder.WithStates("q0", "q1")
	b2 := b1.WithAlphabet('a', 'b')
	b3 := b2.WithAcceptingStates("q1")
	b4 := b3.WithTransition("q0", 'a', "q1")
	b5 := b4.WithTransitions(T("q1", 'b', "q0"))

	// All should be the same instance
	if b1 != builder || b2 != builder || b3 != builder || b4 != builder || b5 != builder {
		t.Error("All builder methods should return the same builder instance for chaining")
	}
}

// TestBuilder_ComplexAutomaton tests building a more complex automaton
func TestBuilder_ComplexAutomaton(t *testing.T) {
	// Build an automaton that accepts binary strings with even number of 1s
	fa := NewBuilder[string, rune]("even").
		WithStates("even", "odd").
		WithAlphabet('0', '1').
		WithAcceptingStates("even").
		WithTransitions(
			T("even", '0', "even"),
			T("even", '1', "odd"),
			T("odd", '0', "odd"),
			T("odd", '1', "even"),
		).
		MustBuild()

	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},       // 0 ones (even)
		{"0", true},      // 0 ones (even)
		{"1", false},     // 1 one (odd)
		{"11", true},     // 2 ones (even)
		{"101", true},    // 2 ones (even)
		{"1010", true},   // 2 ones (even)
		{"1111", true},   // 4 ones (even)
		{"00100", false}, // 1 one (odd)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			accepted, err := fa.ProcessInput([]rune(tt.input))
			if err != nil {
				t.Fatalf("ProcessInput(%q) returned error: %v", tt.input, err)
			}
			if accepted != tt.expected {
				t.Errorf("ProcessInput(%q) = %v, want %v", tt.input, accepted, tt.expected)
			}
		})
	}
}

// BenchmarkBuilder_Construction benchmarks automaton construction
func BenchmarkBuilder_Construction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBuilder[string, rune]("q0").
			WithStates("q0", "q1", "q2").
			WithAlphabet('a', 'b', 'c').
			WithAcceptingStates("q2").
			WithTransitions(
				T("q0", 'a', "q1"),
				T("q1", 'b', "q2"),
				T("q2", 'c', "q0"),
			).
			MustBuild()
	}
}

// BenchmarkBuilder_LargeAutomaton benchmarks building large automata
func BenchmarkBuilder_LargeAutomaton(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewBuilder[int, rune](0)

		// Add many states
		states := make([]int, 100)
		for j := range states {
			states[j] = j
		}
		builder.WithStates(states...)

		// Add alphabet
		builder.WithAlphabet('a', 'b', 'c')

		// Add accepting states
		builder.WithAcceptingStates(99)

		// Add transitions
		transitions := make([]Transition[int, rune], 0, 300)
		for j := 0; j < 100; j++ {
			transitions = append(transitions,
				T(j, 'a', (j+1)%100),
				T(j, 'b', (j+2)%100),
				T(j, 'c', (j+3)%100),
			)
		}
		builder.WithTransitions(transitions...)

		_ = builder.MustBuild()
	}
}
