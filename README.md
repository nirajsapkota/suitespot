<h1 align="center">SuiteSpot</h1>

<p align="center">
  Dependency-free test suites and assertions for Go.
</p>

<p align="center">
  <strong>Prepare the fixture once. Let each test focus on behavior.</strong>
</p>

<p align="center">
  <a href="#quick-start">Quick start</a> &nbsp;·&nbsp;
  <a href="#documentation">Documentation</a> &nbsp;·&nbsp;
  <a href="examples/">Examples</a>
</p>

---

SuiteSpot gives every test a fixture that's ready to use. Create your shared
fixture once in `SetupSuite`, and it's immediately available to every test method
on the suite. Each test can focus on the behavior it checks, with the scenario
setup already taken care of.

- **One place for setup:** keep shared scenario preparation in `SetupSuite`.
- **Focused test cases:** use the prepared fixture directly, keeping setup out of
  individual test methods. When per-test preparation is needed, use `SetupTest`.
- **Standard Go testing:** run with `go test`, choose sequential or parallel
  execution, and use the optional `spot` assertions with each test's `*testing.T`.

## Installation

Requires **Go 1.22 or later**.

```sh
go get github.com/nirajsapkota/suitespot
```

## Quick start

Create `greeting_test.go` in your project:

```go
package greeting_test

import (
    "testing"

    "github.com/nirajsapkota/suitespot/pkg/spot"
    "github.com/nirajsapkota/suitespot/pkg/suite"
)

type GreetingSuite struct {
    greeting string
}

func (s *GreetingSuite) SetupSuite(t *testing.T) {
    s.greeting = "hello"
}

func (s *GreetingSuite) TestValue(t *testing.T) {
    spot.Equal(t, "hello", s.greeting)
}

func (s *GreetingSuite) TestLength(t *testing.T) {
    spot.Equal(t, 5, len(s.greeting))
}

func TestGreetingSuite(t *testing.T) {
    suite.Run(t, &GreetingSuite{})
}
```

Run the tests:

```sh
go test -v ./...
```

## Documentation

- [Creating suites](docs/suites.md): suite definitions, setup, and test discovery.
- [Fixtures](docs/fixtures.md): shared resources and per-test preparation.
- [Parallel execution](docs/parallelism.md): all four execution modes.
- [Assertions](docs/assertions.md): helpers, messages, and comparison behavior.
- [Running tests](docs/running-tests.md): test selection, execution limits, and validation.
- [Runnable examples](examples/): complete suites for each execution mode and assertions.
