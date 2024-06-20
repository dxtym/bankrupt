package gapi

import (
	"fmt"

	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/token"
	"github.com/dxtym/bankrupt/utils"
)

type Server struct {
	pb.UnimplementedBankruptServer
	config utils.Config
	store  db.Store
	token  token.Maker
}

func NewServer(config utils.Config, s db.Store) (*Server, error) {
	token, err := token.NewPasetoMaker([]byte(config.TokenSymmetricKey))
	if err != nil {
		return nil, fmt.Errorf("cannot load token maker: %w", err)
	}
	server := &Server{
		config: config,
		store:  s,
		token:  token,
	}
	return server, nil
}
