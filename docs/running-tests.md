# Running and selecting tests

[SuiteSpot](../README.md)

Run these commands from the root of a SuiteSpot checkout:

```sh
go test ./...
go test -race -parallel 8 ./...
go test -run '^TestWorkspaceSuite$' ./examples/parallel_parallel
go test -run '^TestWorkspaceSuite$/WorkspaceSuite/TestWrite$' ./examples/parallel_parallel
```

The `-parallel` flag limits active parallel tests within each test binary. Suites
in separate packages can still run concurrently even in the sequential example;
use `go test -p 1 ./...` if packages must also run sequentially.

Calling `t.Parallel()` in `Test*Suite` marks that suite entry point as parallel.
Test methods should leave their own `t.Parallel()` calls to the runner.
`suite.Run` and `suite.RunParallel` both accept multiple suites, which run
sequentially within that call. `RunParallel` makes the test methods within each
suite parallel; the entry point's `t.Parallel()` controls suite parallelism.

## Validation

Both runners require at least one suite and fail the calling test if none are
provided. Every supplied suite is validated before its subtest starts, even if
that suite would be excluded by `-run`.

## Testing the runner

`Run` and `RunParallel` accept the small `TestingT` interface. A normal
`*testing.T` satisfies it, so callers still write `suite.Run(t, &GreetingSuite{})`.
Validation failures call `Fatal` automatically; callers have no error to check.

The library's [tests](../suite_test.go) pass a stub that records `Fatal` calls to
assert validation failures directly. The runner stops after reporting a failure,
even when the stub's `Fatal` returns. Setup hooks and test methods still receive
real `*testing.T` values, so their assertions and parallel scheduling use Go's
standard test runner.

The suite runner implementation remains in [suite.go](../suite.go).

See [creating suites](suites.md) and the [four execution modes](parallelism.md).
