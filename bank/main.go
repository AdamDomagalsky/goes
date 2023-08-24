package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/AdamDomagalsky/goes/bank/api"
	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/gapi"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/util"
	_ "github.com/lib/pq" // blank import: side-effect init pg driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
			if err.Error() != "no change" {
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
	runGrpcAPIServer(config, store)
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

	grpcServer := grpc.NewServer()
	pb.RegisterBankServer(grpcServer, server)
	reflection.Register(grpcServer) // TODO kind of self documentation

	if err != grpcServer.Serve(listener) {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
