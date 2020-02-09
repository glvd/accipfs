package service

import "sync"

type Service struct {
	once  *sync.Once
	nodes []Node
}

func (s *Service) Register(node Node) {
	s.nodes = append(s.nodes, node)
}
