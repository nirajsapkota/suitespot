package examples_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nirajsapkota/suitespot/pkg/suite"
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
