package gapi

import (
	"context"

	db "github.com/AdamDomagalsky/goes/bank/db/sqlc"
	"github.com/AdamDomagalsky/goes/bank/proto/pb"
	"github.com/AdamDomagalsky/goes/bank/util"
	pg "github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	violations := validateCreateUserRequest(in)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(in.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %v", err)
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Email:    in.GetEmail(),
		FullName: in.GetFullname(),
		Password: hashedPassword,
		Username: in.GetUsername(),
	})
	if err != nil {
		if pgErr, ok := err.(*pg.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				// ctx.JSON(http.StatusForbidden, errorResponse(err))
				return nil, status.Errorf(codes.AlreadyExists, "cannot create user: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	return &pb.CreateUserResponse{User: convertUser(user)}, nil
}

func validateCreateUserRequest(in *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := ValidateUsername(in.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := ValidateFullname(in.GetFullname()); err != nil {
		violations = append(violations, fieldViolation("fullname", err))
	}
	if err := ValidateEmail(in.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	if err := ValidatePassword(in.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
