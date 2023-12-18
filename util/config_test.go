package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("..")
    fmt.Printf("%+v", config)
    require.NoError(t, err)
    require.NotEmpty(t, config.DBInitSchemaFile)
    require.NotEmpty(t, config.HTTPServerAddress)
    require.NotEmpty(t, config.AccessTokenDuration)
    require.NotEmpty(t, config.Environment)
    require.NotEmpty(t, config.DBSource)
    require.NotEmpty(t, config.TokenSymmetricKey)
}