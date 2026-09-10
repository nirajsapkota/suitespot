package spot

import (
	"fmt"
	"reflect"
	"strings"
)

// Contains checks for a substring, a slice or array element, or a map key.
// Elements and keys use deep equality. String containers require string targets.
func Contains(t TestingT, container, element any, message ...any) bool {
	t.Helper()
	found, err := containsElement(container, element)
	if err != nil {
		return check(t, false, err.Error(), message)
	}
	return check(t, found, fmt.Sprintf("expected %#v to contain %#v", container, element), message)
}

// NotContains checks that a string excludes a substring, a slice or array
// excludes an element, or a map excludes a key.
func NotContains(t TestingT, container, element any, message ...any) bool {
	t.Helper()
	found, err := containsElement(container, element)
	if err != nil {
		return check(t, false, err.Error(), message)
	}
	return check(t, !found, fmt.Sprintf("expected %#v not to contain %#v", container, element), message)
}

// Empty checks for nil, a zero scalar or struct, or a zero-length collection.
func Empty(t TestingT, value any, message ...any) bool {
	t.Helper()
	return check(t, isEmpty(value), fmt.Sprintf("expected %#v to be empty", value), message)
}

// NotEmpty checks that value is not empty.
func NotEmpty(t TestingT, value any, message ...any) bool {
	t.Helper()
	return check(t, !isEmpty(value), fmt.Sprintf("expected %#v not to be empty", value), message)
}

// Len checks the length of a string, array, slice, map, or channel.
func Len(t TestingT, value any, expected int, message ...any) bool {
	t.Helper()
	actual, valid := valueLength(value)
	if !valid {
		return check(t, false, fmt.Sprintf("expected a string, array, slice, map, or channel; got %T", value), message)
	}
	return check(t, actual == expected,
		fmt.Sprintf("expected length %d; got %d for %#v", expected, actual, value), message)
}

func containsElement(container, element any) (bool, error) {
	value := reflect.ValueOf(container)
	switch value.Kind() {
	case reflect.String:
		return containsSubstring(value.String(), element)
	case reflect.Map:
		return containsMapKey(value, element), nil
	case reflect.Array, reflect.Slice:
		return containsListElement(value, element), nil
	default:
		return false, fmt.Errorf("expected a string, array, slice, or map; got %T", container)
	}
}

func containsSubstring(text string, element any) (bool, error) {
	target := reflect.ValueOf(element)
	if target.Kind() != reflect.String {
		return false, fmt.Errorf("expected a string target; got %T", element)
	}
	return strings.Contains(text, target.String()), nil
}

func containsMapKey(value reflect.Value, element any) bool {
	keys := value.MapRange()
	for keys.Next() {
		if reflect.DeepEqual(keys.Key().Interface(), element) {
			return true
		}
	}
	return false
}

func containsListElement(value reflect.Value, element any) bool {
	for index := 0; index < value.Len(); index++ {
		if reflect.DeepEqual(value.Index(index).Interface(), element) {
			return true
		}
	}
	return false
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	actual, valid := indirectValue(value)
	if !valid {
		return isNil(value)
	}
	switch actual.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return actual.Len() == 0
	default:
		return actual.IsZero()
	}
}

func valueLength(value any) (int, bool) {
	actual, valid := indirectValue(value)
	if !valid {
		return 0, false
	}
	switch actual.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return actual.Len(), true
	default:
		return 0, false
	}
}

func indirectValue(value any) (reflect.Value, bool) {
	actual := reflect.ValueOf(value)
	seen := make(map[reference]struct{})
	for actual.IsValid() && (actual.Kind() == reflect.Interface || actual.Kind() == reflect.Pointer) {
		if actual.IsNil() {
			return reflect.Value{}, false
		}
		if actual.Kind() == reflect.Pointer {
			current := reference{actual.Type(), actual.Pointer()}
			if _, found := seen[current]; found {
				return reflect.Value{}, false
			}
			seen[current] = struct{}{}
		}
		actual = actual.Elem()
	}
	return actual, actual.IsValid()
}

type reference struct {
	typeOf  reflect.Type
	pointer uintptr
}

// ElementsMatch checks that arrays or slices contain deeply equal elements with
// equal duplicate counts, regardless of order. Nil is treated as an empty list.
// Inputs are not modified.
func ElementsMatch(t TestingT, expected, actual any, message ...any) bool {
	t.Helper()
	expectedList, expectedValid := collectionValue(expected)
	actualList, actualValid := collectionValue(actual)
	if !expectedValid || !actualValid {
		return check(t, false, fmt.Sprintf("expected arrays, slices, or nil; got %T and %T", expected, actual), message)
	}
	missing, extra := unmatchedElements(expectedList, actualList)
	return check(t, len(missing) == 0 && len(extra) == 0,
		fmt.Sprintf("expected matching elements; missing: %#v; extra: %#v", missing, extra), message)
}

func collectionValue(value any) (reflect.Value, bool) {
	if value == nil {
		return reflect.ValueOf([]any(nil)), true
	}
	collection := reflect.ValueOf(value)
	return collection, collection.Kind() == reflect.Array || collection.Kind() == reflect.Slice
}

func unmatchedElements(expected, actual reflect.Value) (missing, extra []any) {
	matched := make([]bool, actual.Len())
	for index := 0; index < expected.Len(); index++ {
		element := expected.Index(index).Interface()
		match := findUnmatchedElement(actual, element, matched)
		if match < 0 {
			missing = append(missing, element)
			continue
		}
		matched[match] = true
	}
	return missing, remainingElements(actual, matched)
}

func findUnmatchedElement(actual reflect.Value, element any, matched []bool) int {
	for index := 0; index < actual.Len(); index++ {
		if !matched[index] && reflect.DeepEqual(element, actual.Index(index).Interface()) {
			return index
		}
	}
	return -1
}

func remainingElements(actual reflect.Value, matched []bool) []any {
	var remaining []any
	for index, used := range matched {
		if !used {
			remaining = append(remaining, actual.Index(index).Interface())
		}
	}
	return remaining
}
