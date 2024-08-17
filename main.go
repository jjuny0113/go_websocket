package main

import "websocket_chatting/network"

func main() {
	network := network.NewServer()
	network.StartServer()
}
