package grabb3r

import "fmt"

type mockSolutionDestination struct{}

func (mockSolutionDestination) Initialize() error {
	return nil
}

func (mockSolutionDestination) SaveSolution(solution Solution) error {
	fmt.Printf("Solution %s submitted at %s\n", solution.Desc(), solution.Desc().SubmittedTime())
	return nil
}

func (mockSolutionDestination) ContainsSolution(solution SolutionDesc) (bool, error) {
	return false, nil
}

func NewMockDestination() SolutionDestination {
	return mockSolutionDestination{}
}
