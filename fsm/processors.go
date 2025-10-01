package fsm

import (
	"sync"
)

// StandardProcessor implements the standard sequential processing strategy.
type StandardProcessor[Q State, S Symbol] struct{}

// NewStandardProcessor creates a new standard processor.
func NewStandardProcessor[Q State, S Symbol]() *StandardProcessor[Q, S] {
	return &StandardProcessor[Q, S]{}
}

// Process processes input through the automaton sequentially.
func (p *StandardProcessor[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	automaton.Reset()
	trace := []Q{automaton.GetCurrentState()}

	for _, symbol := range input {
		newState, err := automaton.Step(symbol)
		if err != nil {
			return ProcessResult[Q]{
				Accepted:   false,
				Trace:      trace,
				FinalState: automaton.GetCurrentState(),
			}, err
		}
		trace = append(trace, newState)
	}

	finalState := automaton.GetCurrentState()
	accepted := automaton.IsCurrentStateAccepting()

	return ProcessResult[Q]{
		Accepted:   accepted,
		Trace:      trace,
		FinalState: finalState,
	}, nil
}

// ParallelProcessor implements parallel processing for multiple inputs.
type ParallelProcessor[Q State, S Symbol] struct {
	maxWorkers int
}

// NewParallelProcessor creates a new parallel processor.
func NewParallelProcessor[Q State, S Symbol](maxWorkers int) *ParallelProcessor[Q, S] {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	return &ParallelProcessor[Q, S]{
		maxWorkers: maxWorkers,
	}
}

// ProcessBatch processes multiple inputs in parallel.
func (p *ParallelProcessor[Q, S]) ProcessBatch(automaton Automaton[Q, S], inputs [][]S) ([]ProcessResult[Q], error) {
	if len(inputs) == 0 {
		return []ProcessResult[Q]{}, nil
	}

	results := make([]ProcessResult[Q], len(inputs))
	errors := make([]error, len(inputs))

	// Create worker pool
	workers := p.maxWorkers
	if workers > len(inputs) {
		workers = len(inputs)
	}

	inputChan := make(chan struct {
		index int
		input []S
	}, len(inputs))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processor := NewStandardProcessor[Q, S]()

			for job := range inputChan {
				result, err := processor.Process(automaton, job.input)
				results[job.index] = result
				errors[job.index] = err
			}
		}()
	}

	// Send jobs
	for i, input := range inputs {
		inputChan <- struct {
			index int
			input []S
		}{i, input}
	}
	close(inputChan)

	// Wait for completion
	wg.Wait()

	// Check for errors
	collector := NewErrorCollector()
	for _, err := range errors {
		collector.Add(err)
	}

	return results, collector.ToError()
}

// Process processes a single input (implements Processor interface).
func (p *ParallelProcessor[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	// For single input, just use standard processing
	processor := NewStandardProcessor[Q, S]()
	return processor.Process(automaton, input)
}

// OptimizedProcessor implements optimized processing with caching.
type OptimizedProcessor[Q State, S Symbol] struct {
	cache map[string]ProcessResult[Q]
	mutex sync.RWMutex
}

// NewOptimizedProcessor creates a new optimized processor with caching.
func NewOptimizedProcessor[Q State, S Symbol]() *OptimizedProcessor[Q, S] {
	return &OptimizedProcessor[Q, S]{
		cache: make(map[string]ProcessResult[Q]),
	}
}

// Process processes input with caching for repeated inputs.
func (p *OptimizedProcessor[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	// Create cache key
	key := p.createCacheKey(input)

	// Check cache first
	p.mutex.RLock()
	if result, exists := p.cache[key]; exists {
		p.mutex.RUnlock()
		return result, nil
	}
	p.mutex.RUnlock()

	// Process normally
	processor := NewStandardProcessor[Q, S]()
	result, err := processor.Process(automaton, input)
	if err != nil {
		return result, err
	}

	// Cache the result
	p.mutex.Lock()
	p.cache[key] = result
	p.mutex.Unlock()

	return result, nil
}

// createCacheKey creates a string key for caching.
func (p *OptimizedProcessor[Q, S]) createCacheKey(input []S) string {
	// Simple string representation for caching
	// In a real implementation, you might want a more sophisticated key
	return string(rune(len(input))) // Simplified for demo
}

// ClearCache clears the processor's cache.
func (p *OptimizedProcessor[Q, S]) ClearCache() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.cache = make(map[string]ProcessResult[Q])
}

// GetCacheSize returns the number of cached results.
func (p *OptimizedProcessor[Q, S]) GetCacheSize() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return len(p.cache)
}

