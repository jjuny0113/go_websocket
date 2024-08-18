package service

import "websocket_chatting/repository"

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{
		repository: repository,
	}
	return s
}
