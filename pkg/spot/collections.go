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
