# Generic Finite State Automaton Library

A production-ready, generic finite state automaton (FSA) library for Go, with a complete implementation of the mod-three procedure as an example.

## Overview

This project provides:

1. **Generic FSM Library** (`fsm` package): A reusable library for building any finite automaton based on the formal definition
2. **Mod-Three Example** (`examples` package): A complete implementation of modulo-three calculation using the library
3. **Comprehensive Tests**: 100+ test cases covering the library and examples

The library implements the formal definition of a finite automaton as a 5-tuple **(Q, Σ, q₀, F, δ)** where:
- **Q** is a finite set of states
- **Σ** (Sigma) is a finite input alphabet
- **q₀** is the initial state (q₀ ∈ Q)
- **F** is the set of accepting/final states (F ⊆ Q)
- **δ** (delta) is the transition function (δ: Q × Σ → Q)

## Quick Start

### Installation

```bash
go get github.com/dsonic0912/PolicyReporter-FSM
```

### Using the Mod-Three Example

```go
package main

import (
    "fmt"
    "github.com/dsonic0912/PolicyReporter-FSM/examples"
)

func main() {
    // Simple usage
    result, err := examples.ModThree("110")
    if err != nil {
        panic(err)
    }
    fmt.Printf("110 (binary) mod 3 = %d\n", result) // Output: 0

    // With trace
    examples.PrintModThreeTrace("110")
}
```

### Building Your Own FSM

```go
package main

import (
    "fmt"
    "github.com/dsonic0912/PolicyReporter-FSM/fsm"
)

func main() {
    // Create an automaton that accepts binary strings ending in "01"
    fa := fsm.NewBuilder[string, rune]("q0").
        WithStates("q0", "q1", "q2").
        WithAlphabet('0', '1').
        WithAcceptingStates("q2").
        WithTransitions(
            fsm.T("q0", '0', "q1"),
            fsm.T("q0", '1', "q0"),
            fsm.T("q1", '0', "q1"),
            fsm.T("q1", '1', "q2"),
            fsm.T("q2", '0', "q1"),
            fsm.T("q2", '1', "q0"),
        ).
        MustBuild()

    accepted, _ := fa.ProcessInput([]rune("1001"))
    fmt.Printf("Accepted: %v\n", accepted) // Output: true
}
```

## Library API Reference

### Core Types

#### `FiniteAutomaton[Q State, S Symbol]`

The main automaton type. Generic over state type `Q` and symbol type `S`.

**Methods:**
- `AddState(state Q)` - Add a state to Q
- `AddStates(states ...Q)` - Add multiple states to Q
- `AddSymbol(symbol S)` - Add a symbol to Σ
- `AddSymbols(symbols ...S)` - Add multiple symbols to Σ
- `AddAcceptingState(state Q)` - Add a state to F
- `AddAcceptingStates(states ...Q)` - Add multiple states to F
- `AddTransition(from Q, symbol S, to Q)` - Add a transition to δ
- `GetInitialState()` - Get q₀
- `GetCurrentState()` - Get current state during processing
- `IsAcceptingState(state Q)` - Check if state is in F
- `Reset()` - Reset to initial state
- `Step(symbol S)` - Process one symbol and transition
- `ProcessInput(input []S)` - Process a sequence of symbols
- `ProcessInputWithTrace(input []S)` - Process input and return state trace
- `Validate()` - Validate automaton configuration
- `String()` - Get string representation

#### `Builder[Q State, S Symbol]`

Fluent builder for constructing automata.

**Methods:**
- `NewBuilder[Q, S](initialState Q)` - Create a new builder
- `WithStates(states ...Q)` - Set states Q
- `WithAlphabet(symbols ...S)` - Set alphabet Σ
- `WithAcceptingStates(states ...Q)` - Set accepting states F
- `WithTransition(from Q, symbol S, to Q)` - Add a transition
- `WithTransitions(transitions ...Transition[Q, S])` - Add multiple transitions
- `Build()` - Build and validate (returns error)
- `MustBuild()` - Build and validate (panics on error)

**Helper Functions:**
- `T[Q, S](from Q, symbol S, to Q)` - Create a Transition struct

## Mod-Three Example

The mod-three automaton computes the remainder when a binary number is divided by 3.

### Formal Definition

Based on the formal notation:
- **Q** = {S0, S1, S2}
- **Σ** = {'0', '1'}
- **q₀** = S0
- **F** = {S0, S1, S2} (all states are accepting)
- **δ** (Transition function):

