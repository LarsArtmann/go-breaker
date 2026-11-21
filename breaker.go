// Package breaker is a dependency-light, error-family-aware circuit
// breaker: the sibling of go-retry for the FAILURE side of resilience.
//
// Retry answers "should I try again?"; the breaker answers "should I
// even try?" — after enough consecutive trip-worthy failures a
// downstream is presumed unhealthy and calls are refused cheaply until
// the cooldown passes and trial successes prove recovery.
//
// Error-family aware by default: only Transient and Infrastructure
// failures count toward opening the circuit — a Rejection (bad input),
// Conflict (state mismatch), or Corruption (integrity) failure is
// deterministic, and breaking the circuit would not make the next
// identical call succeed. A nil error always counts as a success.
//
// The states form a strict machine:
//
//	Closed --threshold consecutive trip-worthy failures--> Open
//	Open   --ResetTimeout elapsed, next Allow------------> HalfOpen
//	HalfOpen --SuccessThreshold successes----------------> Closed
//	HalfOpen --one trip-worthy failure--------------------> Open
package breaker

import (
	"sync"
	"time"

	errorfamily "github.com/larsartmann/go-error-family"
)

// State is the circuit's position in the Closed/Open/HalfOpen machine.
type State string

const (
	// StateClosed is healthy: calls pass through and failures count.
	StateClosed State = "closed"
	// StateOpen is tripped: calls are refused until the reset timeout
	// elapses.
	StateOpen State = "open"
	// StateHalfOpen is probing: trial calls pass; successes close the
	// circuit, a trip-worthy failure re-opens it.
	StateHalfOpen State = "half_open"
)

const (
	defaultFailureThreshold = 5
	defaultResetTimeout     = 30 * time.Second
	defaultSuccessThreshold = 2

	codeOpen = "breaker.open"
)

// ErrOpen is returned by [Breaker.Allow] and [Breaker.Execute] while the
// circuit is open. It is classified Transient on purpose: an open breaker
// is by nature a temporary condition, so outer retry loops back off and
// re-attempt rather than treating the refusal as terminal.
var ErrOpen error = errorfamily.NewTransient(codeOpen, "circuit breaker is open")

// Config configures a [Breaker]. The zero value is usable and means
// [DefaultConfig].
type Config struct {
	// FailureThreshold is the number of consecutive trip-worthy failures
	// that opens the circuit. Default: 5.
	FailureThreshold int

	// ResetTimeout is how long the circuit stays open before the next
	// Allow moves it to half-open probing. Default: 30s.
	ResetTimeout time.Duration

	// SuccessThreshold is the number of consecutive successes in
	// half-open that close the circuit again. Default: 2.
	SuccessThreshold int

	// TripOn decides whether an error counts toward opening the circuit.
	// Default: Transient and Infrastructure families (error-family
	// aware); deterministic families never trip. Override for
	// domain-specific tripping.
	TripOn func(err error) bool

	// Now is the clock; inject for deterministic tests. Default:
	// time.Now.
	Now func() time.Time

	// OnTransition, when set, is invoked after every state change with
	// the from/to states. Observers run while the breaker's lock is held,
	// so they are serialized in transition order — but they must be cheap,
	// non-blocking, and must never call back into the same Breaker (that
	// would deadlock). Use for metrics and logging, not for control flow.
	OnTransition func(from, to State)
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		FailureThreshold: defaultFailureThreshold,
		ResetTimeout:     defaultResetTimeout,
		SuccessThreshold: defaultSuccessThreshold,
		TripOn:           tripsByFamily,
		Now:              time.Now,
	}
}

// tripsByFamily is the default trip predicate: only failures whose
// error-family is Transient or Infrastructure count. A nil error is a
// success and never reaches the predicate.
func tripsByFamily(err error) bool {
	family := errorfamily.Classify(err)

	return family == errorfamily.Transient || family == errorfamily.Infrastructure
}

// Breaker is a concurrency-safe circuit breaker.
type Breaker struct {
	mu               sync.Mutex
	failureThreshold int
	resetTimeout     time.Duration
	successThreshold int
	tripOn           func(err error) bool
	now              func() time.Time

	state           State
	consecFailures  int
	consecSuccesses int
	openedAt        time.Time
	onTransition    func(from, to State)
}

// New creates a breaker, filling zero Config fields with defaults.
func New(cfg Config) *Breaker {
	def := DefaultConfig()

	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = def.FailureThreshold
	}

	if cfg.ResetTimeout <= 0 {
		cfg.ResetTimeout = def.ResetTimeout
	}

	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = def.SuccessThreshold
	}

	if cfg.TripOn == nil {
		cfg.TripOn = def.TripOn
	}

	if cfg.Now == nil {
		cfg.Now = def.Now
	}

	return &Breaker{
		failureThreshold: cfg.FailureThreshold,
		resetTimeout:     cfg.ResetTimeout,
		successThreshold: cfg.SuccessThreshold,
		tripOn:           cfg.TripOn,
		now:              cfg.Now,
		onTransition:     cfg.OnTransition,
		state:            StateClosed,
	}
}

// State returns the circuit's current state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.currentState()
}

// Allow reports whether a call may proceed, moving an open circuit whose
// reset timeout has elapsed to half-open. Returns [ErrOpen] while the
// circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.resetTimeout {
		b.transition(StateHalfOpen)
		b.consecSuccesses = 0
	}

	if b.state == StateOpen {
		return ErrOpen
	}

	return nil
}

// Record reports the outcome of an allowed call: nil is a success, an
// error is classified by the trip predicate. Successes reset the closed
// circuit's failure count and close a half-open circuit after
// SuccessThreshold consecutive successes; trip-worthy failures count
// toward opening and re-open a half-open circuit immediately.
func (b *Breaker) Record(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err == nil {
		b.recordSuccess()
		return
	}

	if !b.tripOn(err) {
		return
	}

	b.recordFailure()
}

// Execute runs fn under the circuit: it refuses while open ([ErrOpen]),
// otherwise it calls fn once and records the outcome. fn's error is
// returned unmodified (recorded, never wrapped).
func (b *Breaker) Execute(fn func() error) error {
	if err := b.Allow(); err != nil {
		return err
	}

	err := fn()
	b.Record(err)

	return err
}

func (b *Breaker) currentState() State {
	if b.state == StateOpen && b.now().Sub(b.openedAt) >= b.resetTimeout {
		return StateHalfOpen
	}

	return b.state
}

func (b *Breaker) recordSuccess() {
	b.consecFailures = 0

	if b.state == StateHalfOpen {
		b.consecSuccesses++
		if b.consecSuccesses >= b.successThreshold {
			b.transition(StateClosed)
			b.consecSuccesses = 0
		}

		return
	}

	b.consecSuccesses = 0
}

func (b *Breaker) recordFailure() {
	b.consecSuccesses = 0
	b.consecFailures++

	if b.state == StateHalfOpen || b.consecFailures >= b.failureThreshold {
		b.transition(StateOpen)
		b.openedAt = b.now()
		b.consecFailures = 0
	}
}

// transition moves the circuit to newState and fires the OnTransition
// observer. Callers must hold b.mu; the observer runs under the lock so
// events arrive in transition order.
func (b *Breaker) transition(newState State) {
	oldState := b.state
	if oldState == newState {
		return
	}

	b.state = newState

	if b.onTransition != nil {
		b.onTransition(oldState, newState)
	}
}
