package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
    Environment         string        `mapstructure:"ENVIRONMENT"`
    DBSource            string        `mapstructure:"DB_SOURCE"`
    DBInitSchemaFIle    string        `mapstructure:"DB_INIT_SCHEMA_FILE"`
    HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
    TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
    AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}
