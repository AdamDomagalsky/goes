package gapi

import (
	"context"
	"database/sql"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/util"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(in)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

func validateLoginUserRequest(in *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(in.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := ValidatePassword(in.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return
}
