package mocks

import (
	"fmt"
	"testing"
)

type TestingT struct {
	*testing.T
	Failures []string
	Runs     int
}

func (s *TestingT) Fatal(args ...any) {
	s.Failures = append(s.Failures, fmt.Sprint(args...))
}

func (s *TestingT) Run(name string, test func(*testing.T)) bool {
	s.Runs++
	return s.T.Run(name, test)
}
