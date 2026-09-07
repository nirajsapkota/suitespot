# Sequential and parallel execution

[SuiteSpot](../README.md)

Each linked example is a complete runnable package with two suites. Each suite
configures execution directly in its own `Test*Suite` entry point. The snippets
show four alternatives for the [workspace](fixtures.md) and
[greeting](suites.md) suite definitions.

Run the example commands from the root of a SuiteSpot checkout.

| Mode | Suite entry point calls `t.Parallel()` | Runner |
| --- | --- | --- |
| Suites sequential, tests sequential | No | `suite.Run` |
| Suites sequential, tests parallel | No | `suite.RunParallel` |
| Suites parallel, tests sequential | Yes | `suite.Run` |
| Suites parallel, tests parallel | Yes | `suite.RunParallel` |

## 1. Suites sequential, tests sequential

Use `suite.Run` directly in each entry point. Neither the entry points nor the
test methods are marked parallel:

```go
func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, &WorkspaceSuite{})
}

func TestGreetingSuite(t *testing.T) {
	suite.Run(t, &GreetingSuite{})
}
```

`WorkspaceSuite` and `GreetingSuite` run one at a time. Each suite's test methods
also run one at a time, in alphabetical order. Do not rely on the order of the
top-level suite entry points.

Run the [complete sequential-sequential example](../examples/sequential_sequential):

```sh
go test -v ./examples/sequential_sequential
```

## 2. Suites sequential, tests parallel

Use `suite.RunParallel` without calling `t.Parallel()` in the entry points:

```go
func TestWorkspaceSuite(t *testing.T) {
    suite.RunParallel(t, &WorkspaceSuite{})
}

func TestGreetingSuite(t *testing.T) {
    suite.RunParallel(t, &GreetingSuite{})
}
```

`WorkspaceSuite` and `GreetingSuite` run one at a time. Within each suite, the
test methods can overlap. All test methods finish before the next suite starts.
`SetupSuite` completes first; each test's `SetupTest` then runs before its own
test body. Different tests' `SetupTest` calls can overlap.

Run the [complete sequential-parallel example](../examples/sequential_parallel):

```sh
go test -v -race -parallel 8 ./examples/sequential_parallel
```

## 3. Suites parallel, tests sequential

Call `t.Parallel()` in each entry point and use `suite.Run`:

```go
func TestWorkspaceSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &WorkspaceSuite{})
}

func TestGreetingSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &GreetingSuite{})
}
```

`WorkspaceSuite` and `GreetingSuite` can overlap, including their `SetupSuite`
calls. Within each suite, `SetupTest` and its test method finish before the next
test starts. Test methods run in alphabetical order.

Run the [complete parallel-sequential example](../examples/parallel_sequential):

```sh
go test -v -parallel 8 ./examples/parallel_sequential
```

## 4. Suites parallel, tests parallel

Call `t.Parallel()` in each entry point and use `suite.RunParallel`:

```go
func TestWorkspaceSuite(t *testing.T) {
	t.Parallel()

	suite.RunParallel(t, &WorkspaceSuite{})
}

func TestGreetingSuite(t *testing.T) {
	t.Parallel()

	suite.RunParallel(t, &GreetingSuite{})
}
```

Both suites can overlap, and the test methods inside each suite can overlap.
Each suite finishes its `SetupSuite` before any of its `SetupTest` calls start.
Each test's own `SetupTest` completes before its test method starts. There is no
guaranteed execution order between parallel tests.

Run the [complete parallel-parallel example](../examples/parallel_parallel):

```sh
go test -v -race -parallel 8 ./examples/parallel_parallel
```

The `-parallel` flag limits active parallel tests within each test binary. Suites
in different packages can also run concurrently; see [running tests](running-tests.md)
for package-level limits and test selection.
