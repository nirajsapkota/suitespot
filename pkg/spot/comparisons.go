package spot

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
)

func GreaterThan(t TestingT, actual, other any, message ...any) bool {
	t.Helper()
	comparison, valid := compareNumbers(actual, other)
	return check(t, valid && comparison > 0, fmt.Sprintf("expected %#v to be greater than %#v", actual, other), message)
}

func GreaterThanOrEqual(t TestingT, actual, other any, message ...any) bool {
	t.Helper()
	comparison, valid := compareNumbers(actual, other)
	return check(t, valid && comparison >= 0, fmt.Sprintf("expected %#v to be greater than or equal to %#v", actual, other), message)
}

func LessThan(t TestingT, actual, other any, message ...any) bool {
	t.Helper()
	comparison, valid := compareNumbers(actual, other)
	return check(t, valid && comparison < 0, fmt.Sprintf("expected %#v to be less than %#v", actual, other), message)
}

func LessThanOrEqual(t TestingT, actual, other any, message ...any) bool {
	t.Helper()
	comparison, valid := compareNumbers(actual, other)
	return check(t, valid && comparison <= 0, fmt.Sprintf("expected %#v to be less than or equal to %#v", actual, other), message)
}

func compareNumbers(left, right any) (int, bool) {
	a, aOK := number(left)
	b, bOK := number(right)
	if !aOK || !bOK {
		return 0, false
	}
	return a.Cmp(b), true
}

func number(value any) (*big.Rat, bool) {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return new(big.Rat).SetInt64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return new(big.Rat).SetUint64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		if !math.IsInf(v.Float(), 0) && !math.IsNaN(v.Float()) {
			return new(big.Rat).SetFloat64(v.Float()), true
		}
	}
	return nil, false
}
