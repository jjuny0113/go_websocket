package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"websocket_chatting/repository"
	"websocket_chatting/service"
)

type Server struct {
	engine     *gin.Engine
	service    *service.Service
	repository *repository.Repository
	port       string
	ip         string
}

func NewServer(service *service.Service, repository *repository.Repository, port string) *Server {
	s := &Server{
		engine:     gin.New(),
		service:    service,
		repository: repository,
		port:       port,
	}

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery()) // 에러로 인한 서버 다운시 처리
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))
	registerServer(s.engine)
	return s
}

func (s *Server) StartServer() error {
	log.Println("Starting server...")
	return s.engine.Run(s.port)
}
