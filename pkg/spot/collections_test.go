package spot_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

func TestContainsSubstring(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, "hello world", "world")
	requireCollectionPass(t, stub, passed)
}

func TestContainsMissingSubstring(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, "hello", "world")
	requireCollectionFailure(t, stub, passed)
}

func TestContainsEmptySubstring(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, "", "")
	requireCollectionPass(t, stub, passed)
}

func TestContainsNamedStrings(t *testing.T) {
	type text string
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, text("hello"), text("ell"))
	requireCollectionPass(t, stub, passed)
}

func TestContainsInvalidStringTarget(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, "hello", 1)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsNilStringTarget(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, "hello", nil)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []int{1, 2}, 2)
	requireCollectionPass(t, stub, passed)
}

func TestContainsArray(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, [2]int{1, 2}, 1)
	requireCollectionPass(t, stub, passed)
}

func TestContainsMissingElement(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []int{1}, 2)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsDifferentNumericTypes(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []int{1}, int64(1))
	requireCollectionFailure(t, stub, passed)
}

func TestContainsNestedSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, [][]int{{1, 2}}, []int{1, 2})
	requireCollectionPass(t, stub, passed)
}

func TestContainsMapKey(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, map[string]int{"hello": 1}, "hello")
	requireCollectionPass(t, stub, passed)
}

func TestContainsNotMapValue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, map[string]int{"hello": 1}, 1)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsUncomparableTarget(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, map[string]int{"hello": 1}, []int{1})
	requireCollectionFailure(t, stub, passed)
}

func TestContainsNilElement(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []any{nil}, nil)
	requireCollectionPass(t, stub, passed)
}

func TestContainsNilSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []int(nil), 1)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsNilMap(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, map[string]int(nil), "hello")
	requireCollectionFailure(t, stub, passed)
}

func TestContainsNilContainer(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, nil, nil)
	requireCollectionFailure(t, stub, passed)
}

func TestContainsUnsupported(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, 42, 42)
	requireCollectionFailure(t, stub, passed)
}

func TestNotContains(t *testing.T) {
	type namedString string
	cases := []struct {
		name      string
		container any
		element   any
		passes    bool
	}{
		{name: "missing substring", container: "hello", element: "world", passes: true},
		{name: "present substring", container: "hello", element: "ell", passes: false},
		{name: "named strings", container: namedString("hello"), element: namedString("world"), passes: true},
		{name: "missing slice element", container: []int{1, 2}, element: 3, passes: true},
		{name: "present slice element", container: []int{1, 2}, element: 2, passes: false},
		{name: "missing array element", container: [2]int{1, 2}, element: 3, passes: true},
		{name: "present array element", container: [2]int{1, 2}, element: 1, passes: false},
		{name: "missing map key", container: map[string]int{"one": 1}, element: "two", passes: true},
		{name: "present map key", container: map[string]int{"one": 1}, element: "one", passes: false},
		{name: "map value is not a key", container: map[string]int{"one": 1}, element: 1, passes: true},
		{name: "nil slice", container: []int(nil), element: 1, passes: true},
		{name: "unsupported container", container: 1, element: 1, passes: false},
		{name: "invalid string target", container: "hello", element: 1, passes: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.NotContains(stub, test.container, test.element)
			if test.passes {
				requireCollectionPass(t, stub, passed)
				return
			}
			requireCollectionFailure(t, stub, passed)
		})
	}
}

