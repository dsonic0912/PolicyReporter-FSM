package fsm

// Automaton defines the core interface for finite state automata.
// This interface allows for different implementations while maintaining
// a consistent API for users.
type Automaton[Q State, S Symbol] interface {
	// State management
	GetInitialState() Q
	GetCurrentState() Q
	Reset()
	IsAcceptingState(state Q) bool
	IsCurrentStateAccepting() bool

	// Input processing
	Step(symbol S) (Q, error)
	ProcessInput(input []S) (bool, error)
	ProcessInputWithTrace(input []S) ([]Q, bool, error)

	// Configuration
	Validate() error
	String() string
}

// Builder defines the interface for building automata using the builder pattern.
// This allows for different builder implementations and makes the API more flexible.
type Builder[Q State, S Symbol] interface {
	WithStates(states ...Q) Builder[Q, S]
	WithAlphabet(symbols ...S) Builder[Q, S]
	WithAcceptingStates(states ...Q) Builder[Q, S]
	WithTransition(from Q, symbol S, to Q) Builder[Q, S]
	WithTransitions(transitions ...Transition[Q, S]) Builder[Q, S]
	Build() (Automaton[Q, S], error)
	MustBuild() Automaton[Q, S]
}

// Processor defines an interface for different input processing strategies.
// This allows for pluggable processing algorithms.
type Processor[Q State, S Symbol] interface {
	Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error)
}

// ProcessResult encapsulates the result of processing input through an automaton.
type ProcessResult[Q State] struct {
	Accepted   bool
	Trace      []Q
	FinalState Q
}

// Validator defines an interface for automaton validation strategies.
// This allows for different validation rules and extensible validation.
type Validator[Q State, S Symbol] interface {
	Validate(automaton Automaton[Q, S]) error
}

// StateTransitioner defines an interface for state transition logic.
// This allows for different transition strategies and algorithms.
type StateTransitioner[Q State, S Symbol] interface {
	Transition(currentState Q, symbol S) (Q, error)
}

// AutomatonFactory defines an interface for creating automata.
// This supports the factory pattern for different automaton types.
type AutomatonFactory[Q State, S Symbol] interface {
	CreateAutomaton(initialState Q) Automaton[Q, S]
	CreateBuilder(initialState Q) Builder[Q, S]
}

// Serializer defines an interface for automaton serialization.
// This allows for different serialization formats (JSON, XML, binary, etc.).
type Serializer[Q State, S Symbol] interface {
	Serialize(automaton Automaton[Q, S]) ([]byte, error)
	Deserialize(data []byte) (Automaton[Q, S], error)
}

// Observer defines an interface for observing automaton state changes.
// This implements the observer pattern for monitoring and debugging.
type Observer[Q State, S Symbol] interface {
	OnStateChange(from Q, symbol S, to Q)
	OnInputProcessed(input []S, accepted bool)
	OnError(err error)
}

// AutomatonWithObservers extends the basic Automaton interface with observer support.
type AutomatonWithObservers[Q State, S Symbol] interface {
	Automaton[Q, S]
	AddObserver(observer Observer[Q, S])
	RemoveObserver(observer Observer[Q, S])
	NotifyObservers()
}

// Metrics defines an interface for collecting automaton performance metrics.
type Metrics interface {
	IncrementTransitions()
	IncrementInputsProcessed()
	RecordProcessingTime(duration int64)
	GetTransitionCount() int64
	GetInputsProcessedCount() int64
	GetAverageProcessingTime() float64
}

// AutomatonWithMetrics extends the basic Automaton interface with metrics collection.
type AutomatonWithMetrics[Q State, S Symbol] interface {
	Automaton[Q, S]
	GetMetrics() Metrics
}

// Optimizer defines an interface for automaton optimization strategies.
type Optimizer[Q State, S Symbol] interface {
	Optimize(automaton Automaton[Q, S]) (Automaton[Q, S], error)
}

// AutomatonConverter defines an interface for converting between different automaton types.
type AutomatonConverter[Q1, Q2 State, S1, S2 Symbol] interface {
	Convert(source Automaton[Q1, S1]) (Automaton[Q2, S2], error)
}

// Ensure our concrete types implement the interfaces
var (
	_ Automaton[string, rune] = (*FiniteAutomaton[string, rune])(nil)
	_ Builder[string, rune]   = (*AutomatonBuilder[string, rune])(nil)
)
