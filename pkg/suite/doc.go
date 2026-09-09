// Package suite runs struct-based test suites using Go's standard testing package.
// Suites provide SetupSuite and Test methods accepting *testing.T, with an
// optional SetupTest hook. Run executes test methods sequentially; RunParallel
// executes them in parallel against the same suite instance.
package suite