| Current State | Input '0' | Input '1' |
|--------------|-----------|-----------|
| S0 (rem=0)   | S0        | S1        |
| S1 (rem=1)   | S2        | S0        |
| S2 (rem=2)   | S1        | S2        |

The transitions follow: `δ(state, symbol) = (state * 2 + symbol) % 3`

### Example Execution

For input "110" (binary 6, which is 0 mod 3):

```
1. Current state = S0, Input = 1, result state = S1
2. Current state = S1, Input = 1, result state = S0
3. Current state = S0, Input = 0, result state = S0
No more input
Print output value (output for state S0 = 0) <---- This is the answer
```

## Features

### Generic Library Features
- **Type-Safe Generics**: Works with any comparable types for states and symbols
- **Formal Definition**: Directly implements the mathematical definition of FSA
- **Builder Pattern**: Fluent API for easy automaton construction
- **Validation**: Comprehensive validation of automaton configuration
- **Error Handling**: Clear error messages for invalid inputs and configurations
- **Trace Support**: Optional state transition tracing for debugging
- **Zero Dependencies**: Pure Go implementation

### Mod-Three Implementation Features
- **Multiple APIs**: Simple function, trace function, and direct automaton access
- **Comprehensive Testing**: 100+ test cases
- **Performance**: Optimized for speed with zero allocations
- **Well-Documented**: Clear examples and documentation

## Advanced Usage Examples

### Example 1: Even/Odd Parity Checker

```go
// Accepts binary strings with even number of 1s
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

accepted, _ := fa.ProcessInput([]rune("1100")) // true (two 1s)
```

### Example 2: Pattern Matcher

```go
// Accepts strings containing "abc" as a substring
fa := fsm.NewBuilder[int, rune](0).
    WithStates(0, 1, 2, 3).
    WithAlphabet('a', 'b', 'c').
    WithAcceptingStates(3).
    WithTransitions(
        fsm.T(0, 'a', 1), fsm.T(0, 'b', 0), fsm.T(0, 'c', 0),
        fsm.T(1, 'a', 1), fsm.T(1, 'b', 2), fsm.T(1, 'c', 0),
        fsm.T(2, 'a', 1), fsm.T(2, 'b', 0), fsm.T(2, 'c', 3),
        fsm.T(3, 'a', 3), fsm.T(3, 'b', 3), fsm.T(3, 'c', 3),
    ).
    MustBuild()

accepted, _ := fa.ProcessInput([]rune("xyzabc")) // true
```

### Example 3: Using Custom Types

```go
type TrafficLight string
type Event string

const (
    Red    TrafficLight = "RED"
    Yellow TrafficLight = "YELLOW"
    Green  TrafficLight = "GREEN"
)

const (
    Timer     Event = "TIMER"
    Emergency Event = "EMERGENCY"
)

fa := fsm.NewBuilder[TrafficLight, Event](Red).
    WithStates(Red, Yellow, Green).
    WithAlphabet(Timer, Emergency).
    WithAcceptingStates(Red, Yellow, Green).
    WithTransitions(
        fsm.T(Red, Timer, Green),
        fsm.T(Green, Timer, Yellow),
        fsm.T(Yellow, Timer, Red),
        fsm.T(Red, Emergency, Red),
        fsm.T(Green, Emergency, Red),
        fsm.T(Yellow, Emergency, Red),
    ).
    MustBuild()
```

## Running the Demo

```bash
go run main.go
```

Output:
```
=== Mod-Three Finite State Machine ===
Using Generic FSM Library

Example 1: Input = "110"
Input: "110"
1. Current state = S0, Input = 1, result state = S1
2. Current state = S1, Input = 1, result state = S0
3. Current state = S0, Input = 0, result state = S0
No more input
Print output value (output for state S0 = 0) <---- This is the answer

==================================================

Automaton Configuration:

Finite Automaton:
  Q (States): [S2 S0 S1]
  Σ (Alphabet): [49 48]
  q0 (Initial): S0
  F (Accepting): [S0 S1 S2]
  δ (Transitions):
    δ(S0, 48) = S0
    δ(S0, 49) = S1
    δ(S1, 48) = S2
    δ(S1, 49) = S0
    δ(S2, 48) = S1
    δ(S2, 49) = S2
```

## Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Specific Package Tests

```bash
# Test the FSM library
go test -v ./fsm

# Test the mod-three example
go test -v ./examples
```

### Run Benchmarks

```bash
go test -bench=. -benchmem ./...
```

