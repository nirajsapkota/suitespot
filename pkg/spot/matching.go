package spot

import (
	"fmt"
	"reflect"
	"regexp"
)

// Regexp checks that a string is matched by a regular expression. Expression
// may be a string or a compiled *regexp.Regexp.
func Regexp(t TestingT, expression, value any, message ...any) bool {
	t.Helper()
	pattern, err := regularExpression(expression)
	if err != nil {
		return check(t, false, err.Error(), message)
	}
	text, valid := stringValue(value)
	if !valid {
		return check(t, false, fmt.Sprintf("expected a string value; got %T", value), message)
	}
	return check(t, pattern.MatchString(text),
		fmt.Sprintf("expected %#v to match regular expression %q", value, pattern.String()), message)
}

func regularExpression(expression any) (*regexp.Regexp, error) {
	if compiled, ok := expression.(*regexp.Regexp); ok {
		if compiled == nil {
			return nil, fmt.Errorf("expected a regular expression; got nil")
		}
		return compiled, nil
	}
	pattern, valid := stringValue(expression)
	if !valid {
		return nil, fmt.Errorf("expected a regular expression string or *regexp.Regexp; got %T", expression)
	}
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regular expression %q: %v", pattern, err)
	}
	return compiled, nil
}

func stringValue(value any) (string, bool) {
	actual := reflect.ValueOf(value)
	if !actual.IsValid() || actual.Kind() != reflect.String {
		return "", false
	}
	return actual.String(), true
}
