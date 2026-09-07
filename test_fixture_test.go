package suite_test

import (
	"flag"
	"path"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type suiteProbe struct {
	suiteSetups atomic.Int32
	testSetups  callLog
	bodies      callLog
	onTestSetup func(*testing.T)
	onBody      func(*testing.T)
}

func (s *suiteProbe) SetupSuite(t *testing.T) {
	s.suiteSetups.Add(1)
}

func (s *suiteProbe) SetupTest(t *testing.T) {
	s.testSetups.add(path.Base(t.Name()))
	if s.onTestSetup != nil {
		s.onTestSetup(t)
	}
}

func (s *suiteProbe) TestAlpha(t *testing.T) {
	s.recordBody(t)
}

func (s *suiteProbe) TestBeta(t *testing.T) {
	s.recordBody(t)
}

func (s *suiteProbe) Helper(t *testing.T) {
	s.recordBody(t)
}

func (s *suiteProbe) BeforeTest(t *testing.T) {
	s.recordBody(t)
}

func (s *suiteProbe) testHidden(t *testing.T) {
	s.recordBody(t)
}

func (s *suiteProbe) recordBody(t *testing.T) {
	s.bodies.add(path.Base(t.Name()))
	if s.onBody != nil {
		s.onBody(t)
	}
}

type suiteWithoutSetupTest struct {
	bodies atomic.Int32
}

func (*suiteWithoutSetupTest) SetupSuite(*testing.T) {
}

func (s *suiteWithoutSetupTest) TestAlpha(*testing.T) {
	s.bodies.Add(1)
}

func (s *suiteWithoutSetupTest) TestBeta(*testing.T) {
	s.bodies.Add(1)
}

type valueSuite struct {
}

func (valueSuite) SetupSuite(*testing.T) {
}

func (valueSuite) TestValue(*testing.T) {
}

type scalarSuite int

func (*scalarSuite) SetupSuite(*testing.T) {
}

func (*scalarSuite) TestValue(*testing.T) {
}

type suiteWithoutTests struct {
}

func (*suiteWithoutTests) SetupSuite(*testing.T) {
}

type setupWithoutArguments struct {
	valueSuite
}

func (*setupWithoutArguments) SetupTest() {
}

type setupWithWrongArgument struct {
	valueSuite
}

func (*setupWithWrongArgument) SetupTest(int) {
}

type setupWithExtraArgument struct {
	valueSuite
}

func (*setupWithExtraArgument) SetupTest(*testing.T, int) {
}

type setupWithReturnValue struct {
	valueSuite
}

func (*setupWithReturnValue) SetupTest(*testing.T) int {
	return 0
}

type testWithWrongSignature struct {
	valueSuite
}

func (*testWithWrongSignature) TestValue(*testing.T) int {
	return 0
}

type callLog struct {
	mu    sync.Mutex
	names []string
}

func (c *callLog) add(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.names = append(c.names, name)
}

func (c *callLog) snapshot() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return slices.Clone(c.names)
}

func (c *callLog) sortedSnapshot() []string {
	names := c.snapshot()
	slices.Sort(names)
	return names
}

func requireParallelCapacity(t *testing.T) {
	t.Helper()
	limit, err := strconv.Atoi(flag.Lookup("test.parallel").Value.String())
	if err != nil {
		t.Fatal(err)
	}
	if limit < 2 {
		t.Skip("overlap requires -parallel >= 2")
	}
}

func rendezvous() func(*testing.T) {
	var arrived atomic.Int32
	ready := make(chan struct{})
	return func(t *testing.T) {
		t.Helper()
		if arrived.Add(1) == 2 {
			close(ready)
		}
		waitForPeer(t, ready)
	}
}

func waitForPeer(t *testing.T, ready <-chan struct{}) {
	t.Helper()
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()
	select {
	case <-ready:
	case <-timer.C:
		t.Fatal("parallel work did not overlap")
	}
}
