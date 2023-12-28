package gapi

import (
	"context"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
    violations := validateVerifyEmailRequest(req)
    if violations != nil {
        return nil, invalidArgumentError(violations)
    }

    txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
        EmailID: req.GetEmailId(),
        SecretCode: req.GetSecretCode(),
    })
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to verify email")
    }

    resp := &pb.VerifyEmailResponse{
        IsVerified: txResult.User.IsEmailVerified,
    }
    return resp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
    if err := validation.ValidateEmailID(req.GetEmailId()); err != nil {
        violations = append(violations, fieldViolation("email_id", err))
    }

    if err := validation.ValidateSecretCode(req.GetSecretCode()); err != nil {
        violations = append(violations, fieldViolation("secret_code", err))
    }

    return violations
}
