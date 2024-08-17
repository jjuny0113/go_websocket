package network

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"websocket_chatting/types"
)

// http -> websocket으로 upgrade
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.MessageBufferSize,
	WriteBufferSize: types.SocketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Room struct {
	Forward chan *message
	/**
	수신되는 메시지를 보관하는 값
	들어오는 메시지를 다른 클라이언트들에게 전송됨
	*/
	Join    chan *client     // Socket 연결되는 경우에 작동
	Leave   chan *client     // Socket 끊어지는 경우에 작동
	Clients map[*client]bool // 현재 방에 있는 client 정보를 저장
}

type client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		Clients: make(map[*client]bool),
	}
}

func (c *client) Read() {
	// 클라이언트의 들어오는 메시지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			}
			panic(err)
			return
		}
		log.Println("HEAD : ", msg, "client", c.Name)
		msg.Time = time.Now().Unix()
		msg.Name = c.Name
		c.Room.Forward <- &msg
	}
}
func (c *client) Write() {
	// 클라이언트의 메시지를 전송하는 함수
	defer c.Socket.Close()
	for msg := range c.Send {
		log.Println("HEAD : ", msg, "client", c.Name)
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) RunInit() {
	// room에 있는 모든 채널값들을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

func (r *Room) SocketServe(c *gin.Context) {

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}
	client := &client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client

	// connection을 끊어줌
	defer func() {
		r.Leave <- client
	}()

	go client.Read()
	client.Write()
}
