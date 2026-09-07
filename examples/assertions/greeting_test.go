package examples_test

import (
	"errors"
	"fmt"
	"testing"

	suite "github.com/nirajsapkota/suitespot"
	"github.com/nirajsapkota/suitespot/spot"
)

type GreetingSuite struct {
	greeting string
}

func (s *GreetingSuite) SetupSuite(t *testing.T) {
	s.greeting = "hello"
}

func (s *GreetingSuite) TestGreeting(t *testing.T) {
	spot.Equal(t, "hello", s.greeting)
	spot.True(t, len(s.greeting) > 0)
	spot.False(t, s.greeting == "goodbye")
	spot.NotEqual(t, "goodbye", s.greeting)
}

func (s *GreetingSuite) TestErrors(t *testing.T) {
	missing := errors.New("missing greeting")
	err := fmt.Errorf("lookup: %w", missing)
	spot.Err(t, err)
	spot.ErrorContains(t, err, "missing greeting")
	spot.NoErr(t, nil)
}

func (s *GreetingSuite) TestPointers(t *testing.T) {
	var missing *string
	spot.Nil(t, missing)
	spot.NotNil(t, &s.greeting)
}

func TestGreetingSuite(t *testing.T) {
	t.Parallel()
	suite.RunParallel(t, &GreetingSuite{})
}

func TestOrdering(t *testing.T) {
	spot.Greater(t, 3, 2)
	spot.GreaterOrEqual(t, 3, 3)
	spot.LessOrEqual(t, 2, 2)
	spot.Less(t, 2, 3)
}
