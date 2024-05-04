package api

import (
	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(s db.Store) *Server {
	server := &Server{store: s}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
} 

func (s *Server) Start(address string) error {
	return s.router.Run(address)
} 

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}