package spot

import (
	"fmt"
	"time"
)

// Eventually checks condition immediately and at each tick until it succeeds
// or the wait time expires. The condition must return without blocking.
func Eventually(t TestingT, condition func() bool, waitFor, tick time.Duration, message ...any) bool {
	t.Helper()
	if condition == nil {
		return check(t, false, "expected a condition; got nil", message)
	}
	if waitFor <= 0 {
		return check(t, false, fmt.Sprintf("expected a positive wait time; got %s", waitFor), message)
	}
	if tick <= 0 {
		return check(t, false, fmt.Sprintf("expected a positive tick; got %s", tick), message)
	}

	timeout := time.NewTimer(waitFor)
	defer timeout.Stop()
	if condition() {
		return true
	}

	interval := time.NewTicker(tick)
	defer interval.Stop()
	for {
		select {
		case <-interval.C:
			if condition() {
				return true
			}
		case <-timeout.C:
			return check(t, false, fmt.Sprintf("condition was not satisfied within %s", waitFor), message)
		}
	}
}
