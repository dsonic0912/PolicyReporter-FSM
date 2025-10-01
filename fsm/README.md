# FSM Package - Generic Finite State Automaton Library

A production-ready, type-safe finite state automaton library for Go.

## Overview

This package provides a complete implementation of finite state automata (FSA) based on the formal mathematical definition. It uses Go generics to provide type safety while maintaining flexibility.

## Formal Definition

A finite automaton is defined as a 5-tuple **(Q, Σ, q₀, F, δ)** where:

- **Q**: A finite set of states
- **Σ** (Sigma): A finite input alphabet  
- **q₀**: The initial state (q₀ ∈ Q)
- **F**: The set of accepting/final states (F ⊆ Q)
- **δ** (Delta): The transition function (δ: Q × Σ → Q)

This package directly implements this definition with a clean, intuitive API.

## Quick Start

```go
import "github.com/dsonic0912/PolicyReporter-FSM/fsm"

// Create an automaton using the builder
fa := fsm.NewBuilder[string, rune]("q0").
    WithStates("q0", "q1", "q2").
    WithAlphabet('a', 'b').
    WithAcceptingStates("q2").
    WithTransitions(
        fsm.T("q0", 'a', "q1"),
        fsm.T("q1", 'b', "q2"),
    ).
    MustBuild()

// Process input
accepted, err := fa.ProcessInput([]rune("ab"))
```

## API Reference

### Creating an Automaton

#### Using the Builder (Recommended)

```go
fa := fsm.NewBuilder[StateType, SymbolType](initialState).
    WithStates(state1, state2, ...).
    WithAlphabet(symbol1, symbol2, ...).
    WithAcceptingStates(acceptingState1, ...).
    WithTransition(from, symbol, to).
    Build() // or MustBuild()
```

#### Direct Construction

```go
fa := fsm.New[StateType, SymbolType](initialState)
fa.AddStates(state1, state2)
fa.AddSymbols(symbol1, symbol2)
fa.AddAcceptingStates(acceptingState1)
fa.AddTransition(from, symbol, to)
```

### Processing Input

#### Process and Check Acceptance

```go
accepted, err := fa.ProcessInput([]SymbolType{sym1, sym2, ...})
if err != nil {
    // Handle error (invalid symbol or undefined transition)
}
if accepted {
    // Input was accepted (ended in accepting state)
}
```

#### Process with State Trace

```go
trace, accepted, err := fa.ProcessInputWithTrace(input)
// trace contains the sequence of states visited
```

#### Step-by-Step Processing

```go
fa.Reset() // Start from initial state
for _, symbol := range input {
    nextState, err := fa.Step(symbol)
    if err != nil {
        // Handle error
    }
    // Process nextState
}
```

### Querying the Automaton

```go
// Get states
initialState := fa.GetInitialState()
currentState := fa.GetCurrentState()

// Check accepting states
isAccepting := fa.IsAcceptingState(someState)
isCurrentAccepting := fa.IsCurrentStateAccepting()

// Validate configuration
if err := fa.Validate(); err != nil {
    // Automaton is not properly configured
}

// Get string representation
fmt.Println(fa.String())
```

## Type Parameters

The library uses Go generics with two type parameters:

### `Q` - State Type

The type used for states. Must be comparable.

**Examples:**
- `string` - e.g., "q0", "q1", "accepting"
- `int` - e.g., 0, 1, 2
- Custom types - e.g., `type State string`

### `S` - Symbol Type  

The type used for input symbols. Must be comparable.

**Examples:**
- `rune` - for character-based alphabets
- `string` - for string-based alphabets
- `int` - for numeric alphabets
- Custom types - e.g., `type Event string`

## Examples

### Example 1: Binary String Validator

Accepts binary strings with even number of 1s:

```go
fa := fsm.NewBuilder[string, rune]("even").
    WithStates("even", "odd").
    WithAlphabet('0', '1').
    WithAcceptingStates("even").
    WithTransitions(
        fsm.T("even", '0', "even"),
        fsm.T("even", '1', "odd"),
        fsm.T("odd", '0', "odd"),
        fsm.T("odd", '1', "even"),
    ).
    MustBuild()

accepted, _ := fa.ProcessInput([]rune("110011")) // true
accepted, _ = fa.ProcessInput([]rune("101"))     // false
```

### Example 2: Pattern Matcher

Accepts strings ending with "ab":

