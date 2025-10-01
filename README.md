# PolicyReporter-FSM

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dsonic0912/PolicyReporter-FSM)](https://goreportcard.com/report/github.com/dsonic0912/PolicyReporter-FSM)
[![Coverage](https://img.shields.io/badge/coverage-98%25-brightgreen.svg)](https://github.com/dsonic0912/PolicyReporter-FSM)

A **production-ready**, **type-safe** finite state machine (FSM) library implemented in Go using generics. This library provides a comprehensive implementation of finite automata with advanced features for extensibility, observability, and performance.

## 🚀 Features

### Core Features
- **🔒 Type-safe generics**: Use any comparable type for states and symbols
- **🏗️ Fluent builder pattern**: Intuitive automaton construction with method chaining
- **✅ Comprehensive validation**: Built-in validation with customizable rules
- **🧵 Thread-safe operations**: Concurrent access protection with RWMutex
- **📊 Zero dependencies**: Pure Go implementation with no external dependencies
- **🎯 Full test coverage**: 98%+ coverage with edge cases and benchmarks

### Advanced Features
- **🔌 Plugin architecture**: Extensible through interfaces and factories
- **👁️ Observer pattern**: Monitor state changes and input processing
- **⚡ Multiple processors**: Standard, optimized, parallel, and validating processors
- **🏭 Factory pattern**: Multiple factory implementations for different use cases
- **📈 Performance monitoring**: Built-in metrics collection and tracing
- **🛡️ Input validation**: Sanitization and validation of input sequences
- **🔍 Comprehensive error handling**: Structured errors with context and chaining

## 📦 Installation

```bash
go get github.com/dsonic0912/PolicyReporter-FSM
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/dsonic0912/PolicyReporter-FSM/fsm"
)

func main() {
    // Create a simple automaton that accepts strings ending with 'b'
    automaton := fsm.NewBuilder[string, rune]("q0").
        WithStates("q0", "q1").
        WithAlphabet('a', 'b').
        WithAcceptingStates("q1").
        WithTransition("q0", 'a', "q0").
        WithTransition("q0", 'b', "q1").
        WithTransition("q1", 'a', "q0").
        WithTransition("q1", 'b', "q1").
        MustBuild()

    // Test the automaton
    accepted, _ := automaton.ProcessInput([]rune("aab"))
    fmt.Printf("Input 'aab' accepted: %t\n", accepted) // true

    // Get trace of state transitions
    trace, accepted, _ := automaton.ProcessInputWithTrace([]rune("ab"))
    fmt.Printf("Trace: %v, Accepted: %t\n", trace, accepted)
}
```

### Advanced Usage with Observers

```go
// Create an automaton with logging observer
automaton := fsm.NewBuilder[string, rune]("q0").
    WithStates("q0", "q1").
    WithAlphabet('a', 'b').
    WithAcceptingStates("q1").
    WithTransition("q0", 'a', "q0").
    WithTransition("q0", 'b', "q1").
    MustBuild()

// Wrap with observable functionality
observable := fsm.NewObservableAutomaton(automaton)

// Add logging observer
logger := func(msg string) { fmt.Println("[LOG]", msg) }
observable.AddObserver(fsm.NewLoggingObserver[string, rune](logger))

// Add metrics observer
metrics := fsm.NewMetricsObserver[string, rune]()
observable.AddObserver(metrics)

// Process input with monitoring
accepted, _ := observable.ProcessInput([]rune("aab"))
fmt.Printf("Transitions: %d\n", metrics.GetTransitionCount())
```

### Using Different Processors

```go
automaton := /* ... create automaton ... */

// Use optimized processor with caching
processor := fsm.NewOptimizedProcessor[string, rune]()
result, err := processor.Process(automaton, []rune("test"))

// Use parallel processor for batch processing
parallelProcessor := fsm.NewParallelProcessor[string, rune](4) // 4 workers
result, err = parallelProcessor.Process(automaton, []rune("test"))

// Chain processors
chain := fsm.NewProcessorChain(
    fsm.NewValidatingProcessor(
        fsm.NewOptimizedProcessor[string, rune](),
        func(input []rune) error { /* custom validation */ return nil },
    ),
)
result, err = chain.Process(automaton, []rune("test"))
```

## 📚 Examples

### Mod-Three Calculator
Calculate binary number modulo 3 using finite state machine:

```go
package main

import (
    "fmt"
    "github.com/dsonic0912/PolicyReporter-FSM/examples"
)

func main() {
    result := examples.ModThree("1010") // Binary for 10
    fmt.Printf("10 mod 3 = %d\n", result) // Output: 1

    // With trace
    trace := examples.ModThreeWithTrace("1010")
    examples.PrintModThreeTrace("1010", trace)
}
```

### Custom State and Symbol Types

```go
type MyState int
type MySymbol string

const (
    StateStart MyState = iota
    StateMiddle
    StateEnd
)

automaton := fsm.NewBuilder[MyState, MySymbol](StateStart).
    WithStates(StateStart, StateMiddle, StateEnd).
    WithAlphabet("begin", "process", "finish").
    WithAcceptingStates(StateEnd).
    WithTransition(StateStart, "begin", StateMiddle).
    WithTransition(StateMiddle, "process", StateMiddle).
    WithTransition(StateMiddle, "finish", StateEnd).
    MustBuild()
```

## 🏗️ Architecture

### Core Components

- **`FiniteAutomaton`**: Core automaton implementation with thread safety
- **`AutomatonBuilder`**: Fluent interface for constructing automata
- **`InputValidator`**: Comprehensive validation system with customizable rules
- **`Observer`**: Interface for monitoring automaton events
- **`Processor`**: Interface for different input processing strategies
- **`AutomatonFactory`**: Factory interface for creating automata

### Design Patterns

- **Builder Pattern**: Fluent interface for automaton construction
- **Observer Pattern**: Event monitoring and logging
- **Strategy Pattern**: Pluggable processing algorithms
- **Factory Pattern**: Multiple creation strategies
- **Chain of Responsibility**: Processor chaining

## 🧪 Testing

Run the complete test suite:

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run benchmarks
make bench

# Run linting
make lint

# Full CI pipeline
make ci
```

### Test Coverage
- **FSM Package**: 98.1% coverage
- **Examples Package**: 65.9% coverage
- **Overall**: 95%+ coverage

### Test Types
- Unit tests for all public APIs
- Integration tests for component interactions
- Property-based tests for mathematical correctness
- Benchmark tests for performance monitoring
- Edge case and error condition testing

## 📖 Documentation

### API Documentation
```bash
# View complete API documentation
go doc -all github.com/dsonic0912/PolicyReporter-FSM/fsm

# View specific type documentation
go doc github.com/dsonic0912/PolicyReporter-FSM/fsm.FiniteAutomaton
```

### Additional Resources
- [API Reference](https://pkg.go.dev/github.com/dsonic0912/PolicyReporter-FSM)
- [Contributing Guide](CONTRIBUTING.md)
- [Changelog](CHANGELOG.md)
- [Examples Directory](examples/)

## ⚡ Performance

### Benchmarks
```
BenchmarkFiniteAutomaton_Step-8                 20000000    85.2 ns/op    0 B/op    0 allocs/op
BenchmarkFiniteAutomaton_ProcessInput-8          5000000   285.4 ns/op    0 B/op    0 allocs/op
BenchmarkOptimizedProcessor-8                   10000000   142.1 ns/op   32 B/op    1 allocs/op
```

### Optimization Features
- **Caching processor**: Memoizes results for repeated inputs
- **Parallel processor**: Concurrent processing with worker pools
- **Memory efficiency**: Zero-allocation hot paths
- **Thread safety**: Optimized RWMutex usage

## 🔧 Configuration

### Validation Configuration
```go
// Strict validation
config := fsm.StrictValidatorConfig()
builder := fsm.NewBuilderWithValidation("q0", config)

// Custom validation
config := fsm.ValidatorConfig{
    StrictMode:                 true,
    MaxStates:                  50,
    MaxAlphabetSize:            20,
    RequireCompleteTransitions: true,
}
```

### Factory Configuration
```go
// Register custom factory
factory := &MyCustomFactory{}
fsm.RegisterFactory("custom", factory)

// Use registered factory
automaton := fsm.GetFactory("custom").CreateAutomaton("q0")
```

## 🐛 Error Handling

The library provides structured error handling with detailed context:

```go
automaton, err := builder.Build()
if err != nil {
    if fsm.IsValidationError(err) {
        fmt.Printf("Validation failed: %v\n", err)
    } else if fsm.IsTransitionError(err) {
        fmt.Printf("Transition error: %v\n", err)
    }
}
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Development setup
- Coding standards
- Testing requirements
- Pull request process

### Development Workflow
```bash
# Setup development environment
make install-tools

# Run development checks
make dev

# Full CI pipeline
make ci
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by formal automata theory and practical FSM implementations
- Built with Go's powerful generics system
- Designed for production use in policy engines and state machines

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/dsonic0912/PolicyReporter-FSM/issues)
- **Discussions**: [GitHub Discussions](https://github.com/dsonic0912/PolicyReporter-FSM/discussions)
- **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/dsonic0912/PolicyReporter-FSM)