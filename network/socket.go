package network

import "github.com/gorilla/websocket"

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
	Join    chan *Client     // Socket 연결되는 경우에 작동
	Leave   chan *Client     // Socket 끊어지는 경우에 작동
	Clients map[*Client]bool // 현재 방에 있는 Client 정보를 저장
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}
