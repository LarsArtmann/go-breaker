package breaker

import (
	"errors"
	"sync"
	"testing"
	"time"

	errorfamily "github.com/larsartmann/go-error-family"
)

func transientErr() error {
	return errorfamily.NewTransient("test.transient", "transient failure")
}

func infrastructureErr() error {
	return errorfamily.NewInfrastructure("test.infra", "infrastructure failure")
}

func rejectionErr() error {
	return errorfamily.NewRejection("test.rejection", "rejected")
}

func conflictErr() error {
	return errorfamily.NewConflict("test.conflict", "conflict")
}

func TestExecutePassesThroughWhenClosed(t *testing.T) {
	b := New(Config{})

	called := false

	err := b.Execute(func() error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("closed breaker must pass nil errors through, got %v", err)
	}

	if !called {
		t.Fatal("closed breaker must call fn")
	}

	if b.State() != StateClosed {
		t.Fatalf("state = %s, want closed", b.State())
	}
}

func TestThresholdConsecutiveFailuresOpen(t *testing.T) {
	b := New(Config{FailureThreshold: 3})

	for range 2 {
		if err := b.Execute(transientErr); !errors.Is(err, ErrOpen) && err == nil {
			t.Fatalf("unexpected nil error before threshold")
		}
	}

	if b.State() != StateClosed {
		t.Fatalf("two failures under threshold 3 must stay closed, got %s", b.State())
	}

	_ = b.Execute(transientErr)

	if b.State() != StateOpen {
		t.Fatalf("third consecutive failure must open, got %s", b.State())
	}
}

func TestOpenRefusesCheaply(t *testing.T) {
	b := New(Config{FailureThreshold: 1})

	_ = b.Execute(transientErr)

	called := false

	err := b.Execute(func() error {
		called = true
		return nil
	})

	if !errors.Is(err, ErrOpen) {
		t.Fatalf("open breaker must refuse with ErrOpen, got %v", err)
	}

	if called {
		t.Fatal("open breaker must not call fn")
	}
}

func TestSuccessResetsFailureCount(t *testing.T) {
	b := New(Config{FailureThreshold: 2})

	_ = b.Execute(transientErr)
	b.Record(nil)
	_ = b.Execute(transientErr)

	if b.State() != StateClosed {
		t.Fatalf("interleaved success must reset the failure count, got %s", b.State())
	}
}

func TestResetTimeoutMovesToHalfOpen(t *testing.T) {
	now := time.Date(2026, 11, 21, 15, 0, 0, 0, time.UTC)

	b := New(Config{
		FailureThreshold: 1,
		ResetTimeout:     time.Minute,
		Now:              func() time.Time { return now },
	})

	_ = b.Execute(transientErr)

	if b.State() != StateOpen {
		t.Fatalf("state = %s, want open", b.State())
	}

	now = now.Add(time.Minute)

	if b.State() != StateHalfOpen {
		t.Fatalf("after reset timeout state = %s, want half_open", b.State())
	}

	if err := b.Allow(); err != nil {
		t.Fatalf("half-open breaker must allow a trial, got %v", err)
	}
}

func TestHalfOpenClosesAfterSuccesses(t *testing.T) {
	now := time.Date(2026, 11, 21, 15, 0, 0, 0, time.UTC)

	b := New(Config{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		ResetTimeout:     time.Minute,
		Now:              func() time.Time { return now },
	})

	_ = b.Execute(transientErr)

	now = now.Add(2 * time.Minute)

	_ = b.Execute(func() error { return nil })
	_ = b.Execute(func() error { return nil })

	if b.State() != StateClosed {
		t.Fatalf("two half-open successes must close, got %s", b.State())
	}
}

func TestHalfOpenFailureReopensImmediately(t *testing.T) {
	now := time.Date(2026, 11, 21, 15, 0, 0, 0, time.UTC)

	b := New(Config{
		FailureThreshold: 1,
		ResetTimeout:     time.Minute,
		Now:              func() time.Time { return now },
	})

	_ = b.Execute(transientErr)

	now = now.Add(2 * time.Minute)

	if err := b.Allow(); err != nil {
		t.Fatalf("half-open trial must be allowed, got %v", err)
	}

	b.Record(transientErr())

	if b.State() != StateOpen {
		t.Fatalf("half-open failure must re-open immediately, got %s", b.State())
	}
}

func TestDeterministicFamiliesNeverTrip(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"rejection", rejectionErr()},
		{"conflict", conflictErr()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{FailureThreshold: 1})

			err := b.Execute(func() error { return tt.err })

			if !errors.Is(err, tt.err) {
				t.Fatalf("Execute must return fn's error unmodified, got %v", err)
			}

			if b.State() != StateClosed {
				t.Fatalf("deterministic failure must never trip the breaker, got %s", b.State())
			}
		})
	}
}

func TestInfrastructureAndTransientTrip(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"transient", transientErr()},
		{"infrastructure", infrastructureErr()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(Config{FailureThreshold: 1})

			_ = b.Execute(func() error { return tt.err })

			if b.State() != StateOpen {
				t.Fatalf("%s failure must trip, got %s", tt.name, b.State())
			}
		})
	}
}

func TestPlainErrorTripsAsTransient(t *testing.T) {
	b := New(Config{FailureThreshold: 1})

	_ = b.Execute(func() error { return errors.New("boom") })

	if b.State() != StateOpen {
		t.Fatalf("unclassified errors default to Transient (error-family fail-open), must trip, got %s", b.State())
	}
}

func TestErrOpenIsTransientFamily(t *testing.T) {
	if !errorfamily.IsRetryable(ErrOpen) {
		t.Fatal("ErrOpen must be Transient so outer retry loops back off, not abort")
	}
}

func TestExecuteReturnsFnErrorUnwrapped(t *testing.T) {
	b := New(Config{})

	sentinel := errors.New("boom")

	err := b.Execute(func() error { return sentinel })

	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute must return fn's error unmodified, got %v", err)
	}
}

func TestConcurrentUseStaysConsistent(t *testing.T) {
	b := New(Config{FailureThreshold: 10, SuccessThreshold: 2})

	var wg sync.WaitGroup

	for range 50 {

		wg.Go(func() {

			_ = b.Execute(func() error { return transientErr() })
			b.Record(nil)
			_ = b.State()
			_ = b.Allow()
		})
	}

	wg.Wait()

	if state := b.State(); state != StateClosed && state != StateOpen && state != StateHalfOpen {
		t.Fatalf("state %s is not a valid state", state)
	}
}
