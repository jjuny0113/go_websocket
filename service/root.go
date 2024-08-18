package service

import (
	"errors"
	"fmt"
	"log"
	"websocket_chatting/repository"
	"websocket_chatting/types/schema"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{
		repository: repository,
	}
	return s
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	res, err := s.repository.GetChatList(roomName)
	if err != nil {
		log.Println("Failed to get chat list", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	res, err := s.repository.RoomList()
	if err != nil {
		log.Println("Failed to get room list", err.Error())
		return nil, err
	}
	return res, nil
}

func (s *Service) MakeRoom(name string) error {
	err := s.repository.MakeRoom(name)
	if err != nil {
		log.Println("Failed to make new room", err.Error())
		return err
	}
	return err
}

func (s *Service) Room(name string) (*schema.Room, error) {
	res, err := s.repository.Room(name)
	if err != nil {
		log.Println("Failed to get room", err.Error())
		return nil, err
	}

	return res, nil
}

func (s *Service) InsertChatting(user, message, roomName string) error {
	if s == nil || s.repository == nil {
		return errors.New("service or repository is nil")
	}

	err := s.repository.InsertChatting(user, message, roomName)
	if err != nil {
		log.Printf("Failed to insert chatting: user=%s, room=%s, error=%v", user, roomName, err)
		return fmt.Errorf("failed to insert chatting: %w", err)
	}

	log.Printf("Successfully inserted chatting: user=%s, room=%s", user, roomName)
	return nil
}
