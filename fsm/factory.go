package fsm

// StandardAutomatonFactory is the default factory for creating finite automata.
type StandardAutomatonFactory[Q State, S Symbol] struct{}

// NewStandardFactory creates a new standard automaton factory.
func NewStandardFactory[Q State, S Symbol]() *StandardAutomatonFactory[Q, S] {
	return &StandardAutomatonFactory[Q, S]{}
}

// CreateAutomaton creates a new finite automaton with the given initial state.
func (f *StandardAutomatonFactory[Q, S]) CreateAutomaton(initialState Q) Automaton[Q, S] {
	return New[Q, S](initialState)
}

// CreateBuilder creates a new builder for constructing finite automata.
func (f *StandardAutomatonFactory[Q, S]) CreateBuilder(initialState Q) Builder[Q, S] {
	return NewBuilder[Q, S](initialState)
}

// OptimizedAutomatonFactory creates automata with built-in optimizations.
type OptimizedAutomatonFactory[Q State, S Symbol] struct {
	optimizer Optimizer[Q, S]
}

// NewOptimizedFactory creates a factory that applies optimizations to created automata.
func NewOptimizedFactory[Q State, S Symbol](optimizer Optimizer[Q, S]) *OptimizedAutomatonFactory[Q, S] {
	return &OptimizedAutomatonFactory[Q, S]{
		optimizer: optimizer,
	}
}

// CreateAutomaton creates an optimized finite automaton.
func (f *OptimizedAutomatonFactory[Q, S]) CreateAutomaton(initialState Q) Automaton[Q, S] {
	automaton := New[Q, S](initialState)
	if f.optimizer != nil {
		optimized, err := f.optimizer.Optimize(automaton)
		if err == nil {
			return optimized
		}
		// Fall back to unoptimized if optimization fails
	}
	return automaton
}

// CreateBuilder creates a builder that produces optimized automata.
func (f *OptimizedAutomatonFactory[Q, S]) CreateBuilder(initialState Q) Builder[Q, S] {
	return &OptimizedBuilder[Q, S]{
		builder:   NewBuilder[Q, S](initialState),
		optimizer: f.optimizer,
	}
}

// OptimizedBuilder wraps a standard builder and applies optimizations.
type OptimizedBuilder[Q State, S Symbol] struct {
	builder   Builder[Q, S]
	optimizer Optimizer[Q, S]
}

// WithStates adds states to the automaton.
func (b *OptimizedBuilder[Q, S]) WithStates(states ...Q) Builder[Q, S] {
	b.builder = b.builder.WithStates(states...)
	return b
}

// WithAlphabet sets the input alphabet.
func (b *OptimizedBuilder[Q, S]) WithAlphabet(symbols ...S) Builder[Q, S] {
	b.builder = b.builder.WithAlphabet(symbols...)
	return b
}

// WithAcceptingStates sets the accepting states.
func (b *OptimizedBuilder[Q, S]) WithAcceptingStates(states ...Q) Builder[Q, S] {
	b.builder = b.builder.WithAcceptingStates(states...)
	return b
}

// WithTransition adds a single transition.
func (b *OptimizedBuilder[Q, S]) WithTransition(from Q, symbol S, to Q) Builder[Q, S] {
	b.builder = b.builder.WithTransition(from, symbol, to)
	return b
}

// WithTransitions adds multiple transitions.
func (b *OptimizedBuilder[Q, S]) WithTransitions(transitions ...Transition[Q, S]) Builder[Q, S] {
	b.builder = b.builder.WithTransitions(transitions...)
	return b
}

// Build creates and optimizes the automaton.
func (b *OptimizedBuilder[Q, S]) Build() (Automaton[Q, S], error) {
	automaton, err := b.builder.Build()
	if err != nil {
		return nil, err
	}

	if b.optimizer != nil {
		optimized, err := b.optimizer.Optimize(automaton)
		if err == nil {
			return optimized, nil
		}
		// Fall back to unoptimized if optimization fails
	}

	return automaton, nil
}

// MustBuild creates and optimizes the automaton, panicking on error.
func (b *OptimizedBuilder[Q, S]) MustBuild() Automaton[Q, S] {
	automaton, err := b.Build()
	if err != nil {
		panic(err)
	}
	return automaton
}

// FactoryRegistry manages different automaton factories.
type FactoryRegistry[Q State, S Symbol] struct {
	factories      map[string]AutomatonFactory[Q, S]
	defaultFactory AutomatonFactory[Q, S]
}

// NewFactoryRegistry creates a new factory registry.
func NewFactoryRegistry[Q State, S Symbol]() *FactoryRegistry[Q, S] {
	registry := &FactoryRegistry[Q, S]{
		factories: make(map[string]AutomatonFactory[Q, S]),
	}

	// Register default factory
	defaultFactory := NewStandardFactory[Q, S]()
	registry.Register("standard", defaultFactory)
	registry.SetDefault(defaultFactory)

	return registry
}

// Register adds a factory to the registry.
func (r *FactoryRegistry[Q, S]) Register(name string, factory AutomatonFactory[Q, S]) {
	r.factories[name] = factory
}

// Get retrieves a factory by name.
func (r *FactoryRegistry[Q, S]) Get(name string) (AutomatonFactory[Q, S], bool) {
	factory, exists := r.factories[name]
	return factory, exists
}

// GetDefault returns the default factory.
func (r *FactoryRegistry[Q, S]) GetDefault() AutomatonFactory[Q, S] {
	return r.defaultFactory
}

// SetDefault sets the default factory.
func (r *FactoryRegistry[Q, S]) SetDefault(factory AutomatonFactory[Q, S]) {
	r.defaultFactory = factory
}

// CreateAutomaton creates an automaton using the default factory.
func (r *FactoryRegistry[Q, S]) CreateAutomaton(initialState Q) Automaton[Q, S] {
	return r.defaultFactory.CreateAutomaton(initialState)
}

// CreateBuilder creates a builder using the default factory.
func (r *FactoryRegistry[Q, S]) CreateBuilder(initialState Q) Builder[Q, S] {
	return r.defaultFactory.CreateBuilder(initialState)
}

// CreateAutomatonWith creates an automaton using a named factory.
func (r *FactoryRegistry[Q, S]) CreateAutomatonWith(factoryName string, initialState Q) (Automaton[Q, S], error) {
	factory, exists := r.Get(factoryName)
	if !exists {
		return nil, NewError(ErrorTypeInvalidConfiguration, "factory not found: "+factoryName)
	}
	return factory.CreateAutomaton(initialState), nil
}

// CreateBuilderWith creates a builder using a named factory.
func (r *FactoryRegistry[Q, S]) CreateBuilderWith(factoryName string, initialState Q) (Builder[Q, S], error) {
	factory, exists := r.Get(factoryName)
	if !exists {
		return nil, NewError(ErrorTypeInvalidConfiguration, "factory not found: "+factoryName)
	}
	return factory.CreateBuilder(initialState), nil
}

// Global factory registry for convenience
var globalRegistry = NewFactoryRegistry[any, any]()

// RegisterFactory registers a factory globally.
func RegisterFactory[Q State, S Symbol](name string, factory AutomatonFactory[Q, S]) {
	// Type assertion to store in global registry
	globalRegistry.Register(name, factory.(AutomatonFactory[any, any]))
}

// GetFactory retrieves a factory from the global registry.
func GetFactory[Q State, S Symbol](name string) (AutomatonFactory[Q, S], bool) {
	factory, exists := globalRegistry.Get(name)
	if !exists {
		return nil, false
	}
	// Type assertion to return correct type
	return factory.(AutomatonFactory[Q, S]), true
}
