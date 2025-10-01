# Contributing to PolicyReporter-FSM

Thank you for your interest in contributing to PolicyReporter-FSM! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional, but recommended)
- Docker (for containerization)

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/PolicyReporter-FSM.git
   cd PolicyReporter-FSM
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/dsonic0912/PolicyReporter-FSM.git
   ```

4. Install development tools:
   ```bash
   make install-tools
   ```

5. Verify your setup:
   ```bash
   make dev
   ```

## Coding Standards

### Go Style Guide

We follow the standard Go style guide and additional conventions:

- Use `gofmt` for code formatting
- Follow effective Go practices
- Use meaningful variable and function names
- Write clear, concise comments
- Keep functions small and focused

### Code Organization

- **Package Structure**: Follow Go package conventions
- **File Naming**: Use descriptive names with underscores for multi-word files
- **Interface Design**: Prefer small, focused interfaces
- **Error Handling**: Always handle errors appropriately
- **Documentation**: All public APIs must have godoc comments

### Linting

We use `golangci-lint` with a comprehensive configuration. Run linting with:

```bash
make lint
```

Common linting rules:
- Line length should not exceed 120 characters
- Avoid naked returns in functions longer than 30 lines
- Use consistent naming conventions
- Avoid unnecessary complexity

## Testing Guidelines

### Test Coverage

- Maintain minimum 90% test coverage
- Write tests for all public APIs
- Include edge cases and error conditions
- Use table-driven tests where appropriate

### Test Types

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Benchmark Tests**: Performance testing for critical paths
4. **Property-Based Tests**: Test mathematical properties

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run benchmarks
make bench

# Generate coverage report
make coverage
```

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected ExpectedType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Submitting Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-new-processor`
- `fix/memory-leak-in-observer`
- `docs/update-api-documentation`
- `refactor/simplify-builder-interface`

### Commit Messages

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Maintenance tasks

Examples:
```
feat(fsm): add observer pattern support

Add observer interface and implementations for monitoring
state transitions and input processing events.

Closes #123
```

### Pull Request Process

1. **Create Feature Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**: Implement your feature or fix

3. **Run Tests**: Ensure all tests pass
   ```bash
   make ci
   ```

4. **Commit Changes**: Use conventional commit messages

5. **Push Branch**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**: 
   - Use a descriptive title
   - Include detailed description
   - Reference related issues
   - Add appropriate labels

### Pull Request Requirements

- [ ] All tests pass
- [ ] Code coverage maintained or improved
- [ ] Linting passes without errors
- [ ] Documentation updated (if applicable)
- [ ] CHANGELOG.md updated (for significant changes)
- [ ] Backward compatibility maintained (or breaking changes documented)

## Development Workflow

### Daily Development

```bash
# Start development
git checkout main
git pull upstream main
git checkout -b feature/my-feature

# Make changes and test frequently
make dev

# Before committing
make ci
```

### Code Review Checklist

**For Authors:**
- [ ] Self-review your code
- [ ] Run full test suite
- [ ] Update documentation
- [ ] Consider backward compatibility
- [ ] Add appropriate tests

**For Reviewers:**
- [ ] Code follows project conventions
- [ ] Tests are comprehensive
- [ ] Documentation is clear
- [ ] Performance implications considered
- [ ] Security implications considered

## Architecture Guidelines

### Interface Design

- Keep interfaces small and focused
- Use composition over inheritance
- Design for testability
- Consider future extensibility

### Error Handling

- Use structured error types
- Provide meaningful error messages
- Include context in errors
- Handle errors at appropriate levels

### Performance Considerations

- Profile critical paths
- Use benchmarks for performance-sensitive code
- Consider memory allocations
- Optimize for common use cases

## Release Process

### Versioning

We use Semantic Versioning (SemVer):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Release Checklist

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create release tag
4. Automated release via GitHub Actions
5. Update documentation
6. Announce release

## Getting Help

- **Issues**: Create GitHub issues for bugs or feature requests
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check existing documentation first
- **Code Examples**: Look at the examples directory

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to PolicyReporter-FSM!
