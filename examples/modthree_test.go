package examples

import (
	"testing"
)

// TestModThree_SpecificationExample tests the example from the specification
func TestModThree_SpecificationExample(t *testing.T) {
	result, err := ModThree("110")
	if err != nil {
		t.Fatalf("ModThree(\"110\") returned error: %v", err)
	}

	if result != 0 {
		t.Errorf("ModThree(\"110\") = %d, want 0", result)
	}
}

// TestModThree_ValidInputs tests various valid binary inputs
func TestModThree_ValidInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0", 0},       // 0 % 3 = 0
		{"1", 1},       // 1 % 3 = 1
		{"10", 2},      // 2 % 3 = 2
		{"11", 0},      // 3 % 3 = 0
		{"100", 1},     // 4 % 3 = 1
		{"101", 2},     // 5 % 3 = 2
		{"110", 0},     // 6 % 3 = 0
		{"111", 1},     // 7 % 3 = 1
		{"1000", 2},    // 8 % 3 = 2
		{"1001", 0},    // 9 % 3 = 0
		{"1010", 1},    // 10 % 3 = 1
		{"1111", 0},    // 15 % 3 = 0
		{"10000", 1},   // 16 % 3 = 1
		{"101010", 0},  // 42 % 3 = 0
		{"111111", 0},  // 63 % 3 = 0
		{"1000000", 1}, // 64 % 3 = 1
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ModThree(tt.input)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", tt.input, err)
			}

			if result != tt.expected {
				t.Errorf("ModThree(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestModThree_InvalidInputs tests error handling
func TestModThree_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"contains letter", "10a10"},
		{"contains space", "10 10"},
		{"contains 2", "1012"},
		{"only letters", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ModThree(tt.input)
			if err == nil {
				t.Errorf("ModThree(%q) expected error, got nil", tt.input)
			}
		})
	}
}

// TestModThreeWithTrace_StateTransitions tests state transition traces
func TestModThreeWithTrace_StateTransitions(t *testing.T) {
	tests := []struct {
		input          string
		expectedStates []ModThreeState
		expectedResult int
	}{
		{
			input:          "110",
			expectedStates: []ModThreeState{S0, S1, S0, S0},
			expectedResult: 0,
		},
		{
			input:          "101",
			expectedStates: []ModThreeState{S0, S1, S2, S2},
			expectedResult: 2,
		},
		{
			input:          "111",
			expectedStates: []ModThreeState{S0, S1, S0, S1},
			expectedResult: 1,
		},
		{
			input:          "0",
			expectedStates: []ModThreeState{S0, S0},
			expectedResult: 0,
		},
		{
			input:          "1",
			expectedStates: []ModThreeState{S0, S1},
			expectedResult: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, trace, err := ModThreeWithTrace(tt.input)
			if err != nil {
				t.Fatalf("ModThreeWithTrace(%q) returned error: %v", tt.input, err)
			}

			if result != tt.expectedResult {
				t.Errorf("ModThreeWithTrace(%q) result = %d, want %d", tt.input, result, tt.expectedResult)
			}

			if len(trace) != len(tt.expectedStates) {
				t.Fatalf("Trace length = %d, want %d", len(trace), len(tt.expectedStates))
			}

			for i, state := range trace {
				if state != tt.expectedStates[i] {
					t.Errorf("Trace[%d] = %v, want %v", i, state, tt.expectedStates[i])
				}
			}
		})
	}
}

// TestNewModThreeAutomaton_Configuration tests the automaton configuration
func TestNewModThreeAutomaton_Configuration(t *testing.T) {
	fa := NewModThreeAutomaton()

	// Test initial state
	if fa.GetInitialState() != S0 {
		t.Errorf("Initial state = %v, want S0", fa.GetInitialState())
	}

	// Test all states are accepting
	if !fa.IsAcceptingState(S0) {
		t.Error("S0 should be accepting")
	}
	if !fa.IsAcceptingState(S1) {
		t.Error("S1 should be accepting")
	}
	if !fa.IsAcceptingState(S2) {
		t.Error("S2 should be accepting")
	}

	// Test validation
	if err := fa.Validate(); err != nil {
		t.Errorf("Automaton validation failed: %v", err)
	}
}

// TestModThree_LongInputs tests with longer binary strings
func TestModThree_LongInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"11111111", 0},     // 255 % 3 = 0
		{"100000000", 1},    // 256 % 3 = 1
		{"1111111111", 0},   // 1023 % 3 = 0
		{"10000000000", 1},  // 1024 % 3 = 1
		{"101010101010", 0}, // 2730 % 3 = 0
		{"111111111111", 0}, // 4095 % 3 = 0
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ModThree(tt.input)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", tt.input, err)
			}

			if result != tt.expected {
				t.Errorf("ModThree(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkModThree benchmarks the mod-three implementation
func BenchmarkModThree(b *testing.B) {
	input := "101010101010101010101010"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ModThree(input)
	}
}

// BenchmarkModThreeWithTrace benchmarks with trace generation
func BenchmarkModThreeWithTrace(b *testing.B) {
	input := "101010101010101010101010"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ModThreeWithTrace(input)
	}
}

