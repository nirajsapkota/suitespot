package spot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

func TestEventuallyChecksImmediately(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	calls := 0
	passed := spot.Eventually(stub, func() bool {
		calls++
		return true
	}, time.Second, time.Hour)

	if !passed || len(stub.Failures) != 0 {
		t.Fatalf("passed = %t, failures = %v; want a pass", passed, stub.Failures)
	}
	if calls != 1 {
		t.Fatalf("condition calls = %d; want 1", calls)
	}
}

func TestEventuallyRetriesUntilConditionPasses(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	calls := 0
	passed := spot.Eventually(stub, func() bool {
		calls++
		return calls == 3
	}, 250*time.Millisecond, time.Millisecond)

	if !passed || len(stub.Failures) != 0 {
		t.Fatalf("passed = %t, failures = %v; want a pass", passed, stub.Failures)
	}
	if calls != 3 {
		t.Fatalf("condition calls = %d; want 3", calls)
	}
}

func TestEventuallyFailsOnceAfterTimeout(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Eventually(stub, func() bool { return false }, 10*time.Millisecond, time.Millisecond)

	if passed {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "within 10ms") {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestEventuallyValidatesInputs(t *testing.T) {
	cases := []struct {
		name       string
		condition  func() bool
		waitFor    time.Duration
		tick       time.Duration
		diagnostic string
	}{
		{name: "nil condition", condition: nil, waitFor: time.Second, tick: time.Millisecond, diagnostic: "condition; got nil"},
		{name: "zero wait", condition: func() bool { return true }, waitFor: 0, tick: time.Millisecond, diagnostic: "positive wait time"},
		{name: "negative wait", condition: func() bool { return true }, waitFor: -time.Second, tick: time.Millisecond, diagnostic: "positive wait time"},
		{name: "zero tick", condition: func() bool { return true }, waitFor: time.Second, tick: 0, diagnostic: "positive tick"},
		{name: "negative tick", condition: func() bool { return true }, waitFor: time.Second, tick: -time.Millisecond, diagnostic: "positive tick"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.Eventually(stub, test.condition, test.waitFor, test.tick)
			if passed {
				t.Fatal("expected assertion to fail")
			}
			if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], test.diagnostic) {
				t.Fatalf("unexpected failures: %v", stub.Failures)
			}
		})
	}
}

func TestEventuallyIncludesFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.Eventually(stub, func() bool { return false }, 5*time.Millisecond, time.Millisecond, "record %d", 7)
	if len(stub.Failures) != 1 || !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}
