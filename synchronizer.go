package grabb3r

type SolutionSynchronizer interface {
	Synchronize() error
}

type simpleSolutionSynchronizer struct {
	src SolutionSource
	dst SolutionDestination
}

func (s *simpleSolutionSynchronizer) Synchronize() error {
	solutionChan, errChan := s.src.ListSolutions()
	for {
		select {
		case solutionId, ok := <-solutionChan:
			if !ok {
				solutionChan = nil
				goto channelClosed
			}
			if contains, err := s.dst.ContainsSolution(solutionId); err != nil {
				return err
			} else if contains {
				continue
			}
			solution, err := s.src.GetSolution(solutionId)
			if err != nil {
				return err
			}
			if err := s.dst.SaveSolution(solution); err != nil {
				return err
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				goto channelClosed
			} else {
				return err
			}
		}
	channelClosed:
		if solutionChan == nil && errChan == nil {
			break
		}
	}
	return nil
}

func NewSyncronizer(src SolutionSource, dst SolutionDestination) SolutionSynchronizer {
	return &simpleSolutionSynchronizer{src: src, dst: dst}
}
