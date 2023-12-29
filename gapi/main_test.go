package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
    config := util.Config{
        TokenSymmetricKey:   util.RandomString(32),
        AccessTokenDuration: time.Minute,
    }

    server, err := NewServer(config, store, taskDistributor)
    require.NoError(t, err)

    return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
    accessToken, _, err := tokenMaker.CreateToken(username, duration)
    require.NoError(t, err)
    bearerToken := fmt.Sprintf("%s %s", AUTHORIZATION_TYPE_BEARER, accessToken)
    md := metadata.MD{
        AUTHORIZATION_HEADER_KEY: []string{
            bearerToken,
        },
    }
    return metadata.NewIncomingContext(context.Background(), md)
}