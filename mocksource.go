package grabb3r

import (
	"errors"
	"fmt"
	"time"
)

type mockSource struct {
	solutions []simpleSolution
}

func (s *mockSource) Login() error {
	return nil
}

func (s *mockSource) ListSolutions() (chan SolutionId, chan error) {
	resChan := make(chan SolutionId, 10)
	errChan := make(chan error, 1)
	go func() {
		for _, solution := range s.solutions {
			resChan <- solution.Id()
		}
		close(resChan)
	}()
	return resChan, errChan
}

func (s *mockSource) GetSolution(id SolutionId) (Solution, error) {
	for _, solution := range s.solutions {
		if solution.Id() == id { // reference quality!
			return &solution, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Can't find solution '%s'", id))
}

type simpleSolutionId struct {
	id string
}

func (s *simpleSolutionId) Equals(other SolutionId) bool {
	switch v := other.(type) {
	case *simpleSolutionId:
		return s.id == v.id
	default:
		return false
	}
}

func (s *simpleSolutionId) String() string {
	return s.id
}

type simpleSolution struct {
	id            SolutionId
	code          string
	language      Language
	submittedTime time.Time
}

func (s *simpleSolution) Id() SolutionId {
	return s.id
}

func (s *simpleSolution) Code() string {
	return s.code
}

func (s *simpleSolution) SubmittedTime() time.Time {
	return s.submittedTime
}

func (s *simpleSolution) Language() Language {
	return s.language
}

func NewMockSource() SolutionSource {
	return &mockSource{[]simpleSolution{
		{id: &simpleSolutionId{"solution1"}, code: "print 'Hello, World!'", language: "Python", submittedTime: time.Now()},
	}}
}
