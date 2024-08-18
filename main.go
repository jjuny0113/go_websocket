package main

import (
	"flag"
	"fmt"
	"websocket_chatting/config"
	"websocket_chatting/network"
	"websocket_chatting/repository"
	"websocket_chatting/service"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "http server port")

func main() {
	flag.Parse()
	fmt.Println(*pathFlag, *port)
	c := config.NewConfig(*pathFlag)
	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {

		s := network.NewServer(service.NewService(rep), *port)
		s.StartServer()
	}

}
