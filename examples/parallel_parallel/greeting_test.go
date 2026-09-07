package examples_test

import (
	"testing"

	suite "github.com/nirajsapkota/suitespot"
)

func TestGreetingSuite(t *testing.T) {
	t.Parallel()

	suite.RunParallel(t, &GreetingSuite{})
}

type GreetingSuite struct {
	greeting string
}

func (s *GreetingSuite) SetupSuite(t *testing.T) {
	s.greeting = "hello"
}

func (s *GreetingSuite) TestValue(t *testing.T) {
	if s.greeting != "hello" {
		t.Fatalf("greeting = %q; want hello", s.greeting)
	}
}

func (s *GreetingSuite) TestLength(t *testing.T) {
	if len(s.greeting) != 5 {
		t.Fatalf("greeting length = %d; want 5", len(s.greeting))
	}
}