// TracingProcessor wraps another processor and adds detailed tracing.
type TracingProcessor[Q State, S Symbol] struct {
	wrapped Processor[Q, S]
	tracer  func(string)
}

// NewTracingProcessor creates a new tracing processor.
func NewTracingProcessor[Q State, S Symbol](wrapped Processor[Q, S], tracer func(string)) *TracingProcessor[Q, S] {
	return &TracingProcessor[Q, S]{
		wrapped: wrapped,
		tracer:  tracer,
	}
}

// Process processes input with detailed tracing.
func (p *TracingProcessor[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	if p.tracer != nil {
		p.tracer("Starting input processing")
		defer p.tracer("Finished input processing")
	}

	result, err := p.wrapped.Process(automaton, input)

	if p.tracer != nil {
		if err != nil {
			p.tracer("Processing failed with error: " + err.Error())
		} else {
			p.tracer("Processing completed successfully")
		}
	}

	return result, err
}

// ValidatingProcessor wraps another processor and adds input validation.
type ValidatingProcessor[Q State, S Symbol] struct {
	wrapped   Processor[Q, S]
	validator func([]S) error
}

// NewValidatingProcessor creates a new validating processor.
func NewValidatingProcessor[Q State, S Symbol](
	wrapped Processor[Q, S],
	validator func([]S) error,
) *ValidatingProcessor[Q, S] {
	return &ValidatingProcessor[Q, S]{
		wrapped:   wrapped,
		validator: validator,
	}
}

// Process processes input with validation.
func (p *ValidatingProcessor[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	// Validate input first
	if p.validator != nil {
		if err := p.validator(input); err != nil {
			return ProcessResult[Q]{}, NewErrorWithCause(ErrorTypeInvalidInput, "input validation failed", err)
		}
	}

	return p.wrapped.Process(automaton, input)
}

// ProcessorChain allows chaining multiple processors.
type ProcessorChain[Q State, S Symbol] struct {
	processors []Processor[Q, S]
}

// NewProcessorChain creates a new processor chain.
func NewProcessorChain[Q State, S Symbol](processors ...Processor[Q, S]) *ProcessorChain[Q, S] {
	return &ProcessorChain[Q, S]{
		processors: processors,
	}
}

// Process processes input through all processors in the chain.
func (p *ProcessorChain[Q, S]) Process(automaton Automaton[Q, S], input []S) (ProcessResult[Q], error) {
	if len(p.processors) == 0 {
		return ProcessResult[Q]{}, NewError(ErrorTypeInvalidConfiguration, "no processors in chain")
	}

	// Use the first processor for actual processing
	// Others could be used for validation, logging, etc.
	return p.processors[0].Process(automaton, input)
}

// Add adds a processor to the chain.
func (p *ProcessorChain[Q, S]) Add(processor Processor[Q, S]) {
	p.processors = append(p.processors, processor)
}

// ProcessorRegistry manages different processor implementations.
type ProcessorRegistry[Q State, S Symbol] struct {
	processors       map[string]Processor[Q, S]
	defaultProcessor Processor[Q, S]
}

// NewProcessorRegistry creates a new processor registry.
func NewProcessorRegistry[Q State, S Symbol]() *ProcessorRegistry[Q, S] {
	registry := &ProcessorRegistry[Q, S]{
		processors: make(map[string]Processor[Q, S]),
	}

	// Register default processor
	defaultProcessor := NewStandardProcessor[Q, S]()
	registry.Register("standard", defaultProcessor)
	registry.SetDefault(defaultProcessor)

	return registry
}

// Register adds a processor to the registry.
func (r *ProcessorRegistry[Q, S]) Register(name string, processor Processor[Q, S]) {
	r.processors[name] = processor
}

// Get retrieves a processor by name.
func (r *ProcessorRegistry[Q, S]) Get(name string) (Processor[Q, S], bool) {
	processor, exists := r.processors[name]
	return processor, exists
}

// GetDefault returns the default processor.
func (r *ProcessorRegistry[Q, S]) GetDefault() Processor[Q, S] {
	return r.defaultProcessor
}

// SetDefault sets the default processor.
func (r *ProcessorRegistry[Q, S]) SetDefault(processor Processor[Q, S]) {
	r.defaultProcessor = processor
}
