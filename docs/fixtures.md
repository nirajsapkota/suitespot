# Per-suite and per-test fixtures

[SuiteSpot](../README.md)

Put this complete suite in `workspace_test.go`. It uses the sequential entry
point; change only `TestWorkspaceSuite` to select one of the other
[execution modes](parallelism.md).

```go
package examples_test

import (
	"os"
	"path/filepath"
	"testing"

	suite "github.com/nirajsapkota/suitespot"
)

func TestWorkspaceSuite(t *testing.T) {
	suite.Run(t, &WorkspaceSuite{})
}

type WorkspaceSuite struct {
	root string
}

func (s *WorkspaceSuite) SetupSuite(t *testing.T) {
	s.root = t.TempDir()
}

func (s *WorkspaceSuite) SetupTest(t *testing.T) {
	if err := os.MkdirAll(s.directory(t), 0755); err != nil {
		t.Fatal(err)
	}
}

func (s *WorkspaceSuite) TestStartsEmpty(t *testing.T) {
	entries, err := os.ReadDir(s.directory(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("new test directory has %d entries; want 0", len(entries))
	}
}

func (s *WorkspaceSuite) TestWrite(t *testing.T) {
	filename := filepath.Join(s.directory(t), "message.txt")
	if err := os.WriteFile(filename, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "hello" {
		t.Fatalf("contents = %q; want hello", contents)
	}
}

func (s *WorkspaceSuite) directory(t *testing.T) string {
	return filepath.Join(s.root, t.Name())
}
```

`SetupSuite` creates one root directory for this suite execution. The suite
instance stores that root and shares it with all test methods.

`SetupTest` creates a separate directory for each test that runs. It receives the
same `*testing.T` as that test method, so both compute the same directory from
`t.Name()`. `TestWrite` can create a file while `TestStartsEmpty` still sees an
empty directory of its own.

For sequential tests, the calls are:

```text
SetupSuite
SetupTest → TestStartsEmpty
SetupTest → TestWrite
```

For parallel tests, `SetupSuite` still runs first. The two
`SetupTest → test method` pairs can then overlap.

| Method | How often | What to put there |
| --- | --- | --- |
| `SetupSuite(t *testing.T)` | Once per execution of this suite | Shared fixture initialization |
| `SetupTest(t *testing.T)` | Once before each selected suite test method | Preparation specific to that test |

`SetupSuite` is required; it can have an empty body. `SetupTest` is optional.
Both return nothing. A skip or fatal failure in `SetupSuite` prevents the suite's
tests from running; in `SetupTest`, it prevents that test's body from running.
Nested subtests that you create yourself do not get additional `SetupTest` calls.

All methods use the same suite instance. With parallel tests, `SetupTest` must
not reset shared fields such as `s.currentRequest` for each test: another test
could be using them. Use test-local state, distinct resources as above, or
synchronize shared mutable state.

See the [complete workspace example](../examples/sequential_sequential/workspace_test.go)
and [how to select individual tests](running-tests.md).
