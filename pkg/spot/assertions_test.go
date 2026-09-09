package spot_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

func TestTrueAcceptsTrue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.True(stub, true) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestTrueRejectsFalse(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.True(stub, false) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected true") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestFalseAcceptsFalse(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.False(stub, false) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestFalseRejectsTrue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.False(stub, true) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected false") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestNilAcceptsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.Nil(stub, nil) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNilAcceptsTypedNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var value *int
	if !spot.Nil(stub, value) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNilAcceptsNilMap(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var value map[string]int
	if !spot.Nil(stub, value) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNilAcceptsNilSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var value []int
	if !spot.Nil(stub, value) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNilAcceptsNilFunction(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var value func()
	if !spot.Nil(stub, value) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNilRejectsEmptyNonNilSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Nil(stub, []int{}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestNilRejectsZero(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Nil(stub, 0) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestNotNilAcceptsValue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.NotNil(stub, new(int)) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNotNilRejectsTypedNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var value *int
	if spot.NotNil(stub, value) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestErrAcceptsError(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.Err(stub, errors.New("failure")) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestErrAcceptsTypedNilError(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var err *os.PathError
	if !spot.Err(stub, err) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestErrRejectsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Err(stub, nil) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected an error") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestNoErrAcceptsNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.NoErr(stub, nil) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNoErrRejectsError(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.NoErr(stub, errors.New("failure")) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "failure") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestNoErrRejectsTypedNilError(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	var err *os.PathError
	if spot.NoErr(stub, err) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestEqualComparesNestedValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.Equal(stub, map[string][]int{"a": {1, 2}}, map[string][]int{"a": {1, 2}}) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestEqualRejectsDifferentTypes(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Equal(stub, int(1), int64(1)) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestNotEqualAcceptsDifferentValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.NotEqual(stub, 1, 2) {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestNotEqualRejectsEqualValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.NotEqual(stub, []int{1}, []int{1}) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestErrorContainsAcceptsSubstring(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if !spot.ErrorContains(stub, errors.New("read failed"), "failed") {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func TestErrorContainsRejectsMissingSubstring(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.ErrorContains(stub, errors.New("read failed"), "write") {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
}

func TestAssertionFormatsOptionalMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.True(stub, false, "item %d missing", 7) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "item 7 missing") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestEqualFormatsOptionalMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	if spot.Equal(stub, 1, 2, "item %d", 7) {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one fatal failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "item 7") {
		t.Fatalf("unexpected diagnostic: %q", stub.Failures[0])
	}
}

func TestParallelAssertionsUseTheirOwnTest(t *testing.T) {
	for i := 0; i < 16; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			spot.True(stub, false, "test %d", i)
			want := fmt.Sprintf("expected true: test %d", i)
			if len(stub.Failures) != 1 || stub.Failures[0] != want {
				t.Fatalf("failures = %v; want %q", stub.Failures, want)
			}
		})
	}
}
