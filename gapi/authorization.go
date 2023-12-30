package gapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
    AUTHORIZATION_HEADER_KEY  = "authorization"
    AUTHORIZATION_TYPE_BEARER = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, fmt.Errorf("missing metadata")
    }

    values := md.Get(AUTHORIZATION_HEADER_KEY)
    if len(values) == 0 {
        return nil, fmt.Errorf("missing authorization header")
    }

    authHeader := values[0]
    fields := strings.Fields(authHeader)
    if len(fields) < 2 {
        return nil, fmt.Errorf("invalid authorization header format")
    }

    authType := strings.ToLower(fields[0])
    if authType != AUTHORIZATION_TYPE_BEARER {
        return nil, fmt.Errorf("unsupported authorization type: %s", authType)
    }

    accessToken := fields[1]
    payload, err := server.tokenMaker.VerifyToken(accessToken)
    if err != nil {
        return nil, fmt.Errorf("invalid access token: %s", err)
    }

    if !hasPermission(payload.Role, accessibleRoles) {
        return nil, errors.New("permission denied")
    }

    return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
    for _, role := range accessibleRoles {
        if userRole == role {
            return true
        }
    }
    return false
}