package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("..")
    require.NoError(t, err)
    require.NotEmpty(t, config.MigrationURL)
    require.NotEmpty(t, config.HTTPServerAddress)
    require.NotEmpty(t, config.GRPCServerAddress)
    require.NotEmpty(t, config.RedisAddress)
    require.NotEmpty(t, config.AccessTokenDuration)
    require.NotEmpty(t, config.RefreshTokenDuration)
    require.NotEmpty(t, config.EmailSenderName)
    require.NotEmpty(t, config.EmailSenderAddress)
    require.NotEmpty(t, config.EmailSenderPassword)
    require.NotEmpty(t, config.Environment)
    require.NotEmpty(t, config.DBSource)
    require.NotEmpty(t, config.TokenSymmetricKey)
}