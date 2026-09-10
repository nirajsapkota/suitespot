package spot_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/spot"
)

func TestRegexp(t *testing.T) {
	type namedString string
	var nilPattern *regexp.Regexp
	cases := []struct {
		name       string
		expression any
		value      any
		passes     bool
	}{
		{name: "string pattern matches", expression: `^spot$`, value: "spot", passes: true},
		{name: "compiled pattern matches", expression: regexp.MustCompile(`sp.t`), value: "spot", passes: true},
		{name: "named pattern matches named value", expression: namedString(`^spot$`), value: namedString("spot"), passes: true},
		{name: "pattern does not match", expression: `^suite$`, value: "spot", passes: false},
		{name: "invalid pattern", expression: `[`, value: "spot", passes: false},
		{name: "nil compiled pattern", expression: nilPattern, value: "spot", passes: false},
		{name: "unsupported pattern", expression: 1, value: "spot", passes: false},
		{name: "nil pattern", expression: nil, value: "spot", passes: false},
		{name: "unsupported value", expression: `1`, value: 1, passes: false},
		{name: "nil value", expression: `.*`, value: nil, passes: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stub := &mocks.TestingT{T: t}
			passed := spot.Regexp(stub, test.expression, test.value)
			requireRegexpResult(t, stub, passed, test.passes)
		})
	}
}

func TestRegexpReportsInvalidPattern(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.Regexp(stub, `[`, "spot", "record %d", 7)
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], `invalid regular expression "["`) {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
	if !strings.HasSuffix(stub.Failures[0], ": record 7") {
		t.Fatalf("missing custom message: %v", stub.Failures)
	}
}

func TestRegexpReportsPatternAndValue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	spot.Regexp(stub, `^suite$`, "spot")
	if len(stub.Failures) != 1 || !strings.Contains(stub.Failures[0], `expected "spot" to match regular expression "^suite$"`) {
		t.Fatalf("unexpected diagnostic: %v", stub.Failures)
	}
}

func requireRegexpResult(t *testing.T, stub *mocks.TestingT, actual, expected bool) {
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
