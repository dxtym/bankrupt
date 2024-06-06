package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/dxtym/bankrupt/api"
	db "github.com/dxtym/bankrupt/db/sqlc"
	"github.com/dxtym/bankrupt/gapi"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGRPCServer(config, store)
	runGatewayServer(config, store)
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

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankruptHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register service:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	log.Printf("starting HTTP gateway server on %s", config.HTTPAddress)
	listener, err := net.Listen("tcp", config.HTTPAddress)
	if err != nil {
		log.Fatal("cannot listen to address:", err)
	}

	if err := http.Serve(listener, mux); err != nil {
		log.Fatal("cannot start HTTP gateway server", err)
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