func TestLen(t *testing.T) {
	type namedString string
	text := namedString("spot")
	items := []int{1, 2, 3}
	var nilPointer *string
	channel := make(chan int, 2)
	channel <- 1
	cases := []struct {
		name     string
		value    any
		expected int
		passes   bool
	}{
		{name: "string", value: "spot", expected: 4, passes: true},
		{name: "named string", value: namedString("spot"), expected: 4, passes: true},
		{name: "array", value: [2]int{1, 2}, expected: 2, passes: true},
		{name: "slice", value: items, expected: 3, passes: true},
		{name: "nil slice", value: []int(nil), expected: 0, passes: true},
		{name: "map", value: map[string]int{"one": 1}, expected: 1, passes: true},
		{name: "channel", value: channel, expected: 1, passes: true},
		{name: "pointer to string", value: &text, expected: 4, passes: true},
		{name: "pointer to slice", value: &items, expected: 3, passes: true},
		{name: "wrong length", value: "spot", expected: 3, passes: false},
		{name: "nil", value: nil, expected: 0, passes: false},
		{name: "typed nil pointer", value: nilPointer, expected: 0, passes: false},
		{name: "scalar", value: 4, expected: 0, passes: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.Len(stub, test.value, test.expected)
			if test.passes {
				requireCollectionPass(t, stub, passed)
				return
			}
			requireCollectionFailure(t, stub, passed)
		})
	}
}

func TestLenReportsExpectedAndActualLengths(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.Len(stub, []int{1, 2}, 3, "items for %s", "record")
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "expected length 3; got 2") {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
	if !strings.HasSuffix(stub.Failures[0], ": items for record") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func TestNotContainsIncludesFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.NotContains(stub, []int{1}, 1, "record %d", 7)
	requireCollectionFailure(t, stub, passed)
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func TestElementsMatchReordered(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1, 2, 3}, []int{3, 1, 2})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchDuplicates(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1, 1, 2}, []int{2, 1, 1})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchWrongDuplicateCounts(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1, 1, 2}, []int{1, 2, 2})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchMissing(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1, 2}, []int{1})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchExtra(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1}, []int{1, 2})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchArrayAndSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, [2]int{1, 2}, []int{2, 1})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchInterfaceSlice(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1, 2}, []any{2, 1})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchDifferentNumericTypes(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1}, []int64{1})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchNestedValues(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []any{[]int{1}, map[string]int{"a": 2}}, []any{map[string]int{"a": 2}, []int{1}})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchNilElements(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []any{nil, 1}, []any{1, nil})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchNilAndEmpty(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int(nil), []int{})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchUntypedNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, nil, []int{})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchBothNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, nil, nil)
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchNilAndNonempty(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, nil, []int{1})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchNonemptyAndNil(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1}, nil)
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchEmptyArrays(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, [0]int{}, [0]string{})
	requireCollectionPass(t, stub, passed)
}

func TestElementsMatchStringsRejected(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, "", "")
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchMapsRejected(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, map[string]int{}, map[string]int{})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchScalarsRejected(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, 0, 0)
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchPointerRejected(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, new([]int), []int{})
	requireCollectionFailure(t, stub, passed)
}

func TestElementsMatchDoesNotModifyInputs(t *testing.T) {
	expected, actual := []int{3, 1, 2}, []int{2, 3, 1}
	spot.ElementsMatch(t, expected, actual)
	if !reflect.DeepEqual(expected, []int{3, 1, 2}) || !reflect.DeepEqual(actual, []int{2, 3, 1}) {
		t.Fatal("assertion modified its inputs")
	}
}

func TestElementsMatchReportsUnmatchedDuplicates(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.ElementsMatch(stub, []int{1, 1, 2}, []int{1, 2, 2})
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "missing: []interface {}{1}; extra: []interface {}{2}") {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
}

func TestContainsIncludesFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.Contains(stub, []int{1}, 2, "record %d", 7)
	requireCollectionFailure(t, stub, passed)
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func TestElementsMatchIncludesFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	passed := spot.ElementsMatch(stub, []int{1}, []int{2}, "record %d", 7)
	requireCollectionFailure(t, stub, passed)
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func requireCollectionPass(t *testing.T, stub *mocks.TestingT, passed bool) {
	t.Helper()
	if !passed {
		t.Fatal("expected assertion to pass")
	}
	if len(stub.Failures) != 0 {
		t.Fatalf("unexpected failures: %v", stub.Failures)
	}
}

func requireCollectionFailure(t *testing.T, stub *mocks.TestingT, passed bool) {
	t.Helper()
	if passed {
		t.Fatal("expected assertion to fail")
	}
	if len(stub.Failures) != 1 {
		t.Fatalf("expected one fatal failure; got %v", stub.Failures)
	}
}
