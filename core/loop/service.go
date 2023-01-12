package loop

import "chasqi-go/types"

type Service struct {
	engine Engine
}

func NewService(statusGetter Engine) *Service {
	return &Service{statusGetter}
}

func (s *Service) GetStatus(id string) (*types.LoopStatus, error) {
	return s.engine.ById(id)
}

func (s *Service) Enqueue(tree *types.Tree) error {
	return nil
}
