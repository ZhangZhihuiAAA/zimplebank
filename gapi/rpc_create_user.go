package gapi

import (
	"context"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    violations := validateCreateUserRequest(req)
    if violations != nil {
        return nil, invalidArgumentError(violations)
    }

    hashedPassword, err := util.HashPassword(req.GetPassword())
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to hash password")
    }

    arg := db.CreateUserParams{
        Username:       req.GetUsername(),
        HashedPassword: hashedPassword,
        FullName:       req.GetFullName(),
        Email:          req.GetEmail(),
    }

    user, err := server.store.CreateUser(ctx, arg)
    if err != nil {
        errCode := db.ErrorCode(err)
        if errCode == db.UNIQUE_VIOLATION {
            return nil, status.Error(codes.AlreadyExists, "username already exists")
        }

        return nil, status.Error(codes.Internal, "failed to create user")
    }

    resp := &pb.CreateUserResponse{
        User: convertUser(user),
    }
    return resp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
    if err := validation.ValidateUsername(req.GetUsername()); err != nil {
        violations = append(violations, fieldViolation("username", err))
    }

    if err := validation.ValidateFullName(req.GetFullName()); err != nil {
        violations = append(violations, fieldViolation("full_name", err))
    }

    if err := validation.ValidateEmail(req.GetEmail()); err != nil {
        violations = append(violations, fieldViolation("email", err))
    }

    if err := validation.ValidatePassword(req.GetPassword()); err != nil {
        violations = append(violations, fieldViolation("password", err))
    }

    return violations
}