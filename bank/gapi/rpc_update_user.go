package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: add authentication
	violations := validateUpdateUserRequest(in)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := db.UpdateUserParams{
		Username: in.GetUsername(),
		FullName: sql.NullString{
			String: in.GetFullname(),
			Valid:  in.Fullname != nil,
		},
		Email: sql.NullString{
			String: in.GetEmail(),
			Valid:  in.Email != nil,
		},
	}

	if in.Password != nil {
		hashedPassword, err := util.HashPassword(in.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
		}
		args.Password = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		args.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	user, err := server.store.UpdateUser(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot update user: %v", err)
	}

	return &pb.UpdateUserResponse{User: convertUser(user)}, nil
}

func validateUpdateUserRequest(in *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(in.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if in.Password != nil {
		if err := ValidateFullname(in.GetFullname()); err != nil {
			violations = append(violations, fieldViolation("fullname", err))
		}
	}
	if in.Fullname != nil {
		if err := ValidateEmail(in.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if in.Email != nil {
		if err := ValidatePassword(in.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	return violations
}
