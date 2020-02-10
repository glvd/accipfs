package service

import "sync"

// Service ...
type Service struct {
	once  *sync.Once
	nodes []Node
}

// Register ...
func (s *Service) Register(node Node) {
	s.nodes = append(s.nodes, node)
}

// Run ...
func (s *Service) Run() {
	s.once.Do(
		func() {
			for _, node := range s.nodes {
				node.Start()
			}
		})
}
