package fsm

import "testing"

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

	if fa.transitions["q0"]['a'] != "q1" {
		t.Error("Transition q0 --a--> q1 not added")
	}
	if fa.transitions["q1"]['b'] != "q2" {
		t.Error("Transition q1 --b--> q2 not added")
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

