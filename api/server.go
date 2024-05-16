package api

import (
	"fmt"
	"testing"
	"time"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/token"
	"github.com/dxtym/bankrupt/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config utils.Config
	store  db.Store
	token  token.Maker
	router *gin.Engine
}

func newTestServer(t *testing.T, s db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey: utils.RandomString(32),
		TokenDuration:     time.Minute,
	}

	server, err := NewServer(config, s)
	if err != nil {
		t.Fatal("cannot create server:", err)
	}

	return server
}

func NewServer(config utils.Config, s db.Store) (*Server, error) {
	// chose paseto maker (can choose jwt too)
	token, err := token.NewPasetoMaker([]byte(config.TokenSymmetricKey))
	if err != nil {
		return nil, fmt.Errorf("cannot load token maker: %w", err)
	}
	server := &Server{
		config: config,
		store:  s,
		token:  token,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouting()
	return server, nil
}

func (s *Server) setUpRouting() {
	router := gin.Default()

	router.POST("/users", s.CreateUser)
	router.POST("/users/login", s.LoginUser)

	authRoute := router.Group("/").Use(authMiddleware(s.token))

	authRoute.POST("/accounts", s.createAccount)
	authRoute.GET("/accounts/:id", s.getAccount)
	authRoute.GET("/accounts", s.listAccount)

	authRoute.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
