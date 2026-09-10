# Assertions with an explicit test

[SuiteSpot](../README.md)

Import `github.com/nirajsapkota/suitespot/pkg/spot`. Every assertion receives its test
as the first argument. There is no assertion object or stored suite test pointer,
and the package has no external dependencies.

```go
func TestGreeting(t *testing.T) {
    spot.Equal(t, "hello", greeting)
    spot.EqualValues(t, int64(5), len(greeting))
    spot.NotEqual(t, "goodbye", greeting)
    spot.True(t, ready)
    spot.False(t, stopped)
    spot.Nil(t, missing)
    spot.NotNil(t, result)
    spot.Error(t, lookupError)
    spot.NoError(t, saveError)
    spot.ErrorContains(t, lookupError, "not found")
    spot.NotEmpty(t, greeting)
    spot.Len(t, greeting, 5)
    spot.Regexp(t, `^h`, greeting)
}
```

Assertions call `t.Helper()` and stop the current test on failure through `t.Fatal`.
They return `true` on success. A failure returns `false` only when a test double's
`Fatal` returns.

All assertions accept an optional message or format string and arguments:

```go
spot.Equal(t, "ready", status, "resource %s", name)
spot.NoError(t, err, "could not load configuration")
```

## Helpers

| Helpers | Behavior |
| --- | --- |
| `True`, `False` | Check a boolean |
| `Nil`, `NotNil` | Check nil, including typed nil pointers, maps, slices, channels, and functions |
| `Error`, `NoError` | Check `err != nil` or `err == nil` |
| `Equal`, `NotEqual` | Deep equality, with matching types |
| `EqualValues` | Deep equality, then exact equality across numeric types |
| `Empty`, `NotEmpty` | Check nil, zero values, and collection lengths |
| `Contains`, `NotContains` | Check a substring, slice/array element, or map key |
| `Len` | Check the length of a string, array, slice, map, or channel |
| `ElementsMatch` | Compare slice/array elements without order, including duplicate counts |
| `ErrorContains` | Require a non-nil error whose message contains a substring |
| `Regexp` | Match a string with a pattern string or compiled regular expression |
| `Eventually` | Check a condition repeatedly until it passes or times out |
| `GreaterThan`, `GreaterThanOrEqual`, `LessThanOrEqual`, `LessThan` | Compare finite numbers |

`Equal` and `NotEqual` take the expected value before the actual value. Ordering
helpers read left to right:

```go
spot.GreaterThan(t, count, 0)
spot.GreaterThanOrEqual(t, count, minimum)
spot.LessThanOrEqual(t, count, maximum)
spot.LessThan(t, count, limit)
```

Ordering comparisons preserve integer precision and support mixed numeric types.
NaN, infinity, and unsupported ordering inputs report assertion failures.
`EqualValues` also compares mixed numeric types without losing precision. NaN is
never equal, while positive and negative infinities equal an infinity with the
same sign.

`Nil` and `NoError` answer different questions. An `error` interface containing a
typed nil pointer is non-nil: `Nil(t, err)` passes, but `NoError(t, err)` fails,
just like an ordinary `if err != nil` check.

## Collections

```go
spot.Contains(t, "hello world", "world")
spot.Contains(t, []string{"hello", "world"}, "world")
spot.Contains(t, map[string]int{"hello": 1}, "hello") // Checks keys, not values.
spot.NotContains(t, []string{"hello", "world"}, "goodbye")
spot.Len(t, "hello", 5)
spot.ElementsMatch(t, []int{1, 2, 2}, []int{2, 1, 2})
```

`Contains` uses substring matching for strings and deep equality for elements and
map keys. String targets must also be strings. `ElementsMatch` compares each
element using deep equality, ignores order, and requires duplicate counts to
match. Neither helper converts numeric types or modifies its inputs.

`ElementsMatch` accepts arrays and slices (including comparisons between them),
and treats nil as an empty list. Nil and empty slices therefore match. Other
inputs, including strings and maps, fail even when empty. Unsupported inputs
report assertion failures rather than panicking.

`Empty` accepts nil and typed nil values, zero-length strings and collections,
numeric zero, false, zero-value structs, and pointers to empty values. Arrays
are empty only when their length is zero.

## Matching and polling

```go
spot.Regexp(t, `^ready-\d+$`, status)
spot.Regexp(t, regexp.MustCompile(`^ready-\d+$`), status)

spot.Eventually(t, func() bool {
    return service.Ready()
}, time.Second, 10*time.Millisecond)
```

`Regexp` accepts ordinary and named strings for both patterns and values.
Invalid patterns and unsupported values fail the assertion. `Eventually` checks
immediately, then on each tick until the wait time expires. Its wait and tick
must both be positive, and its condition must return without blocking.

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
