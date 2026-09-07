package suite

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type TestingT interface {
	Helper()
	Fatal(args ...any)
	Run(name string, test func(*testing.T)) bool
}

type Suite interface {
	SetupSuite(t *testing.T)
}

type testSetup interface {
	SetupTest(t *testing.T)
}

type testMethod struct {
	name string
	run  func(*testing.T)
}

type suiteDefinition struct {
	name     string
	instance Suite
	tests    []testMethod
}

func Run(t TestingT, suites ...Suite) {
	t.Helper()
	runSuites(t, false, suites...)
}

func RunParallel(t TestingT, suites ...Suite) {
	t.Helper()
	runSuites(t, true, suites...)
}

func runSuites(t TestingT, parallel bool, suites ...Suite) {
	t.Helper()
	if len(suites) == 0 {
		t.Fatal("suite: at least one suite is required")
		return
	}
	for _, instance := range suites {
		if err := runSuite(t, parallel, instance); err != nil {
			t.Fatal(err)
			return
		}
	}
}

func runSuite(t TestingT, parallel bool, instance Suite) error {
	t.Helper()
	definition, err := defineSuite(instance)
	if err != nil {
		return err
	}
	t.Run(definition.name, func(t *testing.T) {
		runSuiteTests(t, parallel, definition)
	})
	return nil
}

func defineSuite(instance Suite) (suiteDefinition, error) {
	name, err := suiteName(instance)
	if err != nil {
		return suiteDefinition{}, err
	}
	tests, err := findTests(instance)
	return suiteDefinition{name: name, instance: instance, tests: tests}, err
}

func suiteName(instance Suite) (string, error) {
	value, err := suiteStruct(instance)
	if err != nil {
		return "", err
	}
	name := value.Type().Name()
	if name == "" {
		return "", fmt.Errorf("suite: expected a named struct, got %T", instance)
	}
	return name, nil
}

func suiteStruct(instance Suite) (reflect.Value, error) {
	value := reflect.ValueOf(instance)
	if !isStructPointer(value) {
		return reflect.Value{}, fmt.Errorf("suite: expected a non-nil pointer to a named struct, got %T", instance)
	}
	return value.Elem(), nil
}

func isStructPointer(value reflect.Value) bool {
	return value.IsValid() &&
		value.Kind() == reflect.Pointer &&
		!value.IsNil() &&
		value.Elem().Kind() == reflect.Struct
}

func runSuiteTests(t *testing.T, parallel bool, definition suiteDefinition) {
	t.Helper()
	definition.instance.SetupSuite(t)
	for _, test := range definition.tests {
		runTest(t, parallel, definition.instance, test)
	}
}

func findTests(instance Suite) ([]testMethod, error) {
	value := reflect.ValueOf(instance)
	if err := validateTestSetup(value); err != nil {
		return nil, err
	}
	tests, err := collectTests(value)
	if err != nil {
		return nil, err
	}
	return tests, requireTests(instance, tests)
}

func validateTestSetup(value reflect.Value) error {
	if !value.MethodByName("SetupTest").IsValid() {
		return nil
	}
	_, err := methodFunction(value, "SetupTest")
	return err
}

func collectTests(value reflect.Value) ([]testMethod, error) {
	var tests []testMethod
	for _, name := range testMethodNames(value.Type()) {
		run, err := methodFunction(value, name)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testMethod{name: name, run: run})
	}
	return tests, nil
}

func testMethodNames(value reflect.Type) []string {
	var names []string
	for i := 0; i < value.NumMethod(); i++ {
		name := value.Method(i).Name
		if strings.HasPrefix(name, "Test") {
			names = append(names, name)
		}
	}
	return names
}

func methodFunction(value reflect.Value, name string) (func(*testing.T), error) {
	run, valid := value.MethodByName(name).Interface().(func(*testing.T))
	if !valid {
		return nil, fmt.Errorf("suite: %s.%s must have signature func(*testing.T)", value.Type(), name)
	}
	return run, nil
}

func requireTests(instance Suite, tests []testMethod) error {
	if len(tests) == 0 {
		return fmt.Errorf("suite: %T has no Test methods", instance)
	}
	return nil
}

func runTest(t *testing.T, parallel bool, instance Suite, test testMethod) {
	t.Helper()
	t.Run(test.name, func(t *testing.T) {
		runTestBody(t, parallel, instance, test)
	})
}

func runTestBody(t *testing.T, parallel bool, instance Suite, test testMethod) {
	t.Helper()
	if parallel {
		t.Parallel()
	}
	setupTest(t, instance)
	test.run(t)
}

func setupTest(t *testing.T, instance Suite) {
	t.Helper()
	if setup, ok := instance.(testSetup); ok {
		setup.SetupTest(t)
	}
}
