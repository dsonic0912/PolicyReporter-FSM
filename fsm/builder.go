package fsm

// Builder provides a fluent interface for constructing finite automata.
// It helps ensure that all required components are properly configured.
type Builder[Q State, S Symbol] struct {
	automaton *FiniteAutomaton[Q, S]
}

// NewBuilder creates a new builder for constructing a finite automaton.
// The initial state q0 must be specified.
func NewBuilder[Q State, S Symbol](initialState Q) *Builder[Q, S] {
	return &Builder[Q, S]{
		automaton: New[Q, S](initialState),
	}
}

// WithStates adds states to the automaton's set Q.
// The initial state is automatically added.
func (b *Builder[Q, S]) WithStates(states ...Q) *Builder[Q, S] {
	b.automaton.AddStates(states...)
	// Ensure initial state is in Q
	b.automaton.AddState(b.automaton.initialState)
	return b
}

// WithAlphabet sets the input alphabet Σ.
func (b *Builder[Q, S]) WithAlphabet(symbols ...S) *Builder[Q, S] {
	b.automaton.AddSymbols(symbols...)
	return b
}

// WithAcceptingStates sets the accepting states F.
func (b *Builder[Q, S]) WithAcceptingStates(states ...Q) *Builder[Q, S] {
	b.automaton.AddAcceptingStates(states...)
	return b
}

// WithTransition adds a single transition δ(from, symbol) = to.
func (b *Builder[Q, S]) WithTransition(from Q, symbol S, to Q) *Builder[Q, S] {
	b.automaton.AddTransition(from, symbol, to)
	return b
}

// WithTransitions adds multiple transitions at once.
// Each transition is specified as a Transition struct.
func (b *Builder[Q, S]) WithTransitions(transitions ...Transition[Q, S]) *Builder[Q, S] {
	for _, t := range transitions {
		b.automaton.AddTransition(t.From, t.Symbol, t.To)
	}
	return b
}

// Build finalizes the automaton and validates its configuration.
// Returns an error if the automaton is not properly configured.
func (b *Builder[Q, S]) Build() (*FiniteAutomaton[Q, S], error) {
	if err := b.automaton.Validate(); err != nil {
		return nil, err
	}
	return b.automaton, nil
}

// MustBuild finalizes the automaton and panics if validation fails.
// Use this when you're certain the configuration is correct.
func (b *Builder[Q, S]) MustBuild() *FiniteAutomaton[Q, S] {
	fa, err := b.Build()
	if err != nil {
		panic(err)
	}
	return fa
}

// Transition represents a single state transition δ(From, Symbol) = To.
type Transition[Q State, S Symbol] struct {
	From   Q
	Symbol S
	To     Q
}

// T is a convenience function for creating Transition structs.
func T[Q State, S Symbol](from Q, symbol S, to Q) Transition[Q, S] {
	return Transition[Q, S]{From: from, Symbol: symbol, To: to}
}

