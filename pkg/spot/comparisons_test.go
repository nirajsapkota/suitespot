package spot_test

import (
	"math"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

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
