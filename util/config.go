package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
    DBInitSchemaFile    string        `mapstructure:"DB_INIT_SCHEMA_FILE"`
    HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
    AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
    Environment         string        `mapstructure:"ENVIRONMENT"`
    DBSource            string        `mapstructure:"DB_SOURCE"`
    TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

// LoadConfig reads configuration from files or environment variables.
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigType("env")

    viper.SetConfigName("app")
    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    viper.SetConfigName("app_prod_overwrite")
    err = viper.MergeInConfig()
    if err != nil {
        return
    }

    viper.AutomaticEnv()    

    err = viper.Unmarshal(&config)
    return
}
