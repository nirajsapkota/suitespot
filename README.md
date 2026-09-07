# SuiteSpot

Simple Go test suites that hit the sweet spot.

SuiteSpot is a small, dependency-free Go testing library for grouping related tests
around a shared fixture. Suites are ordinary structs with independent test methods,
once-per-suite setup, and optional setup before each test. Each suite has its own
Go test entry point, so there is no central suite list to maintain.

Run suites and their tests sequentially or in parallel, with each choice configured
once. The companion `spot` package provides 13 assertions that receive `t` explicitly
and stop the current test on failure. Assertions never rely on a mutable test pointer
stored on the suite.

SuiteSpot requires Go 1.22 or later. It uses Go's standard test runner, and the suite
runner itself fits in [one file](suite.go).

## Documentation

- [Creating suites](docs/suites.md) — installation, a complete first suite, and discovery.
- [Fixtures](docs/fixtures.md) — `SetupSuite`, `SetupTest`, and shared state.
- [Sequential and parallel execution](docs/parallelism.md) — all four execution modes.
- [Assertions](docs/assertions.md) — the `spot` helpers and examples.
- [Running tests](docs/running-tests.md) — selecting suites and tests, execution limits, and validation.

Runnable examples live in [examples](examples/).
