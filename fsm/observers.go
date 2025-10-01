package fsm

import (
	"fmt"
	"sync"
	"time"
)

// LoggingObserver logs all state changes and events.
type LoggingObserver[Q State, S Symbol] struct {
	logger func(string)
}

// NewLoggingObserver creates a new logging observer.
func NewLoggingObserver[Q State, S Symbol](logger func(string)) *LoggingObserver[Q, S] {
	return &LoggingObserver[Q, S]{
		logger: logger,
	}
}

// OnStateChange logs state transitions.
func (o *LoggingObserver[Q, S]) OnStateChange(from Q, symbol S, to Q) {
	if o.logger != nil {
		o.logger(fmt.Sprintf("State transition: %v --(%v)--> %v", from, symbol, to))
	}
}

// OnInputProcessed logs input processing completion.
func (o *LoggingObserver[Q, S]) OnInputProcessed(input []S, accepted bool) {
	if o.logger != nil {
		status := "REJECTED"
		if accepted {
			status = "ACCEPTED"
		}
		o.logger(fmt.Sprintf("Input processed: %v -> %s", input, status))
	}
}

// OnError logs errors.
func (o *LoggingObserver[Q, S]) OnError(err error) {
	if o.logger != nil {
		o.logger(fmt.Sprintf("Error occurred: %v", err))
	}
}

// MetricsObserver collects performance metrics.
type MetricsObserver[Q State, S Symbol] struct {
	transitionCount     int64
	inputsProcessed     int64
	totalProcessingTime int64
	mutex               sync.RWMutex
}

// NewMetricsObserver creates a new metrics observer.
func NewMetricsObserver[Q State, S Symbol]() *MetricsObserver[Q, S] {
	return &MetricsObserver[Q, S]{}
}

// OnStateChange increments transition count.
func (o *MetricsObserver[Q, S]) OnStateChange(from Q, symbol S, to Q) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.transitionCount++
}

// OnInputProcessed increments input processing count.
func (o *MetricsObserver[Q, S]) OnInputProcessed(input []S, accepted bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.inputsProcessed++
}

// OnError does nothing for metrics observer.
func (o *MetricsObserver[Q, S]) OnError(err error) {
	// Metrics observer doesn't handle errors
}

// GetTransitionCount returns the number of transitions.
func (o *MetricsObserver[Q, S]) GetTransitionCount() int64 {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	return o.transitionCount
}

// GetInputsProcessedCount returns the number of inputs processed.
func (o *MetricsObserver[Q, S]) GetInputsProcessedCount() int64 {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	return o.inputsProcessed
}

// RecordProcessingTime records processing time.
func (o *MetricsObserver[Q, S]) RecordProcessingTime(duration int64) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.totalProcessingTime += duration
}

// GetAverageProcessingTime returns average processing time.
func (o *MetricsObserver[Q, S]) GetAverageProcessingTime() float64 {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	if o.inputsProcessed == 0 {
		return 0
	}
	return float64(o.totalProcessingTime) / float64(o.inputsProcessed)
}

// DebugObserver provides detailed debugging information.
type DebugObserver[Q State, S Symbol] struct {
	transitions []TransitionEvent[Q, S]
	inputs      []InputEvent[S]
	errors      []error
	mutex       sync.RWMutex
}

// TransitionEvent represents a state transition event.
type TransitionEvent[Q State, S Symbol] struct {
	Timestamp time.Time
	From      Q
	Symbol    S
	To        Q
}

// InputEvent represents an input processing event.
type InputEvent[S Symbol] struct {
	Timestamp time.Time
	Input     []S
	Accepted  bool
}

// NewDebugObserver creates a new debug observer.
func NewDebugObserver[Q State, S Symbol]() *DebugObserver[Q, S] {
	return &DebugObserver[Q, S]{
		transitions: make([]TransitionEvent[Q, S], 0),
		inputs:      make([]InputEvent[S], 0),
		errors:      make([]error, 0),
	}
}

// OnStateChange records state transition.
func (o *DebugObserver[Q, S]) OnStateChange(from Q, symbol S, to Q) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.transitions = append(o.transitions, TransitionEvent[Q, S]{
		Timestamp: time.Now(),
		From:      from,
		Symbol:    symbol,
		To:        to,
	})
}

// OnInputProcessed records input processing.
func (o *DebugObserver[Q, S]) OnInputProcessed(input []S, accepted bool) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.inputs = append(o.inputs, InputEvent[S]{
		Timestamp: time.Now(),
		Input:     input,
		Accepted:  accepted,
	})
}

// OnError records error.
func (o *DebugObserver[Q, S]) OnError(err error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.errors = append(o.errors, err)
}

// GetTransitions returns all recorded transitions.
func (o *DebugObserver[Q, S]) GetTransitions() []TransitionEvent[Q, S] {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	result := make([]TransitionEvent[Q, S], len(o.transitions))
	copy(result, o.transitions)
	return result
}

// GetInputs returns all recorded input events.
func (o *DebugObserver[Q, S]) GetInputs() []InputEvent[S] {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	result := make([]InputEvent[S], len(o.inputs))
	copy(result, o.inputs)
	return result
}

// GetErrors returns all recorded errors.
func (o *DebugObserver[Q, S]) GetErrors() []error {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	result := make([]error, len(o.errors))
	copy(result, o.errors)
	return result
}

// Clear clears all recorded events.
func (o *DebugObserver[Q, S]) Clear() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.transitions = o.transitions[:0]
	o.inputs = o.inputs[:0]
	o.errors = o.errors[:0]
}

// CompositeObserver combines multiple observers.
type CompositeObserver[Q State, S Symbol] struct {
	observers []Observer[Q, S]
	mutex     sync.RWMutex
}

