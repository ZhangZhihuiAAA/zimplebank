package gapi

import (
	"context"
	"errors"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
    violations := validateLoginUserRequest(req)
    if violations != nil {
        return nil, invalidArgumentError(violations)
    }

    user, err := server.store.GetUser(ctx, req.GetUsername())
    if err != nil {
        if errors.Is(err, db.ErrNoRows) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "failed to find user")
    }

    err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "incorrect password")
    }

    accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create access token")
    }

    refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
        user.Username,
        server.config.RefreshTokenDuration,
    )
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create refresh token")
    }

    metaData := server.extractMetadata(ctx)
    session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
        ID:           refreshPayload.ID,
        Username:     user.Username,
        RefreshToken: refreshToken,
        UserAgent:    metaData.UserAgent,
        ClientIp:     metaData.ClientIP,
        IsBlocked:    false,
        ExpiresAt:    refreshPayload.ExpiresAt,
    })
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create session")
    }

    resp := &pb.LoginUserResponse{
        User:                  convertUser(user),
        SessionId:             session.ID.String(),
        AccessToken:           accessToken,
        AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
        RefreshToken:          refreshToken,
        RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt),
    }

    return resp, nil
}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
    if err := validation.ValidateUsername(req.GetUsername()); err != nil {
        violations = append(violations, fieldViolation("username", err))
    }

    if err := validation.ValidatePassword(req.GetPassword()); err != nil {
        violations = append(violations, fieldViolation("password", err))
    }

    return violations
}