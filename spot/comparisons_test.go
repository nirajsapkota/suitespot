package spot_test

import (
	"math"
	"testing"

	"github.com/nirajsapkota/suitespot/mocks"
	"github.com/nirajsapkota/suitespot/spot"
)

func TestGreaterAcceptsLargerNumber(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.Greater(stub, uint64(math.MaxUint64), int64(math.MaxInt64)) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterRejectsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Greater(stub, 2, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterRejectsNaN(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Greater(stub, math.NaN(), 1) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterRejectsUnsupportedTypes(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Greater(stub, []int{2}, []int{1}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestGreaterOrEqualAcceptsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.GreaterOrEqual(stub, 2, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterOrEqualAcceptsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.GreaterOrEqual(stub, 3, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestGreaterOrEqualRejectsSmallerValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.GreaterOrEqual(stub, 1, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessOrEqualAcceptsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.LessOrEqual(stub, 2, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessOrEqualAcceptsSmallerValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.LessOrEqual(stub, 1, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessOrEqualRejectsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.LessOrEqual(stub, 3, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessAcceptsSmallerNumber(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.Less(stub, 1, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestLessRejectsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Less(stub, 2, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestLessRejectsGreaterValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Less(stub, 3, 2) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}
