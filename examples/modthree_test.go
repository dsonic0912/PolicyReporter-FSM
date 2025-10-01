package examples

import (
	"fmt"
	"strings"
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

// TestModThree_EdgeCases tests edge cases and boundary conditions
func TestModThree_EdgeCases(t *testing.T) {
	t.Run("single bit inputs", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int
		}{
			{"0", 0},
			{"1", 1},
		}

		for _, tt := range tests {
			result, err := ModThree(tt.input)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ModThree(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("leading zeros", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int
		}{
			{"00", 0},
			{"01", 1},
			{"000", 0},
			{"001", 1},
			{"0010", 2},
			{"00110", 0},
		}

		for _, tt := range tests {
			result, err := ModThree(tt.input)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ModThree(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("very long inputs", func(t *testing.T) {
		// Test with very long binary strings
		longInput := strings.Repeat("10", 1000) // 2000 characters
		result, err := ModThree(longInput)
		if err != nil {
			t.Fatalf("ModThree with long input returned error: %v", err)
		}
		// The result should be valid (0, 1, or 2)
		if result < 0 || result > 2 {
			t.Errorf("ModThree result %d is out of range [0, 2]", result)
		}
	})
}

// TestModThree_PropertyBased tests mathematical properties
func TestModThree_PropertyBased(t *testing.T) {
	t.Run("mathematical correctness", func(t *testing.T) {
		// Test that our FSM produces the same results as mathematical mod 3
		testCases := []string{
			"0", "1", "10", "11", "100", "101", "110", "111",
			"1000", "1001", "1010", "1011", "1100", "1101", "1110", "1111",
		}

		for _, binary := range testCases {
			// Convert binary to decimal
			decimal := int64(0)
			for _, bit := range binary {
				decimal = decimal*2 + int64(bit-'0')
			}
			expected := int(decimal % 3)

			result, err := ModThree(binary)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", binary, err)
			}
			if result != expected {
				t.Errorf("ModThree(%q) = %d, want %d (decimal %d)", binary, result, expected, decimal)
			}
		}
	})

	t.Run("step by step processing", func(t *testing.T) {
		// Test that step-by-step processing gives same result as batch processing
		testInputs := []string{"110", "101", "111", "1010"}

		for _, input := range testInputs {
			// Process as batch
			resultBatch, err := ModThree(input)
			if err != nil {
				t.Fatalf("ModThree(%q) returned error: %v", input, err)
			}

			// Process step by step
			fa := NewModThreeAutomaton()
			for _, symbol := range input {
				_, err := fa.Step(symbol)
				if err != nil {
					t.Fatalf("Step(%c) returned error: %v", symbol, err)
				}
			}

			var resultStepwise int
			switch fa.GetCurrentState() {
			case S0:
				resultStepwise = 0
			case S1:
				resultStepwise = 1
			case S2:
				resultStepwise = 2
			}

			if resultBatch != resultStepwise {
				t.Errorf("Step-by-step processing failed: ModThree(%q) = %d, stepwise = %d",
					input, resultBatch, resultStepwise)
			}
		}
	})
}

// TestModThree_StateTransitionCorrectness tests that state transitions follow the mathematical model
func TestModThree_StateTransitionCorrectness(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStates []ModThreeState
	}{
		{
			name:           "binary 6 (110)",
			input:          "110",
			expectedStates: []ModThreeState{S0, S1, S0, S0}, // 0->1->0->0
		},
		{
			name:           "binary 5 (101)",
			input:          "101",
			expectedStates: []ModThreeState{S0, S1, S2, S2}, // 0->1->2->2
		},
		{
			name:           "binary 7 (111)",
			input:          "111",
			expectedStates: []ModThreeState{S0, S1, S0, S1}, // 0->1->0->1
		},
		{
			name:           "all zeros",
			input:          "000",
			expectedStates: []ModThreeState{S0, S0, S0, S0}, // 0->0->0->0
		},
		{
			name:           "alternating",
			input:          "1010",
			expectedStates: []ModThreeState{S0, S1, S2, S2, S1}, // 0->1->2->5->10, mod 3: 0->1->2->2->1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, trace, err := ModThreeWithTrace(tt.input)
			if err != nil {
				t.Fatalf("ModThreeWithTrace(%q) returned error: %v", tt.input, err)
			}

			if len(trace) != len(tt.expectedStates) {
				t.Fatalf("Trace length = %d, want %d", len(trace), len(tt.expectedStates))
			}

			for i, expectedState := range tt.expectedStates {
				if trace[i] != expectedState {
					t.Errorf("State[%d] = %v, want %v", i, trace[i], expectedState)
				}
			}
		})
	}
}

// TestModThree_ErrorHandling tests comprehensive error scenarios
func TestModThree_ErrorHandling(t *testing.T) {
	t.Run("detailed error messages", func(t *testing.T) {
		_, err := ModThree("")
		if err == nil {
			t.Fatal("Expected error for empty string")
		}
		if !strings.Contains(err.Error(), "empty") {
			t.Errorf("Error message should mention 'empty', got: %v", err)
		}
	})

	t.Run("invalid characters with position info", func(t *testing.T) {
		invalidInputs := []string{
			"10a10",
			"1 0",
			"102",
			"abc",
			"10\n10",
			"10\t10",
		}

		for _, input := range invalidInputs {
			_, err := ModThree(input)
			if err == nil {
				t.Errorf("Expected error for invalid input %q", input)
			}
		}
	})
}

// BenchmarkModThree_VariousLengths benchmarks with different input lengths
func BenchmarkModThree_VariousLengths(b *testing.B) {
	lengths := []int{10, 100, 1000, 10000}

	for _, length := range lengths {
		input := strings.Repeat("10", length/2)
		b.Run(fmt.Sprintf("length_%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = ModThree(input)
			}
		})
	}
}

// BenchmarkNewModThreeAutomaton benchmarks automaton creation
func BenchmarkNewModThreeAutomaton(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewModThreeAutomaton()
	}
}
