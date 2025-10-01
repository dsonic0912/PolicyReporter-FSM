package fsm

import (
	"testing"
)

// TestAutomatonInterface tests that FiniteAutomaton implements Automaton interface
func TestAutomatonInterface(t *testing.T) {
	var automaton Automaton[string, rune] = New[string, rune]("q0")

	// Test interface methods
	if automaton.GetInitialState() != "q0" {
		t.Error("GetInitialState() failed")
	}

	if automaton.GetCurrentState() != "q0" {
		t.Error("GetCurrentState() failed")
	}

	if automaton.IsCurrentStateAccepting() {
		t.Error("IsCurrentStateAccepting() should be false initially")
	}

	automaton.Reset()
	if automaton.GetCurrentState() != "q0" {
		t.Error("Reset() failed")
	}
}

// TestBuilderInterface tests that AutomatonBuilder implements Builder interface
func TestBuilderInterface(t *testing.T) {
	var builder Builder[string, rune] = NewBuilder[string, rune]("q0")

	// Test fluent interface
	automaton := builder.
		WithStates("q0", "q1").
		WithAlphabet('a', 'b').
		WithAcceptingStates("q1").
		WithTransition("q0", 'a', "q1").
		MustBuild()

	if automaton.GetInitialState() != "q0" {
		t.Error("Builder interface failed")
	}
}

// TestProcessorInterface tests processor implementations
func TestProcessorInterface(t *testing.T) {
	automaton := NewBuilder[string, rune]("q0").
		WithStates("q0", "q1").
		WithAlphabet('a', 'b').
		WithAcceptingStates("q1").
		WithTransition("q0", 'a', "q1").
		MustBuild()

	// Test StandardProcessor
	var processor Processor[string, rune] = NewStandardProcessor[string, rune]()
	result, err := processor.Process(automaton, []rune("a"))
	if err != nil {
		t.Fatalf("StandardProcessor failed: %v", err)
	}
	if !result.Accepted {
		t.Error("StandardProcessor should accept input 'a'")
	}

	// Test OptimizedProcessor
	processor = NewOptimizedProcessor[string, rune]()
	result, err = processor.Process(automaton, []rune("a"))
	if err != nil {
		t.Fatalf("OptimizedProcessor failed: %v", err)
	}
	if !result.Accepted {
		t.Error("OptimizedProcessor should accept input 'a'")
	}
}

// TestFactoryInterface tests factory implementations
func TestFactoryInterface(t *testing.T) {
	var factory AutomatonFactory[string, rune] = NewStandardFactory[string, rune]()

	// Test CreateAutomaton
	automaton := factory.CreateAutomaton("q0")
	if automaton.GetInitialState() != "q0" {
		t.Error("Factory CreateAutomaton failed")
	}

	// Test CreateBuilder
	builder := factory.CreateBuilder("q0")
	automaton = builder.WithStates("q0").WithAlphabet('a').MustBuild()
	if automaton.GetInitialState() != "q0" {
		t.Error("Factory CreateBuilder failed")
	}
}

// TestErrorTypes tests the error type system
func TestErrorTypes(t *testing.T) {
	// Test error creation
	err := NewValidationError("test validation error")
	if err.Type != ErrorTypeValidation {
		t.Error("NewValidationError should create validation error")
	}

	err = NewTransitionError("q0", 'a', "test transition error")
	if err.Type != ErrorTypeTransition {
		t.Error("NewTransitionError should create transition error")
	}

	// Test error checking functions
	validationErr := NewValidationError("validation")
	if !IsValidationError(validationErr) {
		t.Error("IsValidationError should return true for validation errors")
	}

	transitionErr := NewTransitionError("q0", 'a', "transition")
	if !IsTransitionError(transitionErr) {
		t.Error("IsTransitionError should return true for transition errors")
	}
}

// TestErrorCollector tests the error collector
func TestErrorCollector(t *testing.T) {
	collector := NewErrorCollector()

	if collector.HasErrors() {
		t.Error("New collector should not have errors")
	}

	collector.Add(NewValidationError("error 1"))
	collector.Add(NewValidationError("error 2"))

	if !collector.HasErrors() {
		t.Error("Collector should have errors after adding")
	}

	if len(collector.Errors()) != 2 {
		t.Error("Collector should have 2 errors")
	}

	err := collector.ToError()
	if err == nil {
		t.Error("ToError should return error when collector has errors")
	}
}

