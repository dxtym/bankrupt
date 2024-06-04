package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/dxtym/bankrupt/api"
	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/gapi"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.Driver, config.Source)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(config, store)
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankruptServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("starting gRPC server on %s", config.GRPCAddress)
	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Fatal("cannot listen to grpc address:", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start server", err)
	}
}

func runHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	if err := server.Start(config.HTTPAddress); err != nil {
		log.Fatal("cannot start server", err)
	}
}
