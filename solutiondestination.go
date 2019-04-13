package grabb3r

type SolutionDestination interface {
	Initialize() error
	SaveSolution(solution Solution) error
	ContainsSolution(solution SolutionDesc) (bool, error)
}
