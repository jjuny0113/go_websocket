package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"websocket_chatting/config"
	"websocket_chatting/types/schema"
)

type Repository struct {
	cfg *config.Config
	db  *sql.DB
}

const (
	room   = "chatting.room"
	chat   = "chatting.chat"
	server = "chatting.server"
)

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}
	var err error
	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	qs := query([]string{"SELECT * FROM", chat, "WHERE room = ? ORDER BY `when` DESC LIMIT 10"})
	cursor, err := s.db.Query(qs, roomName)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var results []*schema.Chat

	for cursor.Next() {
		d := new(schema.Chat)
		if err = cursor.Scan(&d.Name, &d.Id, &d.Room, &d.When, &d.Message); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	if len(results) == 0 {
		return []*schema.Chat{}, nil
	}
	return results, nil
}

func (s *Repository) RoomList() ([]*schema.Room, error) {
	// TODO 페이징 추가하기
	qs := query([]string{"SELECT * FROM room"})
	cursor, err := s.db.Query(qs)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var results []*schema.Room

	for cursor.Next() {
		d := new(schema.Room)
		if err = cursor.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	if len(results) == 0 {
		return []*schema.Room{}, nil
	}
	return results, nil
}

func (s *Repository) MakeRoom(name string) error {
	_, err := s.db.Exec("INSERT INTO chatting.room(name) VALUES (?)", name)
	return err
}

func (s *Repository) Room(name string) (*schema.Room, error) {
	d := new(schema.Room)
	qs := query([]string{"Select * From", room, "Where name = ?"})
	/**
	select * from chatting.room where name = ?
	*/

	err := s.db.QueryRow(qs, name).Scan(
		&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt,
	)

	if err == nil {
		return nil, nil
	}

	return d, err
}

func query(qs []string) string {
	return strings.Join(qs, " ") + ";"
}
