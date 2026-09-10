package spot_test

import (
	"strings"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

type emptyRecord struct {
	name string
}

func TestEmpty(t *testing.T) {
	type namedString string
	var nilPointer *int
	zero := 0
	one := 1
	var cycle any
	cycle = &cycle
	var nilFunction func()
	nonemptyChannel := make(chan int, 1)
	nonemptyChannel <- 1
	cases := []struct {
		name   string
		value  any
		passes bool
	}{
		{name: "nil", value: nil, passes: true},
		{name: "typed nil pointer", value: nilPointer, passes: true},
		{name: "typed nil slice", value: []int(nil), passes: true},
		{name: "typed nil map", value: map[string]int(nil), passes: true},
		{name: "typed nil channel", value: (chan int)(nil), passes: true},
		{name: "typed nil function", value: nilFunction, passes: true},
		{name: "empty string", value: "", passes: true},
		{name: "empty named string", value: namedString(""), passes: true},
		{name: "zero length array", value: [0]int{}, passes: true},
		{name: "empty slice", value: []int{}, passes: true},
		{name: "empty map", value: map[string]int{}, passes: true},
		{name: "empty channel", value: make(chan int), passes: true},
		{name: "zero number", value: 0, passes: true},
		{name: "false", value: false, passes: true},
		{name: "zero struct", value: emptyRecord{}, passes: true},
		{name: "pointer to zero", value: &zero, passes: true},
		{name: "pointer to nil slice", value: new([]int), passes: true},
		{name: "nonempty string", value: "spot", passes: false},
		{name: "nonzero number", value: 1, passes: false},
		{name: "true", value: true, passes: false},
		{name: "nonzero struct", value: emptyRecord{name: "spot"}, passes: false},
		{name: "array with zero element", value: [1]int{}, passes: false},
		{name: "nonempty slice", value: []int{0}, passes: false},
		{name: "nonempty channel", value: nonemptyChannel, passes: false},
		{name: "nonempty function", value: func() {}, passes: false},
		{name: "pointer to nonzero", value: &one, passes: false},
		{name: "cyclic interface", value: cycle, passes: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.Empty(stub, test.value)
			requireEmptinessResult(t, stub, passed, test.passes)
		})
	}
}

func TestNotEmpty(t *testing.T) {
	zero := 0
	one := 1
	cases := []struct {
		name   string
		value  any
		passes bool
	}{
		{name: "nil", value: nil, passes: false},
		{name: "empty string", value: "", passes: false},
		{name: "zero length array", value: [0]int{}, passes: false},
		{name: "empty slice", value: []int{}, passes: false},
		{name: "empty map", value: map[string]int{}, passes: false},
		{name: "zero number", value: 0, passes: false},
		{name: "false", value: false, passes: false},
		{name: "zero struct", value: emptyRecord{}, passes: false},
		{name: "pointer to zero", value: &zero, passes: false},
		{name: "nonempty string", value: "spot", passes: true},
		{name: "nonzero number", value: 1, passes: true},
		{name: "true", value: true, passes: true},
		{name: "nonzero struct", value: emptyRecord{name: "spot"}, passes: true},
		{name: "array with zero element", value: [1]int{}, passes: true},
		{name: "nonempty slice", value: []int{0}, passes: true},
		{name: "pointer to nonzero", value: &one, passes: true},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.NotEmpty(stub, test.value)
			requireEmptinessResult(t, stub, passed, test.passes)
		})
	}
}

func TestEmptyIncludesValueAndFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.Empty(stub, 1, "record %d", 7)
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "expected 1 to be empty") {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func TestNotEmptyIncludesValueAndFormattedMessage(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.NotEmpty(stub, nil, "record %d", 7)
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], "expected <nil> not to be empty") {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func requireEmptinessResult(t *testing.T, stub *mocks.TestingT, actual, expected bool) {
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
