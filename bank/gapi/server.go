package gapi

import (
	"context"
	"database/sql"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/token"
	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedBankServer // TODOuse it to be able gradually implement grpc methods
	store                      db.Store
	config                     util.Config
	tokenMaker                 token.Maker
}

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Fullname:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}

}

func (server *Server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(in.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username: in.GetUsername(),
		FullName: in.GetFullname(),
		Email:    in.GetEmail(),
		Password: hashedPassword,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	return &pb.CreateUserResponse{User: convertUser(user)}, nil
}

func (server *Server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, in.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "cannot find user: %v", err) // should return Unauthenticated with invalid username/password but for demo it's ok
		}
		return nil, status.Errorf(codes.Internal, "cannot get user: %v", err)
	}

	err = util.CheckPasswordHash(in.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username/password: %v", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.ACCESS_TOKEN_DURATION,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %v", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.REFRESH_TOKEN_DURATION,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %v", err)
	}

	metadata := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiresAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %v", err)
	}

	return &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiresAt),
	}, nil

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
