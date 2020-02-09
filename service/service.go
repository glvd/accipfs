package service

import "sync"

type Service struct {
	once *sync.Once
}

var _service *Service

func (s *Service) Register(node Node) {

}
