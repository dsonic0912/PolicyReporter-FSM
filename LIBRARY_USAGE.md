# FSM Library - Usage Guide for Developers

This guide shows how to use the generic FSM library to build your own finite state automata.

## Table of Contents

1. [Basic Concepts](#basic-concepts)
2. [Creating Your First Automaton](#creating-your-first-automaton)
3. [Common Patterns](#common-patterns)
4. [Advanced Examples](#advanced-examples)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

## Basic Concepts

### What is a Finite Automaton?

A finite automaton is a mathematical model of computation consisting of:
- **States**: A finite set of conditions the system can be in
- **Alphabet**: A finite set of input symbols
- **Transitions**: Rules for moving between states based on input
- **Initial State**: Where the automaton starts
- **Accepting States**: States that indicate successful completion

### The 5-Tuple Definition

Every automaton in this library is defined by **(Q, Σ, q₀, F, δ)**:

```go
fa := fsm.NewBuilder[StateType, SymbolType](q0).  // q₀: initial state
    WithStates(q1, q2, q3).                        // Q: set of states
    WithAlphabet(sym1, sym2).                      // Σ: alphabet
    WithAcceptingStates(qAccept).                  // F: accepting states
    WithTransition(from, symbol, to).              // δ: transition function
    Build()
```

## Creating Your First Automaton

### Example: Even/Odd Bit Counter

Let's build an automaton that accepts binary strings with an even number of 1s.

#### Step 1: Define the States

```go
// We need two states: even (0, 2, 4... ones) and odd (1, 3, 5... ones)
const (
    Even = "even"
    Odd  = "odd"
)
```

#### Step 2: Define the Alphabet

```go
// Binary alphabet: '0' and '1'
// We'll use rune type for characters
```

#### Step 3: Define Transitions

```
State: Even, Input: '0' → Stay in Even (0 doesn't change parity)
State: Even, Input: '1' → Go to Odd (one more 1)
State: Odd,  Input: '0' → Stay in Odd (0 doesn't change parity)
State: Odd,  Input: '1' → Go to Even (one more 1 makes it even)
```

#### Step 4: Build the Automaton

```go
package main

import (
    "fmt"
    "github.com/dsonic0912/PolicyReporter-FSM/fsm"
)

func main() {
    fa := fsm.NewBuilder[string, rune]("even").
        WithStates("even", "odd").
        WithAlphabet('0', '1').
        WithAcceptingStates("even").  // Accept even number of 1s
        WithTransitions(
            fsm.T("even", '0', "even"),
            fsm.T("even", '1', "odd"),
            fsm.T("odd", '0', "odd"),
            fsm.T("odd", '1', "even"),
        ).
        MustBuild()

    // Test it
    inputs := []string{"11", "101", "1111", "10101"}
    for _, input := range inputs {
        accepted, _ := fa.ProcessInput([]rune(input))
        fmt.Printf("%s: %v\n", input, accepted)
    }
}
```

Output:
```
11: true      (two 1s - even)
101: false    (three 1s - odd)
1111: true    (four 1s - even)
10101: false  (three 1s - odd)
```

## Common Patterns

### Pattern 1: String Suffix Checker

Accepts strings ending with a specific pattern.

```go
// Accepts strings ending in "ing"
fa := fsm.NewBuilder[int, rune](0).
    WithStates(0, 1, 2, 3).
    WithAlphabet('i', 'n', 'g', 'x'). // 'x' represents any other char
    WithAcceptingStates(3).
    WithTransitions(
        // State 0: no match yet
        fsm.T(0, 'i', 1),
        fsm.T(0, 'n', 0),
        fsm.T(0, 'g', 0),
        fsm.T(0, 'x', 0),
        
        // State 1: saw 'i'
        fsm.T(1, 'i', 1),
        fsm.T(1, 'n', 2),
        fsm.T(1, 'g', 0),
        fsm.T(1, 'x', 0),
        
        // State 2: saw "in"
        fsm.T(2, 'i', 1),
        fsm.T(2, 'n', 0),
        fsm.T(2, 'g', 3),
        fsm.T(2, 'x', 0),
        
        // State 3: saw "ing" - accepting
        fsm.T(3, 'i', 1),
        fsm.T(3, 'n', 0),
        fsm.T(3, 'g', 0),
        fsm.T(3, 'x', 0),
    ).
    MustBuild()
```

### Pattern 2: Divisibility Checker

Check if a number (in any base) is divisible by N.

```go
// Check if binary number is divisible by 5
// States represent remainder (0, 1, 2, 3, 4)
fa := fsm.NewBuilder[int, rune](0).
    WithStates(0, 1, 2, 3, 4).
    WithAlphabet('0', '1').
    WithAcceptingStates(0). // Divisible by 5 means remainder 0
    WithTransitions(
        // For each state s and input b: new_state = (s * 2 + b) % 5
        fsm.T(0, '0', 0), fsm.T(0, '1', 1),
        fsm.T(1, '0', 2), fsm.T(1, '1', 3),
        fsm.T(2, '0', 4), fsm.T(2, '1', 0),
        fsm.T(3, '0', 1), fsm.T(3, '1', 2),
        fsm.T(4, '0', 3), fsm.T(4, '1', 4),
    ).
    MustBuild()
```

### Pattern 3: State Machine with Custom Types

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
    WithAcceptingStates(Red, Yellow, Green). // All states valid
    WithTransitions(
        // Normal cycle
        fsm.T(Red, Timer, Green),
        fsm.T(Green, Timer, Yellow),
        fsm.T(Yellow, Timer, Red),
        
        // Emergency: always go to red
        fsm.T(Red, Emergency, Red),
        fsm.T(Green, Emergency, Red),
        fsm.T(Yellow, Emergency, Red),
    ).
    MustBuild()

// Simulate traffic light
events := []Event{Timer, Timer, Emergency, Timer}
for _, event := range events {
    fa.Step(event)
    fmt.Printf("Light is now: %s\n", fa.GetCurrentState())
}
```

## Advanced Examples

### Example 1: Email Validator (Simplified)

```go
type State int
const (
    Start State = iota
    Username
    AtSign
    Domain
    Dot
    TLD
)

fa := fsm.NewBuilder[State, rune](Start).
    WithStates(Start, Username, AtSign, Domain, Dot, TLD).
    WithAlphabet('a', '@', '.', 'x'). // 'a' = letter, 'x' = other
    WithAcceptingStates(TLD).
    WithTransitions(
        fsm.T(Start, 'a', Username),
        fsm.T(Username, 'a', Username),
        fsm.T(Username, '@', AtSign),
        fsm.T(AtSign, 'a', Domain),
        fsm.T(Domain, 'a', Domain),
        fsm.T(Domain, '.', Dot),
        fsm.T(Dot, 'a', TLD),
        fsm.T(TLD, 'a', TLD),
    ).
    MustBuild()
```

### Example 2: Balanced Parentheses (Limited Depth)

```go
// Accepts strings with balanced parentheses (max depth 3)
fa := fsm.NewBuilder[int, rune](0).
    WithStates(0, 1, 2, 3, -1). // -1 = error state
    WithAlphabet('(', ')').
    WithAcceptingStates(0). // Only balanced at depth 0
    WithTransitions(
        fsm.T(0, '(', 1), fsm.T(0, ')', -1),
        fsm.T(1, '(', 2), fsm.T(1, ')', 0),
        fsm.T(2, '(', 3), fsm.T(2, ')', 1),
        fsm.T(3, '(', -1), fsm.T(3, ')', 2),
        fsm.T(-1, '(', -1), fsm.T(-1, ')', -1),
    ).
    MustBuild()
```

### Example 3: Protocol State Machine

```go
type ProtocolState string
type Message string

const (
    Disconnected ProtocolState = "DISCONNECTED"
    Connecting   ProtocolState = "CONNECTING"
    Connected    ProtocolState = "CONNECTED"
    Authenticated ProtocolState = "AUTHENTICATED"
)

const (
    Connect Message = "CONNECT"
    ConnectAck Message = "CONNECT_ACK"
    Auth Message = "AUTH"
    AuthOk Message = "AUTH_OK"
    Disconnect Message = "DISCONNECT"
)

fa := fsm.NewBuilder[ProtocolState, Message](Disconnected).
    WithStates(Disconnected, Connecting, Connected, Authenticated).
    WithAlphabet(Connect, ConnectAck, Auth, AuthOk, Disconnect).
    WithAcceptingStates(Authenticated).
    WithTransitions(
        fsm.T(Disconnected, Connect, Connecting),
        fsm.T(Connecting, ConnectAck, Connected),
        fsm.T(Connected, Auth, Connected),
        fsm.T(Connected, AuthOk, Authenticated),
        fsm.T(Authenticated, Disconnect, Disconnected),
        fsm.T(Connected, Disconnect, Disconnected),
        fsm.T(Connecting, Disconnect, Disconnected),
    ).
    MustBuild()
```

## Best Practices

### 1. Use Meaningful Type Aliases

```go
// ✅ Good
type State string
type Input rune
fa := fsm.NewBuilder[State, Input]("start")

// ❌ Less clear
fa := fsm.NewBuilder[string, rune]("start")
```

### 2. Define Constants for States and Symbols

```go
// ✅ Good
const (
    StateIdle State = "IDLE"
    StateActive State = "ACTIVE"
)

// ❌ Avoid magic strings
fa.AddTransition("IDLE", 'x', "ACTIVE")
```

### 3. Validate Early

```go
// ✅ Good - validate during construction
fa, err := builder.Build()
if err != nil {
    log.Fatalf("Invalid automaton: %v", err)
}

// Or use MustBuild for static automata
fa := builder.MustBuild()
```

### 4. Document Your Automaton

```go
// NewEmailValidator creates an automaton that validates email format.
//
// Formal definition:
//   Q = {start, user, at, domain, dot, tld}
//   Σ = {letter, '@', '.'}
//   q₀ = start
//   F = {tld}
//   δ = defined by transitions below
func NewEmailValidator() *fsm.FiniteAutomaton[State, rune] {
    return fsm.NewBuilder[State, rune](Start).
        // ... transitions
        MustBuild()
}
```

## Troubleshooting

### Error: "initial state X is not in the set of states Q"

**Problem**: You forgot to add the initial state to the states set.

**Solution**:
```go
// ✅ Fix
fa := fsm.NewBuilder[string, rune]("q0").
    WithStates("q0", "q1", "q2"). // Include q0
    // ...
```

### Error: "no transition defined for state X with symbol Y"

**Problem**: Your automaton is incomplete - missing a transition.

**Solution**: Add the missing transition or make sure all possible inputs are handled:
```go
// Add error/sink state for unhandled inputs
fa := fsm.NewBuilder[string, rune]("q0").
    WithStates("q0", "q1", "error").
    WithAlphabet('a', 'b').
    WithTransitions(
        fsm.T("q0", 'a', "q1"),
        fsm.T("q0", 'b', "error"), // Handle 'b' from q0
        fsm.T("q1", 'a', "q1"),
        fsm.T("q1", 'b', "error"), // Handle 'b' from q1
        fsm.T("error", 'a', "error"),
        fsm.T("error", 'b', "error"),
    ).
    MustBuild()
```

### Error: "symbol not in alphabet"

**Problem**: You're trying to process a symbol that wasn't added to the alphabet.

**Solution**: Add all possible symbols to the alphabet:
```go
fa := fsm.NewBuilder[string, rune]("q0").
    WithAlphabet('0', '1', '2', '3', '4', '5', '6', '7', '8', '9').
    // ...
```

## Next Steps

- See [fsm/README.md](fsm/README.md) for complete API reference
- Check [examples/modthree.go](examples/modthree.go) for a complete working example
- Read the main [README.md](README.md) for project overview

## Getting Help

If you encounter issues:
1. Check that your automaton is complete (all states have transitions for all symbols)
2. Validate your automaton with `fa.Validate()`
3. Use `ProcessInputWithTrace()` to debug state transitions
4. Print the automaton with `fa.String()` to visualize the configuration

