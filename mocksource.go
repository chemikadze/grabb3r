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

func (s *mockSource) ListSolutions() (chan SolutionDesc, chan error) {
	resChan := make(chan SolutionDesc, 10)
	errChan := make(chan error, 1)
	go func() {
		for _, solution := range s.solutions {
			resChan <- solution.Desc()
		}
		close(resChan)
	}()
	return resChan, errChan
}

func (s *mockSource) GetSolution(id SolutionDesc) (Solution, error) {
	for _, solution := range s.solutions {
		if solution.Desc() == id { // reference quality!
			return &solution, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Can't find solution '%s'", id))
}

// bare minimum desc implementation
type simpleSolutionDesc struct {
	id            string
	language      Language
	submittedTime time.Time
}

func (s *simpleSolutionDesc) Equals(other SolutionDesc) bool {
	switch v := other.(type) {
	case *simpleSolutionDesc:
		return s.id == v.id
	default:
		return false
	}
}

func (s *simpleSolutionDesc) String() string {
	return s.id
}

func (s *simpleSolutionDesc) ProblemName() string {
	return "The Main Question Of Universe And Everything"
}

func (s *simpleSolutionDesc) SubmittedTime() time.Time {
	return s.submittedTime
}

func (s *simpleSolutionDesc) Language() Language {
	return s.language
}

// bare minimum solution implementation
type simpleSolution struct {
	desc SolutionDesc
	code string
}

func (s *simpleSolution) Desc() SolutionDesc {
	return s.desc
}

func (s *simpleSolution) Code() string {
	return s.code
}

func NewMockSource() SolutionSource {
	return &mockSource{[]simpleSolution{
		{desc: &simpleSolutionDesc{id: "solution1", language: "Python", submittedTime: time.Now()}, code: "print '42'"},
	}}
}
