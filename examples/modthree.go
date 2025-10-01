// Package examples provides example implementations using the FSM library.
package examples

import (
	"fmt"

	"github.com/dsonic0912/PolicyReporter-FSM/fsm"
)

// Modulo three remainders
const (
	RemainderZero = 0
	RemainderOne  = 1
	RemainderTwo  = 2
)

// ModThreeState represents states in the mod-three automaton.
type ModThreeState string

// States for the mod-three automaton
const (
	S0 ModThreeState = "S0" // Remainder 0
	S1 ModThreeState = "S1" // Remainder 1
	S2 ModThreeState = "S2" // Remainder 2
)

// NewModThreeAutomaton creates a finite automaton that computes modulo 3 of binary numbers.
//
// Based on the formal definition:
//
//	Q = {S0, S1, S2}
//	Σ = {'0', '1'}
//	q0 = S0
//	F = {S0, S1, S2} (all states are accepting)
//	δ is defined by the transition table:
//	  δ(S0, '0') = S0; δ(S0, '1') = S1
//	  δ(S1, '0') = S2; δ(S1, '1') = S0
//	  δ(S2, '0') = S1; δ(S2, '1') = S2
func NewModThreeAutomaton() fsm.Automaton[ModThreeState, rune] {
	return fsm.NewBuilder[ModThreeState, rune](S0).
		// Q: Set of states
		WithStates(S0, S1, S2).
		// Σ: Input alphabet
		WithAlphabet('0', '1').
		// F: All states are accepting (we want the final state value)
		WithAcceptingStates(S0, S1, S2).
		// δ: Transition function
		WithTransitions(
			fsm.T(S0, '0', S0), // (0 * 2 + 0) % 3 = 0
			fsm.T(S0, '1', S1), // (0 * 2 + 1) % 3 = 1
			fsm.T(S1, '0', S2), // (1 * 2 + 0) % 3 = 2
			fsm.T(S1, '1', S0), // (1 * 2 + 1) % 3 = 0
			fsm.T(S2, '0', S1), // (2 * 2 + 0) % 3 = 1
			fsm.T(S2, '1', S2), // (2 * 2 + 1) % 3 = 2
		).
		MustBuild()
}

// ModThree computes the remainder when a binary number is divided by 3.
// Returns the remainder (0, 1, or 2) and any error encountered.
func ModThree(binaryString string) (int, error) {
	if binaryString == "" {
		return 0, fmt.Errorf("input string cannot be empty")
	}

	// Create the automaton
	fa := NewModThreeAutomaton()

	// Convert string to rune slice
	input := []rune(binaryString)

	// Process the input
	accepted, err := fa.ProcessInput(input)
	if err != nil {
		return 0, err
	}

	if !accepted {
		return 0, fmt.Errorf("input was not accepted by the automaton")
	}

	// Map final state to remainder value
	finalState := fa.GetCurrentState()
	switch finalState {
	case S0:
		return RemainderZero, nil
	case S1:
		return RemainderOne, nil
	case S2:
		return RemainderTwo, nil
	default:
		return RemainderZero, fmt.Errorf("unexpected final state: %v", finalState)
	}
}

// ModThreeWithTrace computes mod 3 and returns the state transition trace.
func ModThreeWithTrace(binaryString string) (int, []ModThreeState, error) {
	if binaryString == "" {
		return 0, nil, fmt.Errorf("input string cannot be empty")
	}

	fa := NewModThreeAutomaton()
	input := []rune(binaryString)

	trace, accepted, err := fa.ProcessInputWithTrace(input)
	if err != nil {
		return 0, trace, err
	}

	if !accepted {
		return 0, trace, fmt.Errorf("input was not accepted by the automaton")
	}

	finalState := fa.GetCurrentState()
	var result int
	switch finalState {
	case S0:
		result = 0
	case S1:
		result = 1
	case S2:
		result = 2
	}

	return result, trace, nil
}

// PrintModThreeTrace prints the state transitions for a binary input.
func PrintModThreeTrace(binaryString string) {
	fmt.Printf("Input: \"%s\"\n", binaryString)

	result, trace, err := ModThreeWithTrace(binaryString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print transitions
	for i := 0; i < len(binaryString); i++ {
		fmt.Printf("%d. Current state = %s, Input = %c, result state = %s\n",
			i+1, trace[i], binaryString[i], trace[i+1])
	}

	fmt.Println("No more input")
	fmt.Printf("Print output value (output for state %s = %d) <---- This is the answer\n",
		trace[len(trace)-1], result)
}
