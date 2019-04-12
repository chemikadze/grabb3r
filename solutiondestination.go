package grabb3r

type SolutionDestination interface {
	SaveSolution(solution Solution) error
	ContainsSolution(solution SolutionId) (bool, error)
}
