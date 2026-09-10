package spot_test

import (
	"math"
	"strings"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

func TestEqualValues(t *testing.T) {
	type namedInt int
	type namedUint uint64
	cases := []struct {
		name     string
		expected any
		actual   any
		passes   bool
	}{
		{name: "strict equal strings", expected: "spot", actual: "spot", passes: true},
		{name: "strict equal collections", expected: []int{1, 2}, actual: []int{1, 2}, passes: true},
		{name: "different numeric types", expected: int(1), actual: int64(1), passes: true},
		{name: "named and built in numbers", expected: namedInt(2), actual: int8(2), passes: true},
		{name: "signed and unsigned", expected: int64(math.MaxInt64), actual: uint64(math.MaxInt64), passes: true},
		{name: "maximum unsigned values", expected: namedUint(math.MaxUint64), actual: uint64(math.MaxUint64), passes: true},
		{name: "exact floats", expected: float32(1.5), actual: float64(1.5), passes: true},
		{name: "different values", expected: 1, actual: 2, passes: false},
		{name: "negative and unsigned", expected: int64(-1), actual: uint64(math.MaxUint64), passes: false},
		{name: "large integer and rounded float", expected: int64(1<<53 + 1), actual: float64(1<<53 + 1), passes: false},
		{name: "float32 and float64 precision", expected: float32(0.1), actual: float64(0.1), passes: false},
		{name: "numeric and string", expected: 1, actual: "1", passes: false},
		{name: "not a number", expected: math.NaN(), actual: math.NaN(), passes: false},
		{name: "positive infinities", expected: float32(math.Inf(1)), actual: math.Inf(1), passes: true},
		{name: "negative infinities", expected: float32(math.Inf(-1)), actual: math.Inf(-1), passes: true},
		{name: "opposite infinities", expected: math.Inf(-1), actual: math.Inf(1), passes: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.EqualValues(stub, test.expected, test.actual)
			requireComparisonResult(t, stub, passed, test.passes)
		})
	}
}

func TestEqualValuesReportsExpectedActualAndFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.EqualValues(stub, 1, 2, "record %d", 7)
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "expected 1; got 2") {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func requireComparisonResult(t *testing.T, stub *mocks.TestingT, actual, expected bool) {
	t.Helper()
	if actual != expected {
		t.Fatalf("assertion result = %t; want %t", actual, expected)
	}
	wantFailures := 0
	if !expected {
		wantFailures = 1
	}
	if len(stub.Failures) != wantFailures {
		t.Fatalf("failures = %v; want %d", stub.Failures, wantFailures)
	}
}

func TestGreaterThanAcceptsLargerNumber(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.GreaterThan(stub, uint64(math.MaxUint64), int64(math.MaxInt64)) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterThanRejectsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, 2, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterThanRejectsNaN(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, math.NaN(), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterThanRejectsUnsupportedTypes(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, []int{2}, []int{1}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterThanOrEqualAcceptsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.GreaterThanOrEqual(stub, 2, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualAcceptsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.GreaterThanOrEqual(stub, 3, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualRejectsSmallerValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThanOrEqual(stub, 1, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessThanOrEqualAcceptsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.LessThanOrEqual(stub, 2, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessThanOrEqualAcceptsSmallerValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.LessThanOrEqual(stub, 1, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessThanOrEqualRejectsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThanOrEqual(stub, 3, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessThanAcceptsSmallerNumber(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.LessThan(stub, 1, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessThanRejectsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, 2, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessThanRejectsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, 3, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterThanRejectsStrings(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, "b", "a") {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanRejectsNonNumericValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, struct{}{}, struct{}{}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanRejectsInfinity(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, math.Inf(1), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanRejectsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThan(stub, nil, 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualRejectsStrings(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThanOrEqual(stub, "b", "a") {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualRejectsNonNumericValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThanOrEqual(stub, struct{}{}, struct{}{}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualRejectsInfinity(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThanOrEqual(stub, math.Inf(1), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestGreaterThanOrEqualRejectsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterThanOrEqual(stub, nil, 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanRejectsStrings(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, "b", "a") {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanRejectsNonNumericValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, struct{}{}, struct{}{}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanRejectsInfinity(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, math.Inf(1), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanRejectsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThan(stub, nil, 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanOrEqualRejectsStrings(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThanOrEqual(stub, "b", "a") {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanOrEqualRejectsNonNumericValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThanOrEqual(stub, struct{}{}, struct{}{}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanOrEqualRejectsInfinity(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThanOrEqual(stub, math.Inf(1), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}

func TestLessThanOrEqualRejectsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessThanOrEqual(stub, nil, 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}
