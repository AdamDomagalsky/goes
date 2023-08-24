package gapi

import (
	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/token"
	"github.com/AdamDomagalsky/goes/bank/util"
)

type Server struct {
	pb.UnimplementedBankServer // TODO use it to be able gradually implement grpc methods
	store                      db.Store
	config                     util.Config
	tokenMaker                 token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPASETOMaker(config.SYMMETRIC_KEY)
	// tokenMaker, err := token.NewJWTMaker(config.SYMMETRIC_KEY)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