```go
fa := fsm.NewBuilder[int, rune](0).
    WithStates(0, 1, 2).
    WithAlphabet('a', 'b', 'c').
    WithAcceptingStates(2).
    WithTransitions(
        fsm.T(0, 'a', 1), fsm.T(0, 'b', 0), fsm.T(0, 'c', 0),
        fsm.T(1, 'a', 1), fsm.T(1, 'b', 2), fsm.T(1, 'c', 0),
        fsm.T(2, 'a', 1), fsm.T(2, 'b', 0), fsm.T(2, 'c', 0),
    ).
    MustBuild()

accepted, _ := fa.ProcessInput([]rune("xyzab"))  // true
accepted, _ = fa.ProcessInput([]rune("abc"))     // false
```

### Example 3: Custom Types

```go
type State string
type Input int

const (
    Start State = "START"
    End   State = "END"
)

const (
    Zero Input = 0
    One  Input = 1
)

fa := fsm.NewBuilder[State, Input](Start).
    WithStates(Start, End).
    WithAlphabet(Zero, One).
    WithAcceptingStates(End).
    WithTransition(Start, One, End).
    MustBuild()
```

## Error Handling

The library provides detailed error messages for common issues:

### Validation Errors

```go
fa, err := builder.Build()
if err != nil {
    // Possible errors:
    // - "initial state X is not in the set of states Q"
    // - "accepting state X is not in the set of states Q"
    // - "transition from state X, but state is not in Q"
    // - "transition uses symbol X, but symbol is not in Σ"
    // - "transition to state X, but state is not in Q"
}
```

### Runtime Errors

```go
state, err := fa.Step(symbol)
if err != nil {
    // Possible errors:
    // - "symbol not in alphabet: X"
    // - "no transition defined for state X with symbol Y"
}
```

## Best Practices

### 1. Use the Builder Pattern

```go
// ✅ Good - Clear and validated
fa := fsm.NewBuilder[string, rune]("q0").
    WithStates("q0", "q1").
    WithAlphabet('a', 'b').
    WithAcceptingStates("q1").
    WithTransition("q0", 'a', "q1").
    MustBuild()

// ❌ Avoid - Manual construction without validation
fa := fsm.New[string, rune]("q0")
fa.AddState("q0")
// ... easy to forget validation
```

### 2. Use MustBuild for Static Automata

```go
// ✅ Good - Automaton structure is known at compile time
var ModThreeFA = fsm.NewBuilder[string, rune]("S0").
    WithStates("S0", "S1", "S2").
    // ...
    MustBuild() // Panic if invalid - fail fast

// ✅ Good - Automaton structure is dynamic
fa, err := fsm.NewBuilder[string, rune](initialState).
    WithStates(dynamicStates...).
    // ...
    Build() // Return error for runtime handling
```

### 3. Validate Input Before Processing

```go
// ✅ Good - Validate automaton configuration
if err := fa.Validate(); err != nil {
    return fmt.Errorf("invalid automaton: %w", err)
}

// Process input
accepted, err := fa.ProcessInput(input)
```

### 4. Use Appropriate Types

```go
// ✅ Good - Semantic types
type State string
type Event rune

fa := fsm.NewBuilder[State, Event]("idle")

// ❌ Avoid - Generic types lose meaning
fa := fsm.NewBuilder[string, rune]("idle")
```

## Performance Considerations

- **State Transitions**: O(1) map lookup
- **Input Processing**: O(n) where n is input length
- **Memory**: O(|Q| × |Σ|) for transition table
- **Allocations**: Zero allocations in hot path (ProcessInput)

## Thread Safety

The `FiniteAutomaton` type is **not thread-safe**. If you need to process multiple inputs concurrently:

```go
// Create separate instances for each goroutine
for i := 0; i < numWorkers; i++ {
    go func() {
        fa := NewModThreeAutomaton() // Each goroutine gets its own instance
        // Process inputs...
    }()
}
```

## Limitations

1. **Deterministic Only**: This library implements DFA (Deterministic Finite Automata), not NFA
2. **No ε-transitions**: Epsilon transitions are not supported
3. **Single Transition**: Only one transition per (state, symbol) pair
4. **Comparable Types**: State and symbol types must be comparable (no slices, maps, or functions)

## See Also

- [Examples Package](../examples/) - Real-world usage examples including mod-three
- [Main README](../README.md) - Project overview and quick start
- [Go Generics Documentation](https://go.dev/doc/tutorial/generics)

