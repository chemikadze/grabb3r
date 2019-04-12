package grabb3r

import "fmt"

type mockSolutionDestination struct{}

func (mockSolutionDestination) SaveSolution(solution Solution) error {
	fmt.Printf("Solution %s submitted at %s\n", solution.Id(), solution.SubmittedTime())
	return nil
}

func (mockSolutionDestination) ContainsSolution(solution SolutionId) (bool, error) {
	return false, nil
}

func NewMockDestination() SolutionDestination {
	return mockSolutionDestination{}
}
