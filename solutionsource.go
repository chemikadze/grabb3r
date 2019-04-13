package grabb3r

import "time"

type SolutionSource interface {
	Login() error
	ListSolutions() (chan SolutionDesc, chan error)
	GetSolution(id SolutionDesc) (Solution, error)
}

type SolutionDesc interface {
	String() string
	ProblemName() string
	Language() Language
	SubmittedTime() time.Time
	Equals(other SolutionDesc) bool
}

type Solution interface {
	Desc() SolutionDesc
	Code() string
}

type Language string

func NewSource() SolutionSource {
	return NewMockSource()
}
