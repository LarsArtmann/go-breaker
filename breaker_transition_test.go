package breaker

import (
	"testing"
	"time"
)

// OnTransition observer tests: the Config.OnTransition callback must fire
// exactly once per committed state change, in machine order
// (Closed -> Open -> HalfOpen -> Closed, and HalfOpen -> Open on failure).

func TestOnTransitionFiresInOrderAcrossLifecycle(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 21, 15, 0, 0, 0, time.UTC)

	type move struct{ from, to State }

	var moves []move

	b := New(Config{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		ResetTimeout:     time.Minute,
		Now:              func() time.Time { return now },
		OnTransition:     func(from, to State) { moves = append(moves, move{from, to}) },
	})

	_ = b.Execute(transientErr)

	if len(moves) != 0 {
		t.Fatalf("below threshold there is no transition, got %v", moves)
	}

	_ = b.Execute(transientErr)

	if len(moves) != 1 || moves[0].from != StateClosed || moves[0].to != StateOpen {
		t.Fatalf("threshold must fire closed->open, got %v", moves)
	}

	now = now.Add(2 * time.Minute)

	if err := b.Allow(); err != nil {
		t.Fatalf("post-timeout trial must be allowed, got %v", err)
	}

	if len(moves) != 2 || moves[1].from != StateOpen || moves[1].to != StateHalfOpen {
		t.Fatalf("Allow after timeout must fire open->half_open, got %v", moves)
	}

	_ = b.Execute(func() error { return nil })

	if len(moves) != 2 {
		t.Fatalf("first half-open success must not transition yet, got %v", moves)
	}

	_ = b.Execute(func() error { return nil })

	if len(moves) != 3 || moves[2].from != StateHalfOpen || moves[2].to != StateClosed {
		t.Fatalf("second half-open success must fire half_open->closed, got %v", moves)
	}
}

func TestOnTransitionHalfOpenFailureFiresReopen(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 21, 15, 0, 0, 0, time.UTC)

	var tos []State

	b := New(Config{
		FailureThreshold: 1,
		ResetTimeout:     time.Minute,
		Now:              func() time.Time { return now },
		OnTransition:     func(_, to State) { tos = append(tos, to) },
	})

	_ = b.Execute(transientErr)

	now = now.Add(2 * time.Minute)

	if err := b.Allow(); err != nil {
		t.Fatalf("half-open trial must be allowed, got %v", err)
	}

	b.Record(transientErr())

	want := []State{StateOpen, StateHalfOpen, StateOpen}
	if len(tos) != len(want) {
		t.Fatalf("transitions = %v, want %v", tos, want)
	}

	for i := range want {
		if tos[i] != want[i] {
			t.Fatalf("transitions = %v, want %v", tos, want)
		}
	}
}
