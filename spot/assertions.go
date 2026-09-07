package spot

import (
	"fmt"
	"reflect"
	"strings"
)

type TestingT interface {
	Helper()
	Fatal(args ...any)
}

func True(t TestingT, value bool, message ...any) bool {
	t.Helper()
	return check(t, value, "expected true", message)
}

func False(t TestingT, value bool, message ...any) bool {
	t.Helper()
	return check(t, !value, "expected false", message)
}

func Nil(t TestingT, value any, message ...any) bool {
	t.Helper()
	return check(t, isNil(value), fmt.Sprintf("expected nil; got %#v", value), message)
}

func NotNil(t TestingT, value any, message ...any) bool {
	t.Helper()
	return check(t, !isNil(value), "expected a non-nil value", message)
}

func Err(t TestingT, err error, message ...any) bool {
	t.Helper()
	return check(t, err != nil, "expected an error", message)
}

func NoErr(t TestingT, err error, message ...any) bool {
	t.Helper()
	return check(t, err == nil, fmt.Sprintf("expected no error; got %v", err), message)
}

func Equal(t TestingT, expected, actual any, message ...any) bool {
	t.Helper()
	return check(t, reflect.DeepEqual(expected, actual), fmt.Sprintf("expected %#v (%T); got %#v (%T)", expected, expected, actual, actual), message)
}

func NotEqual(t TestingT, expected, actual any, message ...any) bool {
	t.Helper()
	return check(t, !reflect.DeepEqual(expected, actual), fmt.Sprintf("expected values to differ; both are %#v", actual), message)
}

func ErrorContains(t TestingT, err error, text string, message ...any) bool {
	t.Helper()
	return check(t, err != nil && strings.Contains(err.Error(), text), fmt.Sprintf("expected error containing %q; got %v", text, err), message)
}

func check(t TestingT, passed bool, failure string, message []any) bool {
	t.Helper()
	if passed {
		return true
	}
	if len(message) != 0 {
		failure += ": " + formatMessage(message)
	}
	t.Fatal(failure)
	return false
}

func formatMessage(message []any) string {
	if format, ok := message[0].(string); ok && len(message) > 1 {
		return fmt.Sprintf(format, message[1:]...)
	}
	return fmt.Sprint(message...)
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	case reflect.UnsafePointer:
		return v.Pointer() == 0
	default:
		return false
	}
}
