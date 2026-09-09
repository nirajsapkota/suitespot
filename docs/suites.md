# Creating suites

[SuiteSpot](../README.md)

SuiteSpot requires Go 1.22 or later and has no external dependencies.
Add it to your Go module:

```sh
go get github.com/nirajsapkota/suitespot
```

Put this complete suite in `greeting_test.go`. Its entry point selects sequential
execution. `SetupSuite` creates the shared fixture once, and there is no per-test
preparation, so this suite does not need `SetupTest`:

```go
package examples_test

import (
	"testing"

	"github.com/nirajsapkota/suitespot/pkg/suite"
)

func TestGreetingSuite(t *testing.T) {
	suite.Run(t, &GreetingSuite{})
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
```

Go discovers `TestGreetingSuite` automatically. The library then discovers
`TestLength` and `TestValue` on that instance. Exported method names starting with
`Test` must have signature `func(*testing.T)`; other methods are ignored except
for the setup hooks.

Run your suite with `go test ./...`. Add more suites by placing a `Test*Suite`
entry point beside each suite struct; there is no central list to maintain.

Continue with [per-suite and per-test fixtures](fixtures.md),
[parallel execution](parallelism.md), or [assertions](assertions.md).
