package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/dxtym/bankrupt/api"
	db "github.com/dxtym/bankrupt/db/sqlc"
	_ "github.com/dxtym/bankrupt/doc/statik"
	"github.com/dxtym/bankrupt/gapi"
	"github.com/dxtym/bankrupt/pb"
	"github.com/dxtym/bankrupt/utils"
	"github.com/dxtym/bankrupt/worker"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config: %s", err)
	}

	if config.Environment == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.Driver, config.Source)
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s", err)
	}

	runDbMigration(config.MigrateURL, config.Source)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runProcessor(redisOpt, store)
	go runGRPCServer(config, store, taskDistributor)
	runGatewayServer(config, store, taskDistributor)
}

func runDbMigration(migrateURL, source string) {
	m, err := migrate.New(migrateURL, source)
	if err != nil {
		log.Fatal().Msgf("cannot create migration instance: %s", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("failed to migrate up: %s", err)
	}

	log.Info().Msg("db migrated succesfully")
}

func runProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	rtp := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("starting processor")
	if err := rtp.Run(); err != nil {
		log.Fatal().Msgf("cannot start processor: %s", err)
	}
}

func runGRPCServer(config utils.Config, store db.Store, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	logger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(logger)
	pb.RegisterBankruptServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Info().Msgf("starting gRPC server on %s", config.GRPCAddress)
	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		log.Fatal().Msgf("cannot listen to grpc address: %s", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msgf("cannot start server: %s", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store, td worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, td)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
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
		log.Fatal().Msgf("cannot register service: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// fs := http.FileServer(http.Dir("./doc/swagger"))
	// serve from server memory
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Msgf("cannot create file system: %s", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	log.Info().Msgf("starting HTTP gateway server on %s", config.HTTPAddress)
	listener, err := net.Listen("tcp", config.HTTPAddress)
	if err != nil {
		log.Fatal().Msgf("cannot listen to address: %s", err)
	}

	handler := gapi.HttpLogger(mux)
	if err := http.Serve(listener, handler); err != nil {
		log.Fatal().Msgf("cannot start HTTP gateway server: %s", err)
	}
}

func runHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	if err := server.Start(config.HTTPAddress); err != nil {
		log.Fatal().Msgf("cannot start server: %s", err)
	}
}
