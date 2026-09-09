# Assertions with an explicit test

[SuiteSpot](../README.md)

Import `github.com/nirajsapkota/suitespot/pkg/spot`. Every assertion receives its test
as the first argument. There is no assertion object or stored suite test pointer,
and the package has no external dependencies.

```go
func TestGreeting(t *testing.T) {
    spot.Equal(t, "hello", greeting)
    spot.NotEqual(t, "goodbye", greeting)
    spot.True(t, ready)
    spot.False(t, stopped)
    spot.Nil(t, missing)
    spot.NotNil(t, result)
    spot.Err(t, lookupError)
    spot.NoErr(t, saveError)
    spot.ErrorContains(t, lookupError, "not found")
}
```

Assertions call `t.Helper()` and stop the current test on failure through `t.Fatal`.
They return `true` on success. A failure returns `false` only when a test double's
`Fatal` returns.

All assertions accept an optional message or format string and arguments:

```go
spot.Equal(t, "ready", status, "resource %s", name)
spot.NoErr(t, err, "could not load configuration")
```

## Helpers

| Helpers | Behavior |
| --- | --- |
| `True`, `False` | Check a boolean |
| `Nil`, `NotNil` | Check nil, including typed nil pointers, maps, slices, channels, and functions |
| `Err`, `NoErr` | Check `err != nil` or `err == nil` |
| `Equal`, `NotEqual` | Deep equality, with matching types |
| `ErrorContains` | Require a non-nil error whose message contains a substring |
| `Greater`, `GreaterOrEqual`, `LessOrEqual`, `Less` | Compare finite numbers |

`Equal` and `NotEqual` take the expected value before the actual value. Ordering
helpers read left to right:

```go
spot.Greater(t, count, 0)
spot.GreaterOrEqual(t, count, minimum)
spot.LessOrEqual(t, count, maximum)
spot.Less(t, count, limit)
```

Numeric comparisons preserve integer precision and support mixed numeric types.
NaN, infinity, and unsupported comparison inputs report assertion failures.

`Nil` and `NoErr` answer different questions. An `error` interface containing a
typed nil pointer is non-nil: `Nil(t, err)` passes, but `NoErr(t, err)` fails,
just like an ordinary `if err != nil` check.

## Parallel suites

```go
func (s *GreetingSuite) TestValue(t *testing.T) {
    spot.Equal(t, "hello", s.greeting)
}

func TestGreetingSuite(t *testing.T) {
    t.Parallel()
    suite.RunParallel(t, &GreetingSuite{})
}
```

Use each method's `t` for its assertions. Shared fixtures still need immutable
state or synchronization; explicit test arguments do not make shared data safe
to mutate concurrently.

See the [runnable example](../examples/assertions/greeting_test.go).

See [creating suites](suites.md) and [parallel execution](parallelism.md).
