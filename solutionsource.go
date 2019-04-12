package grabb3r

import "time"

type SolutionSource interface {
	Login() error
	ListSolutions() (chan SolutionId, chan error)
	GetSolution(id SolutionId) (Solution, error)
}

type SolutionId interface {
	String() string
	Equals(other SolutionId) bool
}

type Solution interface {
	Id() SolutionId
	Code() string
	Language() Language
	SubmittedTime() time.Time
}

type Language string

func NewSource() SolutionSource {
	return NewMockSource()
}
