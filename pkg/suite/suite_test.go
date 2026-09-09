package suite_test

import (
	"path"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/nirajsapkota/suitespot/internal/mocks"
	"github.com/nirajsapkota/suitespot/pkg/suite"
)

func TestRunNoSuitesFails(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "suite: at least one suite is required") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "suite: at least one suite is required")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunInstanceIsNotValid(t *testing.T) {
	var instance suite.Suite
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected a non-nil pointer to a named struct") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "expected a non-nil pointer to a named struct")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunInstanceIsNotPointer(t *testing.T) {
	instance := valueSuite{}
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected a non-nil pointer to a named struct") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "expected a non-nil pointer to a named struct")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunInstanceIsNil(t *testing.T) {
	var instance *suiteProbe
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected a non-nil pointer to a named struct") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "expected a non-nil pointer to a named struct")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunInstanceIsNotStruct(t *testing.T) {
	instance := new(scalarSuite)
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected a non-nil pointer to a named struct") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "expected a non-nil pointer to a named struct")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunInstanceNameIsBlank(t *testing.T) {
	instance := &struct{ valueSuite }{}
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "expected a named struct") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "expected a named struct")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunSetupTestHasNoArguments(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &setupWithoutArguments{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], ".SetupTest must have signature func(*testing.T)") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], ".SetupTest must have signature func(*testing.T)")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunSetupTestHasWrongArgument(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &setupWithWrongArgument{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], ".SetupTest must have signature func(*testing.T)") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], ".SetupTest must have signature func(*testing.T)")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunSetupTestHasExtraArgument(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &setupWithExtraArgument{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], ".SetupTest must have signature func(*testing.T)") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], ".SetupTest must have signature func(*testing.T)")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunSetupTestHasReturnValue(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &setupWithReturnValue{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], ".SetupTest must have signature func(*testing.T)") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], ".SetupTest must have signature func(*testing.T)")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunTestMethodHasWrongSignature(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &testWithWrongSignature{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], ".TestValue must have signature func(*testing.T)") {
		t.Errorf("unexpected failure: %q", stub.Failures[0])
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunRequiresAtLeastOneFoundTestMethod(t *testing.T) {
	instance := &suiteWithoutTests{}
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, instance)

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "has no Test methods") {
		t.Errorf("failure = %q; want %q", stub.Failures[0], "has no Test methods")
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunStopsAfterInvalidSuite(t *testing.T) {
	stub := &mocks.TestingT{T: t}
	suite.Run(stub, &suiteWithoutTests{}, &suiteProbe{})

	if len(stub.Failures) != 1 {
		t.Fatalf("failures = %v; want one validation failure", stub.Failures)
	}
	if !strings.Contains(stub.Failures[0], "has no Test methods") {
		t.Errorf("unexpected failure: %q", stub.Failures[0])
	}
	if stub.Runs != 0 {
		t.Errorf("started %d suites after validation failed", stub.Runs)
	}
}

func TestRunOnlyCollectsTestsPrefixedCorrectly(t *testing.T) {
	instance := &suiteProbe{}
	suite.Run(t, instance)
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunWithoutSetupTest(t *testing.T) {
	instance := &suiteWithoutSetupTest{}
	suite.Run(t, instance)
	if got := instance.bodies.Load(); got != 2 {
		t.Errorf("executed bodies = %d; want 2", got)
	}
}

func TestRunCallsSetupSuiteExactlyOnce(t *testing.T) {
	instance := &suiteProbe{}
	instance.onBody = func(t *testing.T) {
		if got := instance.suiteSetups.Load(); got != 1 {
			t.Errorf("suite setups before test body = %d; want 1", got)
		}
	}
	suite.Run(t, instance)
	if got := instance.suiteSetups.Load(); got != 1 {
		t.Errorf("suite setups = %d; want 1", got)
	}
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunCallsSetupTestOncePerTest(t *testing.T) {
	instance := &suiteProbe{}
	suite.Run(t, instance)
	if got := instance.testSetups.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.testSetups = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunSingleSuite(t *testing.T) {
	instance := &suiteProbe{}
	suite.Run(t, instance)
	if got := instance.suiteSetups.Load(); got != 1 {
		t.Errorf("suite setups = %d; want 1", got)
	}
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunMultipleSuites(t *testing.T) {
	first, second := &suiteProbe{}, &suiteProbe{}
	suite.Run(t, first, second)
	if got := first.suiteSetups.Load(); got != 1 {
		t.Errorf("suite setups = %d; want 1", got)
	}
	if got := first.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("first.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
	if got := second.suiteSetups.Load(); got != 1 {
		t.Errorf("suite setups = %d; want 1", got)
	}
	if got := second.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("second.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunExecutesTestBodies(t *testing.T) {
	var executed atomic.Int32
	suite.Run(t, &suiteProbe{onBody: func(t *testing.T) { executed.Add(1) }})
	if got := executed.Load(); got != 2 {
		t.Errorf("executed bodies = %d; want 2", got)
	}
}

func TestRunExecutesTestBodiesSequentially(t *testing.T) {
	var events callLog
	suite.Run(t, &suiteProbe{onBody: func(t *testing.T) {
		name := path.Base(t.Name())
		events.add(name + " start")
		t.Run("nested", func(t *testing.T) {
			t.Parallel()
			events.add(name + " end")
		})
	}})
	want := []string{"TestAlpha start", "TestAlpha end", "TestBeta start", "TestBeta end"}
	if got := events.snapshot(); !slices.Equal(got, want) {
		t.Errorf("events = %v; want %v", got, want)
	}
}

func TestRunParallelExecutesTestBodiesConcurrently(t *testing.T) {
	requireParallelCapacity(t)
	instance := &suiteProbe{onBody: rendezvous()}
	suite.RunParallel(t, instance)
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunParallelExecutesSetupTestConcurrently(t *testing.T) {
	requireParallelCapacity(t)
	instance := &suiteProbe{onTestSetup: rendezvous()}
	suite.RunParallel(t, instance)
	if got := instance.testSetups.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.testSetups = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunParallelCompletesSetupBeforeEachTestBody(t *testing.T) {
	var prepared sync.Map
	instance := &suiteProbe{onTestSetup: func(t *testing.T) { prepared.Store(t.Name(), true) }}
	instance.onBody = func(t *testing.T) {
		if got := instance.suiteSetups.Load(); got != 1 {
			t.Errorf("suite setups before test body = %d; want 1", got)
		}
		if _, ready := prepared.Load(t.Name()); !ready {
			t.Error("test body ran before its SetupTest completed")
		}
	}
	suite.RunParallel(t, instance)
	if got := instance.bodies.sortedSnapshot(); !slices.Equal(got, []string{"TestAlpha", "TestBeta"}) {
		t.Errorf("instance.bodies = %v; want %v", got, []string{"TestAlpha", "TestBeta"})
	}
}

func TestRunParallelWaitsForAllTestBodies(t *testing.T) {
	var completed atomic.Int32
	suite.RunParallel(t, &suiteProbe{onBody: func(t *testing.T) {
		defer completed.Add(1)
	}})
	if got := completed.Load(); got != 2 {
		t.Errorf("completed bodies when runner returned = %d; want 2", got)
	}
}
