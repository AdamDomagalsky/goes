package main

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/AdamDomagalsky/goes/bank/api"
	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/gapi"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

//go:embed docs/swagger/dist/*
var swaggerFiles embed.FS

func main() {
	config, err := util.LoadConfig(".")
	fmt.Printf("config: %v\n", config)
	if err != nil {
		log.Fatal("Cannot load env config:", err)
	}
	conn, err := sql.Open(config.DATABASE_DRVIER, util.DbURL(config))
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	if config.GIN_MODE != "release" { // TODO generalize flag to work Gin <> Grpc wise
		err = db.MigrateUp(conn, config.DATABASE_NAME)
		if err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(
					fmt.Sprintf("GIN-%s, MigratingUp(%s) - failed:", config.GIN_MODE, config.DATABASE_NAME),
					err)
			} else {
				log.Printf("GIN-%s, MigratingUp(%s) - no change\n", config.GIN_MODE, config.DATABASE_NAME)
			}
		} else {
			log.Printf("GIN-%s, MigratingUp(%s) - succeed \n", config.GIN_MODE, config.DATABASE_NAME)
		}
	}

	store := db.NewStore(conn)
	go runGinAPIServer(config, store) // run GIN API server in separate goroutine
	go runGrpcAPIServer(config, store)
	runGrpcGatewayAPIServer(config, store)
}

func runGinAPIServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}

func runGrpcAPIServer(config util.Config, store db.Store) {
	listener, err := net.Listen("tcp", config.GRPC_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("could not create grpc server: %v", err)
	}

	serverOption := gapi.Setup0logGrpcLogger(config.ENVIROMENT)
	// serverOption, err := gapi.SetupZapGrpcLogger(config.ENVIROMENT)
	if err != nil {
		log.Fatalf("could not setup the logger: %v", err)
	}
	grpcServer := grpc.NewServer(serverOption...)

	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer) // TODO kind of self documentation

	log.Printf("start gRPC server at %s\n", listener.Addr().String())

	if err != grpcServer.Serve(listener) {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func runGrpcGatewayAPIServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("could not create grpc server: %v", err)
	}

	grpcMux := runtime.NewServeMux(
		// snake_case API
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// pb.RegisterBankHandlerFromEndpoint(...) // TODO split GW to separate service, see: https://blog.logrocket.com/guide-to-grpc-gateway/)
	err = pb.RegisterBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerPath := "docs/swagger/dist"
	// any change will be reflected immediately as it server files from disk (not embedded)
	// http://localhost:2138/swagger/
	indexFsDynamicDir := http.FileServer(http.Dir(swaggerPath))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", indexFsDynamicDir))

	// any change will be reflected after recompilation as it embeds files into binary
	// http://localhost:2138/swagger-embedded/
	indexFsEmbedded, err := fs.Sub(swaggerFiles, swaggerPath)
	if err != nil {
		log.Fatalf("cannot get subdirectory %v\n", err)
	}
	mux.Handle(
		"/swagger-embedded/",
		http.StripPrefix(
			"/swagger-embedded/",
			http.FileServer(http.FS(indexFsEmbedded)),
		),
	)

	listener, err := net.Listen("tcp", config.GRPC_API_GATEWAY_SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("start HTTP gRPC Gateway REST API server at %s\n", listener.Addr().String())

	wrappedMusLoggerHandler := gapi.HttpLogger(mux)
	if err != http.Serve(listener, wrappedMusLoggerHandler) {
		log.Fatalf("Failed to serve HTTP gRPC Gateway REST API: %v\n", err)
	}
}