// TestProcessorRegistry tests the processor registry
func TestProcessorRegistry(t *testing.T) {
	registry := NewProcessorRegistry[string, rune]()

	// Test default processor
	defaultProcessor := registry.GetDefault()
	if defaultProcessor == nil {
		t.Error("Registry should have default processor")
	}

	// Test getting standard processor
	processor, exists := registry.Get("standard")
	if !exists {
		t.Error("Registry should have standard processor")
	}
	if processor == nil {
		t.Error("Standard processor should not be nil")
	}

	// Test registering new processor
	customProcessor := NewOptimizedProcessor[string, rune]()
	registry.Register("optimized", customProcessor)

	_, exists = registry.Get("optimized")
	if !exists {
		t.Error("Registry should have optimized processor after registration")
	}
}

// TestFactoryRegistry tests the factory registry
func TestFactoryRegistry(t *testing.T) {
	registry := NewFactoryRegistry[string, rune]()

	// Test default factory
	defaultFactory := registry.GetDefault()
	if defaultFactory == nil {
		t.Error("Registry should have default factory")
	}

	// Test getting standard factory
	factory, exists := registry.Get("standard")
	if !exists {
		t.Error("Registry should have standard factory")
	}
	if factory == nil {
		t.Error("Standard factory should not be nil")
	}

	// Test creating automaton with default factory
	automaton := registry.CreateAutomaton("q0")
	if automaton.GetInitialState() != "q0" {
		t.Error("Registry CreateAutomaton failed")
	}

	// Test creating builder with default factory
	builder := registry.CreateBuilder("q0")
	automaton = builder.WithStates("q0").WithAlphabet('a').MustBuild()
	if automaton.GetInitialState() != "q0" {
		t.Error("Registry CreateBuilder failed")
	}
}

// TestProcessorChain tests the processor chain
func TestProcessorChain(t *testing.T) {
	automaton := NewBuilder[string, rune]("q0").
		WithStates("q0", "q1").
		WithAlphabet('a').
		WithAcceptingStates("q1").
		WithTransition("q0", 'a', "q1").
		MustBuild()

	// Create processor chain
	chain := NewProcessorChain[string, rune](
		NewStandardProcessor[string, rune](),
		NewOptimizedProcessor[string, rune](),
	)

	result, err := chain.Process(automaton, []rune("a"))
	if err != nil {
		t.Fatalf("ProcessorChain failed: %v", err)
	}
	if !result.Accepted {
		t.Error("ProcessorChain should accept input 'a'")
	}

	// Test empty chain
	emptyChain := NewProcessorChain[string, rune]()
	_, err = emptyChain.Process(automaton, []rune("a"))
	if err == nil {
		t.Error("Empty processor chain should return error")
	}
}

// TestValidatingProcessor tests the validating processor
func TestValidatingProcessor(t *testing.T) {
	automaton := NewBuilder[string, rune]("q0").
		WithStates("q0", "q1").
		WithAlphabet('a').
		WithAcceptingStates("q1").
		WithTransition("q0", 'a', "q1").
		MustBuild()

	// Create validating processor with validator that rejects empty input
	validator := func(input []rune) error {
		if len(input) == 0 {
			return NewInvalidInputError(rune(0), 0, "empty input not allowed")
		}
		return nil
	}

	processor := NewValidatingProcessor[string, rune](
		NewStandardProcessor[string, rune](),
		validator,
	)

	// Test valid input
	result, err := processor.Process(automaton, []rune("a"))
	if err != nil {
		t.Fatalf("ValidatingProcessor failed with valid input: %v", err)
	}
	if !result.Accepted {
		t.Error("ValidatingProcessor should accept valid input")
	}

	// Test invalid input
	_, err = processor.Process(automaton, []rune{})
	if err == nil {
		t.Error("ValidatingProcessor should reject invalid input")
	}
	if !IsInvalidInputError(err) {
		t.Error("ValidatingProcessor should return InvalidInputError")
	}
}
