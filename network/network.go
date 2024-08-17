package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Network struct {
	engine *gin.Engine
}

func NewServer() *Network {
	n := &Network{
		engine: gin.New(),
	}

	n.engine.Use(gin.Logger())
	n.engine.Use(gin.Recovery()) // 에러로 인한 서버 다운시 처리
	n.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowCredentials: true,
	}))

	return n
}

func (n *Network) StartServer() error {
	log.Println("Starting server...")
	return n.engine.Run(":8080")
}