Expected output:
```
BenchmarkModThree-16              	 5000000	       230 ns/op	       0 B/op	       0 allocs/op
BenchmarkModThreeWithTrace-16     	 2000000	       650 ns/op	     128 B/op	       2 allocs/op
```

## Design Principles & Best Practices

This library demonstrates professional software engineering practices:

### 1. **Generic Programming**
- Type-safe generics for maximum flexibility
- Works with any comparable types (strings, ints, custom types)
- No runtime type assertions or reflection

### 2. **API Design**
- **Fluent Builder Pattern**: Intuitive, chainable API
- **Method Chaining**: All builder methods return `*Builder` for chaining
- **Clear Naming**: Methods directly map to formal FSA notation (Q, Σ, q₀, F, δ)
- **Multiple Construction Styles**: Direct construction or builder pattern

### 3. **Error Handling**
- Comprehensive validation with descriptive error messages
- Two build modes: `Build()` (returns error) and `MustBuild()` (panics)
- Input validation at every step
- Clear error messages referencing formal notation

### 4. **Testing**
- **100+ Test Cases**: Comprehensive coverage of all functionality
- **Table-Driven Tests**: Easy to add new test cases
- **Edge Cases**: Empty inputs, invalid symbols, undefined transitions
- **Benchmarks**: Performance testing included
- **Example Tests**: Real-world usage examples

### 5. **Documentation**
- **Package Documentation**: Complete godoc comments
- **Formal Definition**: Maps directly to mathematical notation
- **Usage Examples**: Multiple real-world examples
- **API Reference**: Complete method documentation

### 6. **Code Organization**
```
.
├── fsm/                    # Generic FSM library
│   ├── automaton.go       # Core automaton implementation
│   ├── builder.go         # Builder pattern
│   ├── automaton_test.go  # Library tests
│   └── builder_test.go    # Builder tests
├── examples/              # Example implementations
│   ├── modthree.go       # Mod-three implementation
│   └── modthree_test.go  # Mod-three tests
├── main.go               # Demo application
├── go.mod                # Module definition
└── README.md             # This file
```

### 7. **Performance**
- Zero allocations in hot paths
- Efficient map-based transition lookups
- Optional trace generation (only when needed)
- Benchmarked and optimized

## How It Works: Mod-Three Algorithm

The mod-three automaton demonstrates an elegant property of finite state machines: computing properties of arbitrarily large numbers without storing the entire number.

### Mathematical Foundation

For a binary number, each new bit effectively:
1. Doubles the current value (left shift)
2. Adds the new bit (0 or 1)

Therefore: `new_value = old_value * 2 + bit`

Taking modulo 3: `new_remainder = (old_remainder * 2 + bit) % 3`

### Why This Works

The FSM maintains only the **remainder** (0, 1, or 2), not the actual number. This means:
- **Constant Memory**: Only 3 states needed regardless of input length
- **No Overflow**: Never converts to decimal, avoiding integer overflow
- **Linear Time**: O(n) where n is the length of the binary string

### State Transition Logic

```
State S0 (remainder 0):
  - Input '0': (0 * 2 + 0) % 3 = 0 → Stay in S0
  - Input '1': (0 * 2 + 1) % 3 = 1 → Go to S1

State S1 (remainder 1):
  - Input '0': (1 * 2 + 0) % 3 = 2 → Go to S2
  - Input '1': (1 * 2 + 1) % 3 = 0 → Go to S0

State S2 (remainder 2):
  - Input '0': (2 * 2 + 0) % 3 = 1 → Go to S1
  - Input '1': (2 * 2 + 1) % 3 = 2 → Stay in S2
```

## Library Design Philosophy

### Why Generics?

The library uses Go generics to provide:
- **Type Safety**: Compile-time checking of state and symbol types
- **Flexibility**: Works with any comparable types
- **Performance**: No runtime type assertions or boxing
- **Clarity**: Clear type signatures in API

### Why Builder Pattern?

The builder pattern provides:
- **Readability**: Clear, declarative automaton construction
- **Validation**: Centralized validation in `Build()`
- **Flexibility**: Optional validation with `MustBuild()`
- **Discoverability**: IDE autocomplete guides usage

### Why Formal Notation?

Using formal FSA notation (Q, Σ, q₀, F, δ):
- **Academic Alignment**: Matches textbook definitions
- **Precision**: Unambiguous specification
- **Learning**: Helps users understand FSA theory
- **Documentation**: Self-documenting code

## Requirements

- Go 1.24.7 or later

## License

This is a demonstration project for educational purposes.