// NewCompositeObserver creates a new composite observer.
func NewCompositeObserver[Q State, S Symbol](observers ...Observer[Q, S]) *CompositeObserver[Q, S] {
	return &CompositeObserver[Q, S]{
		observers: observers,
	}
}

// AddObserver adds an observer to the composite.
func (o *CompositeObserver[Q, S]) AddObserver(observer Observer[Q, S]) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.observers = append(o.observers, observer)
}

// RemoveObserver removes an observer from the composite.
func (o *CompositeObserver[Q, S]) RemoveObserver(observer Observer[Q, S]) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	for i, obs := range o.observers {
		if obs == observer {
			o.observers = append(o.observers[:i], o.observers[i+1:]...)
			break
		}
	}
}

// OnStateChange notifies all observers of state change.
func (o *CompositeObserver[Q, S]) OnStateChange(from Q, symbol S, to Q) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	for _, observer := range o.observers {
		observer.OnStateChange(from, symbol, to)
	}
}

// OnInputProcessed notifies all observers of input processing.
func (o *CompositeObserver[Q, S]) OnInputProcessed(input []S, accepted bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	for _, observer := range o.observers {
		observer.OnInputProcessed(input, accepted)
	}
}

// OnError notifies all observers of errors.
func (o *CompositeObserver[Q, S]) OnError(err error) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	for _, observer := range o.observers {
		observer.OnError(err)
	}
}

// ObservableAutomaton wraps an automaton with observer support.
type ObservableAutomaton[Q State, S Symbol] struct {
	automaton Automaton[Q, S]
	observers []Observer[Q, S]
	mutex     sync.RWMutex
}

// NewObservableAutomaton creates a new observable automaton.
func NewObservableAutomaton[Q State, S Symbol](automaton Automaton[Q, S]) *ObservableAutomaton[Q, S] {
	return &ObservableAutomaton[Q, S]{
		automaton: automaton,
		observers: make([]Observer[Q, S], 0),
	}
}

// AddObserver adds an observer.
func (oa *ObservableAutomaton[Q, S]) AddObserver(observer Observer[Q, S]) {
	oa.mutex.Lock()
	defer oa.mutex.Unlock()
	oa.observers = append(oa.observers, observer)
}

// RemoveObserver removes an observer.
func (oa *ObservableAutomaton[Q, S]) RemoveObserver(observer Observer[Q, S]) {
	oa.mutex.Lock()
	defer oa.mutex.Unlock()
	for i, obs := range oa.observers {
		if obs == observer {
			oa.observers = append(oa.observers[:i], oa.observers[i+1:]...)
			break
		}
	}
}

// NotifyObservers is a placeholder for interface compliance.
func (oa *ObservableAutomaton[Q, S]) NotifyObservers() {
	// This method is part of the interface but not used in this implementation
}

// Delegate methods to wrapped automaton
func (oa *ObservableAutomaton[Q, S]) GetInitialState() Q {
	return oa.automaton.GetInitialState()
}

func (oa *ObservableAutomaton[Q, S]) GetCurrentState() Q {
	return oa.automaton.GetCurrentState()
}

func (oa *ObservableAutomaton[Q, S]) Reset() {
	oa.automaton.Reset()
}

func (oa *ObservableAutomaton[Q, S]) IsAcceptingState(state Q) bool {
	return oa.automaton.IsAcceptingState(state)
}

func (oa *ObservableAutomaton[Q, S]) IsCurrentStateAccepting() bool {
	return oa.automaton.IsCurrentStateAccepting()
}

func (oa *ObservableAutomaton[Q, S]) Step(symbol S) (Q, error) {
	from := oa.automaton.GetCurrentState()
	to, err := oa.automaton.Step(symbol)

	if err != nil {
		oa.notifyError(err)
		return to, err
	}

	oa.notifyStateChange(from, symbol, to)
	return to, nil
}

func (oa *ObservableAutomaton[Q, S]) ProcessInput(input []S) (bool, error) {
	accepted, err := oa.automaton.ProcessInput(input)

	if err != nil {
		oa.notifyError(err)
	} else {
		oa.notifyInputProcessed(input, accepted)
	}

	return accepted, err
}

func (oa *ObservableAutomaton[Q, S]) ProcessInputWithTrace(input []S) ([]Q, bool, error) {
	trace, accepted, err := oa.automaton.ProcessInputWithTrace(input)

	if err != nil {
		oa.notifyError(err)
	} else {
		oa.notifyInputProcessed(input, accepted)
	}

	return trace, accepted, err
}

func (oa *ObservableAutomaton[Q, S]) Validate() error {
	return oa.automaton.Validate()
}

func (oa *ObservableAutomaton[Q, S]) String() string {
	return oa.automaton.String()
}

// Helper methods for notifying observers
func (oa *ObservableAutomaton[Q, S]) notifyStateChange(from Q, symbol S, to Q) {
	oa.mutex.RLock()
	defer oa.mutex.RUnlock()
	for _, observer := range oa.observers {
		observer.OnStateChange(from, symbol, to)
	}
}

func (oa *ObservableAutomaton[Q, S]) notifyInputProcessed(input []S, accepted bool) {
	oa.mutex.RLock()
	defer oa.mutex.RUnlock()
	for _, observer := range oa.observers {
		observer.OnInputProcessed(input, accepted)
	}
}

func (oa *ObservableAutomaton[Q, S]) notifyError(err error) {
	oa.mutex.RLock()
	defer oa.mutex.RUnlock()
	for _, observer := range oa.observers {
		observer.OnError(err)
	}
}
